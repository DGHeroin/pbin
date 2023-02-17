package server

import (
    "errors"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/glebarez/sqlite"
    "github.com/patrickmn/go-cache"
    "github.com/spf13/cobra"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "io"
    "net/http"
    "os"
    "path"
    "pbin/common/ider"
    "pbin/common/logger"
    "pbin/common/utils"
    "strings"
    "sync/atomic"
    "time"
)

var (
    Cmd = &cobra.Command{
        Use: "server <args>",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runServer()
        },
    }
)

var (
    address       string
    dataDir       string
    dbConn        string
    hostname      string
    keepTime      time.Duration
    MaxBytes      int64
    cc            = cache.New(5*time.Minute, 10*time.Minute)
    CapacityBytes int64
    useCapacity   int64
)

func init() {
    Cmd.PersistentFlags().StringVar(&address, "address", "127.0.0.1:18080", "serve address")
    Cmd.PersistentFlags().StringVar(&dataDir, "dir", "./data/", "data dir")
    Cmd.PersistentFlags().StringVar(&dbConn, "db", "sqlite://data/pbin.db", "data dir")
    Cmd.PersistentFlags().DurationVar(&keepTime, "keep", time.Hour*24*3, "keep file")
    Cmd.PersistentFlags().Int64Var(&MaxBytes, "max", 30<<20, "max file size file")
    Cmd.PersistentFlags().Int64Var(&CapacityBytes, "capacity", 1<<10, "max capacity storage")
    Cmd.PersistentFlags().StringVar(&hostname, "hostname", "http://127.0.0.1:18080", "")
}
func initDB() (db *gorm.DB, err error) {
    if strings.HasPrefix(dbConn, "sqlite://") {
        connStr := strings.Replace(dbConn, "sqlite://", "", 1)
        db, err = gorm.Open(sqlite.Open(connStr), &gorm.Config{})

    }
    if strings.HasPrefix(dbConn, "mysql://") {
        connStr := strings.Replace(dbConn, "mysql://", "", 1)
        db, err = gorm.Open(mysql.Open(connStr), &gorm.Config{})

    }
    if strings.HasPrefix(dbConn, "postgres://") {
        connStr := strings.Replace(dbConn, "postgres://", "", 1)
        db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
    }
    if db == nil {
        return nil, errors.New("unknown db")
    }
    err = db.AutoMigrate(&FileInfo{})
    useCapacity = GetFileTotalSize(db)
    return db, err
}

func runServer() error {
    db, err := initDB()
    if err != nil {
        return err
    }
    if db == nil {
        panic(err)
    }
    go gcOldFile(db)

    r := gin.Default()
    r.Use(gin.Recovery(), cors.Default())

    r.GET("/get/:id", func(c *gin.Context) {
        id := c.Param("id")
        var (
            fi  *FileInfo
            err error
        )
        p, ok := cc.Get(id)
        if ok {
            fi = p.(*FileInfo)
        } else {
            fi, err = GetFileWithId(db, id)
            if err != nil {
                c.String(http.StatusBadRequest, err.Error())
                logger.Info(err)
                return
            }
        }
        if fi == nil {
            c.Status(http.StatusNotFound)
            return
        }

        FileUpdateTime(db, fi)
        cc.SetDefault(id, fi)
        c.Writer.WriteHeader(http.StatusOK)
        c.Header("Content-Disposition", "attachment; filename="+fi.Filename)
        c.Header("Content-Type", fi.ContentType)
        c.File(fi.Path)
    })
    r.POST("/put", func(c *gin.Context) {
        // c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBytes)
        file, header, err := c.Request.FormFile("upload")
        if err != nil {
            c.String(http.StatusBadRequest, err.Error())
            logger.Info("get upload file error:", err)
            return
        }
        fileLength := header.Size
        if fileLength > MaxBytes {
            c.String(http.StatusBadRequest, "request body too large")
            return
        }
        if fileLength+useCapacity > CapacityBytes {
            c.String(http.StatusBadRequest, "server capacity full")
            return
        }
        atomic.AddInt64(&useCapacity, fileLength)
        contentType := c.GetHeader("X-Content-Type")
        if contentType == "" {
            contentType = c.ContentType()
        }
        var fi = &FileInfo{
            TrackID:     ider.Generate(),
            ContentType: contentType,
            Filename:    header.Filename,
            IP:          GetClientIPByHeaders(c.Request),
            Length:      fileLength,
        }
        needClean := false
        bDir := path.Join(dataDir, time.Now().Format("2006-01-02"), fi.TrackID)
        if err := os.MkdirAll(bDir, os.ModePerm); err != nil {
            c.String(http.StatusBadRequest, "Bad request")
            return
        }
        defer func() {
            recover()
            if needClean {
                _ = os.RemoveAll(bDir)
                atomic.AddInt64(&useCapacity, -fileLength)
            }
        }()

        savePath := path.Join(bDir, "file")
        out, err := os.Create(savePath)
        if err != nil {
            c.String(http.StatusInternalServerError, "internal error")
            logger.Info("create file error:", err)
            needClean = true
            return
        }
        defer func() {
            _ = out.Close()
        }()
        _, err = io.Copy(out, file)
        if err != nil {
            c.String(http.StatusInternalServerError, "internal error")
            logger.Info("copy file error:", err)
            needClean = true
            return
        }
        fi.Path = savePath
        err = AddFileInfo(db, fi)
        if err != nil {
            c.String(http.StatusInternalServerError, "internal error")
            logger.Info("save db file error:", err)
            needClean = true
            return
        }
        c.String(http.StatusOK, "%v/get/%v", hostname, fi.TrackID)
    })

    return r.Run(address)
}

func gcOldFile(db *gorm.DB) {
    for {
        files := GetExpiredFile(db, time.Now().Add(-keepTime))
        for _, fi := range files {
            err := os.RemoveAll(path.Dir(fi.Path))
            if err != nil {
                logger.Info("删除失败:", err)
            } else {
                atomic.AddInt64(&useCapacity, -fi.Length)
                logger.Infof("[回收] id:%v content type:%v size:%v 可用空间:%v/%v",
                    fi.TrackID, fi.ContentType, utils.ByteSize(fi.Length),
                    utils.ByteSize(CapacityBytes-useCapacity), utils.ByteSize(CapacityBytes),
                )
            }

        }
        RemFiles(db, files)
        time.Sleep(time.Minute)
    }
}
