package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/suconghou/utilgo"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/youtubeVideoParser"
)

var mediaroute = []routeInfo{
	{regexp.MustCompile(`^(large|medium|small)/([\w\-]{6,12})\.(mp4|flv|webm|3gp|json)$`), youtubeVideo},
	{regexp.MustCompile(`^(large|medium|small)/([\w\-]{6,12})\.(jpg|webp)$`), youtubeImage},
	{regexp.MustCompile(`^([\w\-]{6,12})/(\d+)\.(mp4|flv|webm|3gp|mp3)$`), youtubeVideoItag},
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
		bs, err := youtubeVideoParser.GetYoutubeVideoInfo(id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		utilgo.JSONPut(w, bs, true, 86400)
	} else {
		if url, err := youtubeVideoParser.GetYoutubeVideoURL(id, ext, quality); err == nil {
			fastload.Pipe(w, r, url, usecachefilter, 30, nil)
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
}

func youtubeVideoItag(w http.ResponseWriter, r *http.Request, match []string) {
	id, itag := match[1], match[2]
	info, err := youtubeVideoParser.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if v, ok := info.Streams[itag]; ok {
		fastload.Pipe(w, r, v.URL, usecachefilter, 30, nil)
	} else {
		http.Error(w, fmt.Sprintf("itag %s for %s not found", itag, id), 404)
	}
}

func youtubeImage(w http.ResponseWriter, r *http.Request, match []string) {
	quality, id, ext := match[1], match[2], match[3]
	url := youtubeVideoParser.GetYoutubeImageURL(id, ext, quality)
	fastload.Pipe(w, r, url, usecachefilter, 30, nil)
}
