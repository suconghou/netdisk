package fslayer

import (
	"io"
	"net/http"
	"netdisk/config"
	"netdisk/layers/baidudisk"
	"netdisk/util"
	"path"

	"os"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/utilgo"
)

func GetInfo() {

}

func ListDir(filePath string) error {
	return nil
}

// Get file form backend
func Get(filePath string, saveas string, reqHeader http.Header, thread int32, thunk int64, start int64, end int64) error {
	if config.IsPcs() {
		url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(filePath)
		return WgetURL(url, saveas, reqHeader, thread, thunk, start, end)
	}
	// 腾讯云 cos
	return nil
}

// WgetURL download a url file
func WgetURL(url string, saveas string, reqHeader http.Header, thread int32, thunk int64, start int64, end int64) error {
	file, cstart, err := utilgo.GetContinue(saveas)
	if err != nil {
		return err
	}
	if start == -1 {
		start = cstart
	}
	loader := fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(path.Base(file.Name())+" ", " ", nil, nil), nil)
	resp, total, thread, err := loader.Load(start, end)
	if err != nil {
		return err
	}
	util.Log.Printf("下载中,线程%d大小%d", thread, total)
	n, err := io.Copy(file, resp)
	if err != nil {
		return err
	}
	util.Log.Printf("下载完毕,%d/%d", n, total)
	file.Close()
	return nil
}

// Play play a backend file
func Play(filePath string, saveas string, reqHeader http.Header, thread int32, thunk int64, start int64, end int64, stdout bool) error {
	if config.IsPcs() {
		url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(filePath)
		return PlayURL(url, saveas, reqHeader, thread, thunk, start, end, stdout)
	}
	return nil
}

// PlayURL play a url media
func PlayURL(url string, saveas string, reqHeader http.Header, thread int32, thunk int64, start int64, end int64, stdout bool) error {
	var (
		loader *fastload.Fastloader
		file   io.WriteCloser
		err    error
		cstart int64
	)
	if stdout {
		util.Log.SetOutput(os.Stderr)
		if start == -1 {
			start = 0
		}
		file = os.Stdout
		loader = fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(path.Base(saveas)+" ", " ", nil, os.Stderr), nil)
	} else {
		file, cstart, err = utilgo.GetContinue(saveas)
		if err != nil {
			return err
		}
		if start == -1 {
			start = cstart
		}
		hook := func(loaded float64, speed float64, remain float64) {
			if loaded > 5 {
				utilgo.CallPlayer(saveas)
			}
		}
		loader = fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(path.Base(saveas)+" ", " ", hook, nil), nil)
	}
	resp, total, thread, err := loader.Load(start, end)
	if err != nil {
		return err
	}
	util.Log.Printf("下载中,线程%d,大小%d", thread, total)
	n, err := io.Copy(file, resp)
	if err != nil {
		return err
	}
	util.Log.Printf("下载完毕,%d/%d", n, total)
	file.Close()
	return nil
}

func GetFileInfo(filePath string, noprint bool) {

}

func PutFile(filePath string, savePath string, fileSize uint64, ondup string) {

}

func PutFileRapid(filePath string, savePath string, fileSize uint64, ondup string, md5Str string) {

}

func Mkdir(path string) error {
	return nil
}

func DeleteFile(path string) error {
	return nil
}

func MoveFile(source string, target string) error {
	return nil
}

func CopyFile(source string, target string) error {
	return nil
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
