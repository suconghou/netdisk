package fslayer

import (
	"fmt"
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

func Pwd() error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Pwd(config.Cfg.Pcs.Path)
	}
	return nil
}

func GetInfo() error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Info()
	}
	return nil
}

func ListDir(filePath string, keep bool) error {
	if config.IsPcs() {
		if filePath == "" {
			filePath = config.Cfg.Pcs.Path
		}
		err := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Ls(filePath)
		if keep && err == nil && filePath != config.Cfg.Pcs.Path {
			config.Cfg.Pcs.Path = filePath
			config.Cfg.Save()
		}
		return err
	}
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
	var (
		startstr string
		thunkstr string
	)
	file, cstart, err := utilgo.GetContinue(saveas)
	defer file.Close()
	if err != nil {
		return err
	}
	if start == -1 {
		start = cstart
	}
	loader := fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(path.Base(file.Name())+" ", " ", nil, nil), nil)
	resp, total, filesize, thread, err := loader.Load(start, end)
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("该文件已经下载完毕")
		}
		return err
	}
	if start != 0 || end != 0 {
		startstr = fmt.Sprintf(",%d-%d", start, end)
	}
	if thread > 1 {
		thunkstr = fmt.Sprintf(",分块%dKB", thunk/1024)
	}
	util.Log.Printf("下载中,线程%d%s,大小%s/%s(%d/%d)%s", thread, thunkstr, utilgo.ByteFormat(uint64(total)), utilgo.ByteFormat(uint64(filesize)), total, filesize, startstr)
	n, err := io.Copy(file, resp)
	if err != nil {
		return err
	}
	util.Log.Printf("\n下载完毕,%d/%d", n, total)
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
		loader   *fastload.Fastloader
		file     io.WriteCloser
		err      error
		cstart   int64
		player   bool
		startstr string
		thunkstr string
	)
	if stdout {
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
			if loaded > 2 && !player {
				utilgo.CallPlayer(saveas)
				player = true
			}
		}
		loader = fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(path.Base(saveas)+" ", " ", hook, nil), nil)
	}
	defer file.Close()
	resp, total, filesize, thread, err := loader.Load(start, end)
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("该文件已经下载完毕")
		}
		return err
	}
	if start != 0 || end != 0 {
		startstr = fmt.Sprintf(",%d-%d", start, end)
	}
	if thread > 1 {
		thunkstr = fmt.Sprintf(",分块%dKB", thunk/1024)
	}
	util.Log.Printf("下载中,线程%d%s,大小%s/%s(%d/%d)%s", thread, thunkstr, utilgo.ByteFormat(uint64(total)), utilgo.ByteFormat(uint64(filesize)), total, filesize, startstr)
	n, err := io.Copy(file, resp)
	if err != nil {
		return err
	}
	util.Log.Printf("\n下载完毕,%d/%d", n, total)
	return nil
}

func GetFileInfo(filePath string, dlink bool) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Fileinfo(filePath, dlink)
	}
	return nil
}

func PutFile(filePath string, savePath string, fileSize uint64, ondup string) {

}

func PutFileRapid(filePath string, savePath string, fileSize uint64, ondup string, md5Str string) {

}

func Mkdir(path string) error {
	return nil
}

func DeleteFile(fileName string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Rm(fileName)
	}
	return nil
}

func MoveFile(source string, target string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Mv(source, target)
	}
	return nil
}

func CopyFile(source string, target string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Cp(source, target)
	}
	return nil
}

func SearchFile(fileName string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Search(fileName)
	}
	return nil
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
