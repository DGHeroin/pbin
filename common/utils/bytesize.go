package utils

import "fmt"

func ByteSize(n int64) string {
    v := float64(n)
    u := "B"
    if v > 1024 {
        v /= 1024
        u = "K"
    }

    if v > 1024 {
        v /= 1024
        u = "M"
    }
    if v > 1024 {
        v /= 1024
        u = "G"
    }
    if v > 1024 {
        v /= 1024
        u = "T"
    }
    if v > 1024 {
        v /= 1024
        u = "P"
    }

    return fmt.Sprintf("%.2f%v", v, u)
}
