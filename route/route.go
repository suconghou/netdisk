package route

import (
	"net/http"
	"netdisk/middleware"
	"regexp"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string)
}

// RoutePath defines all routes
var RoutePath = []routeInfo{
	{regexp.MustCompile(`^/files/(.+)$`), files},
	{regexp.MustCompile(`^/imgs/(.+)$`), imgs},
	{regexp.MustCompile(`^/api/media/(.+)$`), middleware.MediaStreamAPI},
	{regexp.MustCompile(`^/media/(.+)$`), middleware.MediaStream},
	{regexp.MustCompile(`^/net/(.+)$`), middleware.NetStreamAPI},
	{regexp.MustCompile(`^/(?:(https?):/?)?(?:\w+\.)+\w+(?::\d+)?/.*$`), middleware.Pipe},
}

func files(w http.ResponseWriter, r *http.Request, match []string) {

}

func imgs(w http.ResponseWriter, r *http.Request, match []string) {

}
