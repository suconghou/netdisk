package util

import (
	"io"
	"net/http"
	"time"
)

var (
	client = &http.Client{Timeout: time.Minute}

	fwdHeaders = []string{
		"User-Agent",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"If-Modified-Since",
		"If-None-Match",
		"Range",
		"Content-Length",
		"Content-Type",
	}
	exposeHeaders = []string{
		"Accept-Ranges",
		"Content-Range",
		"Content-Length",
		"Content-Type",
		"Content-Encoding",
		"Date",
		"Expires",
		"Last-Modified",
		"Etag",
		"Cache-Control",
	}
)

// ProxyURL proxy target resource
func ProxyURL(w http.ResponseWriter, r *http.Request, target string, headers http.Header) error {
	var (
		reqHeader = http.Header{}
	)
	req, err := http.NewRequest(r.Method, target, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	req.Header = copyHeader(r.Header, reqHeader, fwdHeaders)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer resp.Body.Close()
	to := w.Header()
	copyHeader(resp.Header, to, exposeHeaders)
	for k, v := range headers {
		if k != "" && len(v) == 1 {
			to.Set(k, v[0])
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	return err
}

func copyHeader(from http.Header, to http.Header, headers []string) http.Header {
	for _, k := range headers {
		if v := from.Get(k); v != "" {
			to.Set(k, v)
		}
	}
	return to
}
