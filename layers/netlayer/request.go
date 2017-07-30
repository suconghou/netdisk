package netlayer

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/suconghou/utilgo"
)

var (
	range1 = regexp.MustCompile(`^--range:(\d+)-(\d+)$`)
	range2 = regexp.MustCompile(`^--range:(\d+)-$`)
)

// Get return response data
func Get(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Post retrun response data
func Post(url string, contentType string, body io.Reader) ([]byte, error) {
	response, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bodyStr, err
}

// PostWait wait a second and retrun response data
func PostWait(url string, contentType string, body io.Reader) ([]byte, error) {
	response, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	time.Sleep(time.Second)
	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bodyStr, nil
}

// ParseThreadThunkStartEnd return thread thunk start end
func ParseThreadThunkStartEnd(thread int32, thunk int64, start int64, end int64) (int32, int64, int64, int64) {
	if utilgo.HasFlag("--most") {
		thread = thread * 4
	} else if utilgo.HasFlag("--fast") {
		thread = thread * 2
	} else if utilgo.HasFlag("--slow") {
		thread = thread / 2
	}
	if utilgo.HasFlag("--thin") {
		thunk = thunk / 8
	} else if utilgo.HasFlag("--fat") {
		thunk = thunk * 4
	}
	for _, item := range os.Args {
		if range1.MatchString(item) {
			matches := range1.FindStringSubmatch(item)
			start, _ = strconv.ParseInt(matches[1], 10, 64)
			end, _ = strconv.ParseInt(matches[2], 10, 64)
			break
		} else if range2.MatchString(item) {
			matches := range2.FindStringSubmatch(item)
			start, _ = strconv.ParseInt(matches[1], 10, 64)
			break
		}
	}
	return thread, thunk, start, end
}

// ParseCookieUaRefer return http.Header
func ParseCookieUaRefer() http.Header {
	reqHeader := http.Header{}
	if value, err := utilgo.GetParam("--cookie"); err == nil {
		reqHeader.Add("Cookie", value)
	}
	if value, err := utilgo.GetParam("--ua"); err == nil {
		reqHeader.Add("User-Agent", value)
	} else {
		reqHeader.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	}
	if value, err := utilgo.GetParam("--refer"); err == nil {
		reqHeader.Add("Referer", value)
	}
	return reqHeader
}
