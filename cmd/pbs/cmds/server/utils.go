package server

import (
    "net/http"
    "strings"
)

func GetClientIPByHeaders(req *http.Request) string {
    var ipSlice []string
    ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
    ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
    ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))
    ipSlice = append(ipSlice, req.RemoteAddr)

    for _, v := range ipSlice {
        if v != "" {
            ips := strings.Split(v, ",")
            if len(ips) == 0 {
                return ""
            }
            return strings.TrimSpace(ips[0])
        }
    }
    return HTTPGetClientIPByHeader(req.Header)
}
func HTTPGetClientIPByHeader(Header http.Header) string {
    var (
        ipSlice []string
        all     []string
    )

    ipSlice = append(ipSlice, Header.Get("X-Forwarded-For"))
    ipSlice = append(ipSlice, Header.Get("x-forwarded-for"))
    ipSlice = append(ipSlice, Header.Get("X-FORWARDED-FOR"))

    for _, v := range ipSlice {
        if v != "" {
            ips := strings.Split(v, ",")
            if len(ips) == 0 {
                return ""
            }
            all = append(all, ips...)
        }
    }

    for _, v := range all {
        if strings.HasPrefix(v, "127.0.0.1") {
            continue
        }
        return v
    }
    return ""
}
