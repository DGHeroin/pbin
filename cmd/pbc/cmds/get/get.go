package get

import (
    "fmt"
    "github.com/spf13/cobra"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "pbin/common/logger"
    "pbin/common/utils"
    "strings"
)

var (
    Cmd = &cobra.Command{
        Use: "get",
        RunE: func(cmd *cobra.Command, args []string) error {
            return download()
        },
    }
)
var (
    id        string
    baseUrl   string
    cipherKey string
)

func init() {
    Cmd.PersistentFlags().StringVar(&id, "id", "", "upload id")
    Cmd.PersistentFlags().StringVar(&baseUrl, "h", "https://pbin.moyuta.com", "server")

    if cipherKey == "" {
        if hostname, err := os.Hostname(); err != nil {
            cipherKey = "pbin"
        } else {
            cipherKey = hostname
        }
    }
}
func download() error {
    resp, err := http.Get(baseUrl + "/get/" + id)
    if err != nil {
        return err
    }
    if resp.StatusCode != http.StatusOK {
        data, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("%v:%v", resp.Status, string(data))
    }
    filename := id
    if disp := resp.Header.Get("Content-Disposition"); disp != "" {
        filename = strings.ReplaceAll(disp, "attachment; filename=", "")
    }
    logger.Info("文件名", filename)
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    w, err := utils.AESWriter(cipherKey, file)
    _, err = io.Copy(w, resp.Body)
    return err
}
