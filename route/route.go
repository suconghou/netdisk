package route

import (
	"net/http"
	"regexp"

	"netdisk/middleware"
)

// 路由定义
type routeInfo struct {
	Reg     *regexp.Regexp
	Handler func(http.ResponseWriter, *http.Request, []string) error
}

// RoutePath defines all routes
var RoutePath = []routeInfo{
	{regexp.MustCompile(`^/net/(.+)$`), middleware.NetStreamAPI},
	{regexp.MustCompile(`^/(?:(https?):/?)?(?:[\w\-]+\.)+[\w\-]+(?::\d+)?/.*$`), middleware.Pipe},
}
