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
	{regexp.MustCompile(`^/files/(.+)$`), middleware.Files},
	{regexp.MustCompile(`^/imgs/(.+)$`), middleware.Imgs},
	{regexp.MustCompile(`^/play/(.+)$`), middleware.MediaStreamAPI},
	{regexp.MustCompile(`^/net/(.+)$`), middleware.NetStreamAPI},
	{regexp.MustCompile(`^/(?:(https?):/?)?(?:\w+\.)+\w+(?::\d+)?/.*$`), middleware.Pipe},
}
