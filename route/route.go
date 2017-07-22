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

// 路由添加
var RoutePath = []routeInfo{
	{regexp.MustCompile(`^/files/(.+)$`), files},
	{regexp.MustCompile(`^/imgs/(.+)$`), imgs},
	{regexp.MustCompile(`^/api/media/(.+)$`), middleware.MediaStreamApi},
	{regexp.MustCompile(`^/media/(.+)$`), middleware.MediaStream},
	{regexp.MustCompile(`^/api/net/(.+)$`), middleware.NetStreamApi},
	{regexp.MustCompile(`^/net/(.+)$`), middleware.NetStream},
}

func files(w http.ResponseWriter, r *http.Request, match []string) {

}

func imgs(w http.ResponseWriter, r *http.Request, match []string) {

}
