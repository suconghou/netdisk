package route

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/suconghou/netdisk/middleware"
)

var (
	apiv1, _ = url.Parse("http://api.suconghou.cn/")
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
