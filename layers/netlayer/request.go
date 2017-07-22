package netlayer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"netdisk/util"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var userRangeFullFormatReg = regexp.MustCompile(`^--range:(\d+)-(\d+)$`)
var userRangeHalfFormatReg = regexp.MustCompile(`^--range:(\d+)-$`)

func init() {

}

func Get(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("%s", err)))
		os.Exit(1)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func Post(url string, contentType string, body io.Reader) []byte {
	response, err := http.Post(url, contentType, body)
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("%s", err)))
		os.Exit(1)
	}
	defer response.Body.Close()
	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return bodyStr
}

func PostWait(url string, contentType string, body io.Reader) []byte {
	response, err := http.Post(url, contentType, body)
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("%s", err)))
		os.Exit(1)
	}
	defer response.Body.Close()
	time.Sleep(time.Second)
	bodyStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return bodyStr
}

type WriteCounter struct {
	Total     uint64
	Size      uint64
	StartTime time.Time
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	var i int = int(float64(wc.Total) / float64(wc.Size) * 100)
	duration := time.Since(wc.StartTime).Seconds()
	speed := float64(wc.Total) / 1024 / duration
	leftTime := float64(wc.Size)/1024/speed - duration
	fmt.Printf("\r%s%d%% %s %.2fKB/s %.1f %.1f  ", util.Bar(i, 25), i, util.ByteFormat(wc.Total), speed, duration, leftTime)
	return n, nil
}

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

func Download(url string, saveas string, size uint64, hash string) {
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(saveas)
	if err != nil {
		panic(err)
	}
	counter := &WriteCounter{Size: size, StartTime: start}
	src := io.TeeReader(res.Body, counter)
	count, err := io.Copy(f, src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if count < 1 {
		os.Exit(1)
	}
	end := time.Since(start)
	speed := float64(size/1024) / end.Seconds()
	fmt.Printf("\n下载完毕,耗时%s,%.2fKB/s,校验MD5中...\n", end.String(), speed)
	util.PrintMd5(saveas)
}

func callPlayer(file string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("PotPlayerMini.exe", file)
		cmd.Start()
	} else {
		cmd := exec.Command("mpv", file)
		cmd.Start()
	}
}

func tryGetUserRange(start uint64, end uint64, rangeAble bool) (uint64, uint64, bool) {
	var userRangeStart uint64
	var userRangeEnd uint64
	var matched bool = false
	for _, item := range os.Args {
		if userRangeHalfFormatReg.MatchString(item) {
			matches := userRangeHalfFormatReg.FindStringSubmatch(item)
			startInt, _ := strconv.ParseUint(matches[1], 10, 64)
			userRangeStart = startInt
			userRangeEnd = end
			matched = true
			break
		} else if userRangeFullFormatReg.MatchString(item) {
			matches := userRangeFullFormatReg.FindStringSubmatch(item)
			startInt, _ := strconv.ParseUint(matches[1], 10, 64)
			endInt, _ := strconv.ParseUint(matches[2], 10, 64)
			userRangeStart = startInt
			userRangeEnd = endInt
			matched = true
			break
		}
	}
	if matched {
		if userRangeStart > userRangeEnd {
			panic("error range: start no less than end")
		} else if (userRangeEnd > end) || (userRangeStart > end) {
			panic("error range: range out of file end")
		}
		if rangeAble {
			return userRangeStart, userRangeEnd, true
		} else {
			panic("download is not rangeable , reset to default")
			return start, end, false
		}
	}
	if start > end {
		panic("error range: start no less than end")
	}
	return start, end, false
}

func getThreadThunk(rangeAble bool) (uint8, uint32) {
	var thread uint8 = 8
	var thunk uint32 = 524288 * 4
	if util.HasFlag("--most") {
		thread = thread * 4
	} else if util.HasFlag("--fast") {
		thread = thread * 2
	} else if util.HasFlag("--slow") {
		thread = thread / 2
	}
	if util.HasFlag("--thin") {
		thunk = thunk / 8
	} else if util.HasFlag("--fat") {
		thunk = thunk * 4
	}
	if !rangeAble {
		panic("download is not rangeable , using one thread")
		thread = 1
	}
	return thread, thunk
}

func ParseCookieUaRefer() {
	var headers = map[string]string{}
	if value, err := util.GetParam("--cookie"); err == nil {
		headers["Cookie"] = value
	}
	if value, err := util.GetParam("--ua"); err == nil {
		headers["User-Agent"] = value
	}
	if value, err := util.GetParam("--refer"); err == nil {
		headers["Referer"] = value
	}
}
