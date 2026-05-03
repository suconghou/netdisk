package middleware

import (
	"net/http"
	"regexp"
)

type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string) error
}

func dispatch(w http.ResponseWriter, r *http.Request, match []string, route []routeInfo, fallback func(w http.ResponseWriter, r *http.Request, match []string) error) error {
	uri := match[1]
	for _, p := range route {
		matches := p.Reg.FindStringSubmatch(uri)
		if matches != nil {
			return p.Handler(w, r, matches)
		}
	}
	return fallback(w, r, match)
}
