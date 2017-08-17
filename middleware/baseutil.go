package middleware

import (
	"net/http"
	"regexp"

	"github.com/suconghou/utilgo"
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
	Handler func(http.ResponseWriter, *http.Request, []string)
}

var usecachefilter = func(out http.Header, res *http.Response) int {
	out.Del("Set-Cookie")
	utilgo.UseHTTPCache(out, 604800)
	return res.StatusCode
}

func dispatch(w http.ResponseWriter, r *http.Request, match []string, route []routeInfo, fallback func(w http.ResponseWriter, r *http.Request, match []string)) {
	uri := match[1]
	found := false
	for _, p := range route {
		if p.Reg.MatchString(uri) {
			found = true
			p.Handler(w, r, p.Reg.FindStringSubmatch(uri))
			break
		}
	}
	if !found {
		fallback(w, r, match)
	}
}

func cleanHeader(r *http.Request, headers []string) *http.Request {
	for _, k := range headers {
		r.Header.Del(k)
	}
	return r
}
