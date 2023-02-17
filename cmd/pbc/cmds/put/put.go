package put

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/spf13/cobra"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "net/url"
    "os"
    "path"
    "pbin/common/utils"
)

var (
    Cmd = &cobra.Command{
        Use:  "put",
        RunE: RunAsClient,
    }
)
var (
    cipherKey   string
    contentType string
    plaintext   string
    baseUrl     string
    limit       int
)

func init() {
    Cmd.PersistentFlags().StringVar(&cipherKey, "k", "", "cipher key")
    Cmd.PersistentFlags().StringVar(&contentType, "c", "application/octet-stream", "content type")
    Cmd.PersistentFlags().StringVar(&plaintext, "s", "", "cipher content")
    Cmd.PersistentFlags().StringVar(&baseUrl, "h", "https://pbin.moyuta.com", "server")
    Cmd.PersistentFlags().IntVar(&limit, "limit", 0, "limit download")

    if cipherKey == "" {
        if hostname, err := os.Hostname(); err != nil {
            cipherKey = "pbin"
        } else {
            cipherKey = hostname
        }
    }
}

var (
    ErrInvalidArgs = errors.New("invalid args")
)

func RunAsClient(_ *cobra.Command, args []string) error {
    var (
        isFile = true
        r      io.Reader
    )
    isFile = plaintext == ""
    if !isFile {
        r = bytes.NewBufferString(plaintext)
    } else {
        if len(args) != 1 {
            return ErrInvalidArgs
        }
        inFile, err := os.Open(args[0])
        if err != nil {
            return err
        }
        defer inFile.Close()
        r = inFile
    }
    reader, err := utils.AESReader(cipherKey, r)
    if err != nil {
        return err
    }

    buf := new(bytes.Buffer)
    w := multipart.NewWriter(buf)
    if isFile {
        part, err := w.CreateFormFile("upload", path.Base(args[0]))
        if err != nil {
            return err
        }
        _, err = io.Copy(part, reader)
        if err != nil {
            return err
        }
    } else {
        part, err := w.CreateFormFile("upload", "raw.txt")
        if err != nil {
            return err
        }
        _, err = io.Copy(part, reader)
        if err != nil {
            return err
        }
    }
    err = w.Close()
    if err != nil {
        return err
    }

    request, err := http.NewRequest(http.MethodPost, baseUrl+"/put", buf)
    if err != nil {
        return err
    }
    request.Header.Set("X-Content-Type", contentType)
    request.Header.Add("Content-Type", w.FormDataContentType())

    resp, err := http.DefaultClient.Do(request)
    if err != nil {
        return err
    }
    if resp.StatusCode != http.StatusOK {
        data, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("%v:%v", resp.Status, string(data))
    }
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    u, err := url.Parse(string(data))
    if err != nil {
        return err
    }
    fmt.Println(u)
    fmt.Println("pbc get", path.Base(u.Path))
    return nil
}
