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

type WriteCounter struct {
	Total     uint64
	Size      uint64
	StartTime time.Time
}

// func (wc *WriteCounter) Write(p []byte) (int, error) {
// 	n := len(p)
// 	wc.Total += uint64(n)
// 	var i int = int(float64(wc.Total) / float64(wc.Size) * 100)
// 	duration := time.Since(wc.StartTime).Seconds()
// 	speed := float64(wc.Total) / 1024 / duration
// 	leftTime := float64(wc.Size)/1024/speed - duration
// 	fmt.Printf("\r%s%d%% %s %.2fKB/s %.1f %.1f  ", util.Bar(i, 25), i, util.ByteFormat(wc.Total), speed, duration, leftTime)
// 	return n, nil
// }

type PutprogressReporter struct {
	R         io.Reader
	Total     uint64
	Size      uint64
	StartTime time.Time
}

func (pr *PutprogressReporter) Read(p []byte) (int, error) {
	return pr.R.Read(p)
	// n, err := pr.R.Read(p)
	// pr.Total += uint64(n)
	// var i int = int(float64(pr.Total) / float64(pr.Size) * 100)
	// duration := time.Since(pr.StartTime).Seconds()
	// speed := float64(pr.Total) / 1024 / duration
	// leftTime := float64(pr.Size)/1024/speed - duration
	// fmt.Printf("\r%s%d%% %s %.2fKB/s %.1f %.1f  ", util.Bar(i, 25), i, util.ByteFormat(pr.Total), speed, duration, leftTime)
	// return n, err
}

// func Download(url string, saveas string, size uint64, hash string) {
// 	start := time.Now()
// 	res, err := http.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	f, err := os.Create(saveas)
// 	if err != nil {
// 		panic(err)
// 	}
// 	counter := &WriteCounter{Size: size, StartTime: start}
// 	src := io.TeeReader(res.Body, counter)
// 	count, err := io.Copy(f, src)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	if count < 1 {
// 		os.Exit(1)
// 	}
// 	end := time.Since(start)
// 	speed := float64(size/1024) / end.Seconds()
// 	fmt.Printf("\n下载完毕,耗时%s,%.2fKB/s,校验MD5中...\n", end.String(), speed)
// 	util.PrintMd5(saveas)
// }

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
