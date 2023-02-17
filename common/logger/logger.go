package logger

import (
    "fmt"
    "github.com/sirupsen/logrus"
    "runtime"
    "strings"
)

var logger = logrus.New()

func init() {

}
func SetLogFormatter(formatter logrus.Formatter) {
    logger.Formatter = formatter
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
    if logger.Level >= logrus.InfoLevel {
        entry := logger.WithFields(logrus.Fields{})
        entry.Data["file"] = fileInfo(2)
        entry.Info(args...)
    }
}
func Infof(format string, args ...interface{}) {
    if logger.Level >= logrus.InfoLevel {
        entry := logger.WithFields(logrus.Fields{})
        entry.Data["file"] = fileInfo(2)
        entry.Infof(format, args...)
    }
}
func fileInfo(skip int) string {
    _, file, line, ok := runtime.Caller(skip)
    if !ok {
        file = "<???>"
        line = 1
    } else {
        slash := strings.LastIndex(file, "/")
        if slash >= 0 {
            file = file[slash+1:]
        }
    }
    return fmt.Sprintf("%s:%d", file, line)
}
