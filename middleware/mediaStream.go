package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/youtubeVideoParser"
)

var mediaroute = []routeInfo{
	{regexp.MustCompile(`^/(large|medium|small)/([\w\-]{6,12})\.(mp4|flv|webm|3gp|mp3|json)$`), youtubeVideo},
	{regexp.MustCompile(`^/(large|medium|small)/([\w\-]{6,12})\.(jpg|webp)$`), youtubeImage},
}

// MediaStreamAPI response json data
func MediaStreamAPI(w http.ResponseWriter, r *http.Request, match []string) {
	dispatch(w, r, match, mediaroute, func(w http.ResponseWriter, r *http.Request, match []string) {
		http.NotFound(w, r)
	})
}

func youtubeVideo(w http.ResponseWriter, r *http.Request, match []string) {
	quality, id, ext := match[1], match[2], match[3]
	if ext == "json" {
		info, err := youtubeVideoParser.GetYoutubeVideoInfo(id)
	} else {
		if url, err := youtubeVideoParser.GetYoutubeVideoURL(id, ext, quality); err == nil {
			fastload.Pipe(w, r, url, usecachefilter, 30, nil)
		} else {
			http.Error(w, fmt.Sprintf("%s", err), 500)
		}
	}
}

func youtubeImage(w http.ResponseWriter, r *http.Request, match []string) {
	quality, id, ext := match[1], match[2], match[3]
	url := youtubeVideoParser.GetYoutubeImageURL(id, ext, quality)
	fastload.Pipe(w, r, url, usecachefilter, 30, nil)
}
