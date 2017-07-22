package fslayer

import (
	"fmt"
	"io"
	"net/http"
	"netdisk/config"
	"netdisk/layers/baidudisk"
	"netdisk/layers/netlayer"
	"netdisk/util"
	"os"
	"path"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/utilgo"
)

func GetInfo() {

}

func ListDir(filePath string) {

}

func Get(filePath string) {
	if config.Cfg.Driver == "pcs" {

		resp, err := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Get(filePath)

	}

}

// GetURL download a url file
func GetURL(url string) error {
	saveas, err := utilgo.GetStorePath(url)
	if err != nil {
		return err
	}
	file, start, err := utilgo.GetContinue(saveas)
	if err != nil {
		return err
	}
	resp, total, thread, err := fastload.Get(url, 0, 0, utilgo.ProgressBar(path.Base(file.Name())+" ", " "))
	if err != nil {
		return err
	}
	n, err := io.Copy(file, resp)
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("\r\n%s : download ok use thread %d size %d/%d \r\n", url, thread, n, total))
	file.Close()
	return nil
}

// Wget file form baidu disk
func Wget(filePath string) error {
	if true { // 百度网盘
		url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Pcs.Token, filePath)
		return WgetURL(url, filePath)
	}
	// 腾讯云 cos
	return nil
}

// WgetURL download a url file
func WgetURL(url string, saveas string) error {
	// netlayer.ParseCookieUaRefer()
	file, start, err := utilgo.GetContinue(saveas)
	if err != nil {
		return err
	}
	reqHeader := http.Header{}
	reqHeader.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	loader := fastload.NewLoader(url, 8, 1048576, reqHeader, utilgo.ProgressBar(path.Base(filePath)+" ", " "), nil)

	resp, total, thread, err := fastload.Get(url, 0, 0, utilgo.ProgressBar(path.Base(file.Name())+" ", " "))
	if err != nil {
		return err
	}
	n, err := io.Copy(file, resp)
	if err != nil {
		return err
	}
	os.Stdout.WriteString(fmt.Sprintf("\r\n%s : download ok use thread %d size %d/%d \r\n", url, thread, n, total))
	file.Close()
	return nil
}

// PlayURL play a url media
func PlayURL(url string, dist string, stdout bool) {
	netlayer.ParseCookieUaRefer()
	size, rangeAble, err := fastload.GetUrlInfo(url, util.HasFlag("--GET"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var hash string
	var msg string
	if size > 0 {
		msg = dist + " " + util.ByteFormat(size)
	} else {
		msg = dist + " 远程文件大小未知"
	}
	if stdout {
		os.Stderr.Write([]byte(msg))
	} else {
		fmt.Println(msg)
	}
	netlayer.PlayStream(url, dist, size, hash, stdout, rangeAble)
}

func GetPlayStream(filePath string, dist string, size uint64, hash string, stdout bool) {

}

func GetFileInfo(filePath string, noprint bool) (bool, uint64, string) {

}

func PutFile(filePath string, savePath string, fileSize uint64, ondup string) {

}

func PutFileRapid(filePath string, savePath string, fileSize uint64, ondup string, md5Str string) {

}

func Mkdir(path string) {

}

func DeleteFile(path string) {

}

func MoveFile(source string, target string) {

}

func CopyFile(source string, target string) {

}

func SearchFile(fileName string) {

}

func Empty() {
}

func GetTaskList() {

}

func AddTask(savePath string, sourceUrl string) {

}

func RemoveTask(id string) {

}

func GetTaskInfo(ids string) {

}

func showTaskStatus(status string) string {
	if v, ok := taskStatusMap[status]; ok {
		return v
	} else {
		return "未知状态"
	}
}
