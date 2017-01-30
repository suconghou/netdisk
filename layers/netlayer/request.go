package netlayer

import (
	"fmt"
	"github.com/suconghou/fastload/fastload"
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

func Get(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("%s",err)))
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
		os.Stderr.Write([]byte(fmt.Sprintf("%s",err)))
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
		os.Stderr.Write([]byte(fmt.Sprintf("%s",err)))
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

// WriteCounter counts the number of bytes written to it.
type WriteCounter struct {
	Total uint64 // Total # of bytes written
	Size  uint64
	Part  string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	var per float64 = float64(wc.Total) / float64(wc.Size)
	var i int = int(per * 100)
	fmt.Printf("\r%s%d%% %s %s", util.Bar(i, 25), i, util.ByteFormat(wc.Total), wc.Part)
	return n, nil
}

// WriteCounter counts the number of bytes written to it.
type PutWriteCounter struct {
	Total uint64 // Total # of bytes written
	Size  uint64
	Part  string
}

func (wc *PutWriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	var per float64 = float64(wc.Total) / float64(wc.Size)
	var i int = int(per * 100)
	fmt.Printf("\r%s%d%% %s %s  ", util.Bar(i, 25), i, util.ByteFormat(wc.Total), wc.Part)
	return n, nil
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
	counter := &WriteCounter{}
	counter.Size = size
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

func WgetDownload(url string, saveas string, size uint64, hash string) {
	var thread uint8 = 8
	var thunk uint32 = 1048576 * 2
	start := fastload.GetContinue(saveas)
	end := size
	fmt.Printf("下载中...线程%d,分块大小%dKB\n", thread, thunk/1024)
	startTime := time.Now()
	if start >= size {
		fmt.Println("\n已下载完毕,校验MD5中...")
		util.PrintMd5(saveas)
		os.Exit(0)
	} else {
		fastload.Load(url, saveas, start, end, thread, thunk, false, nil)
	}
	endTime := time.Since(startTime)
	speed := float64((size-start)/1024) / endTime.Seconds()
	fmt.Printf("\n下载完毕,耗时%s,%.2fKB/s,校验MD5中...\n", endTime.String(), speed)
	util.PrintMd5(saveas)
}

func PlayStream(url string, saveas string, size uint64, hash string, stdout bool) {
	var thread uint8 = 4
	var thunk uint32 = 524288 * 4
	var startContinue uint64 = 0
	if len(os.Args) >= 4 {
		var fl string
		if len(os.Args) > 4 {
			fl = os.Args[4]
		} else {
			fl = os.Args[3]
		}
		match, matchStart, matchEnd := getRange(fl, startContinue, size)
		if match {
			size = matchEnd
			startContinue = matchStart
		}
	}
	var start uint64 = 0
	var end uint64 = size
	if startContinue > 1 {
		start = startContinue
	} else {
		start = fastload.GetContinue(saveas)
	}
	if !stdout {
		fmt.Printf("下载中...线程%d,分块大小%dKB\n", thread, thunk/1024)
	}
	startTime := time.Now()
	if start >= size {
		fmt.Printf("\n已下载完毕,校验MD5中...\n")
		util.PrintMd5(saveas)
		os.Exit(0)
	} else {
		var playerRun bool = false
		fastload.Load(url, saveas, start, end, thread, thunk, stdout, func(percent int, downloaded uint64) {
			if percent > 5 && !playerRun {
				playerRun = true
				callPlayer(saveas)
			}
		})
	}
	if !stdout {
		endTime := time.Since(startTime)
		speed := float64((size-start)/1024) / endTime.Seconds()
		fmt.Printf("\n下载完毕,耗时%s,%.2fKB/s %s \n", endTime.String(), speed, util.BoolString(stdout, "", ",校验MD5中..."))
		util.PrintMd5(saveas)
	}
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

func getRange(str string, start uint64, end uint64) (bool, uint64, uint64) {
	var rangeXp = regexp.MustCompile(`^(\d+)-(\d+)$`)
	var rangeXpAll = regexp.MustCompile(`^(\d+)-$`)
	var matched bool = false
	var matchStart uint64
	var matchEnd uint64
	if rangeXp.MatchString(str) {
		match := rangeXp.FindStringSubmatch(str)
		strSt1, _ := strconv.Atoi(match[1])
		matchStart = uint64(strSt1)
		strSt2, _ := strconv.Atoi(match[2])
		matchEnd = uint64(strSt2)
		if matchStart < 8192 {
			matchStart = matchStart * 1048576
		}
		if matchEnd < 8192 {
			matchEnd = matchEnd * 1048576
		}
		if (matchEnd < matchStart) || (matchStart >= end) {
			fmt.Println("error range")
			os.Exit(1)
		} else if matchEnd > end {
			matchEnd = end
		}
		matched = true

	} else if rangeXpAll.MatchString(str) {
		match := rangeXpAll.FindStringSubmatch(str)
		strSt1, _ := strconv.Atoi(match[1])
		matchStart = uint64(strSt1)
		if matchStart < 8192 {
			matchStart = matchStart * 1048576
		}
		if matchStart > end {
			fmt.Println("error range")
			os.Exit(1)
		} else {
			matchEnd = end
		}
		matched = true
	}
	return matched, matchStart, matchEnd

}
