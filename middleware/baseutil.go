package middleware

import (
	"net/http"
	"regexp"
)

var (
	xheaders = []string{
		"X-Forwarded-For",
		"X-Forwarded-Host",
		"X-Forwarded-Server",
		"X-Forwarded-Port",
		"X-Forwarded-Proto",
		"X-Client-Ip",
	}
)

type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string) error
}

var usecachefilter = func(out *http.Header, res *http.Header, status int) int {
	out.Del("Set-Cookie")
	out.Del("Content-Disposition")
	out.Set("Access-Control-Allow-Origin", "*")
	out.Set("Access-Control-Allow-Headers", "Range, Origin, X-Requested-With, Content-Type, Content-Length, Accept, Accept-Encoding, Cache-Control, Expires")
	return status
}

func dispatch(w http.ResponseWriter, r *http.Request, match []string, route []routeInfo, fallback func(w http.ResponseWriter, r *http.Request, match []string) error) error {
	uri := match[1]
	for _, p := range route {
		if p.Reg.MatchString(uri) {
			return p.Handler(w, r, p.Reg.FindStringSubmatch(uri))
		}
	}
	return fallback(w, r, match)
}

func cleanHeader(r *http.Request, headers []string) *http.Request {
	for _, k := range headers {
		r.Header.Del(k)
	}
	return r
}
