package server

import (
    "gorm.io/gorm"
    "log"
    "time"
)

type FileInfo struct {
    gorm.Model
    TrackID     string
    ContentType string
    Filename    string
    IP          string
    Path        string
    Length      int64
}

func AddFileInfo(db *gorm.DB, fi *FileInfo) error {
    return db.Create(fi).Error
}
func GetExpiredFile(db *gorm.DB, ts time.Time) []*FileInfo {
    var files []*FileInfo
    db.Where("updated_at < ?", ts).Find(&files)
    return files
}
func RemFiles(db *gorm.DB, files []*FileInfo) {
    if len(files) == 0 {
        return
    }
    db.Delete(files)
}
func GetFileWithId(db *gorm.DB, id string) (*FileInfo, error) {
    var fi FileInfo
    tx := db.Where("track_id = ?", id).Find(&fi)
    if tx.RowsAffected != 1 {
        return nil, nil
    }
    return &fi, tx.Error
}
func FileUpdateTime(db *gorm.DB, fi *FileInfo) {
    db.Model(&FileInfo{}).
        Where("id = ?", fi.ID).
        Update("updated_at", time.Now())
}
func GetFileTotalSize(db *gorm.DB) int64 {
    result := struct {
        Total int64
    }{}
    db.Model(&FileInfo{}).Select("id, sum(length) as total").Scan(&result)
    log.Println("总数", result)
    return result.Total
}
