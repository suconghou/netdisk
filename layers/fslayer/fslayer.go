package fslayer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/netdisk/config"
	"github.com/suconghou/netdisk/layers/baidudisk"
	"github.com/suconghou/netdisk/util"
	"github.com/suconghou/utilgo"
)

// Pwd print current path
func Pwd() error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Pwd(config.Cfg.Pcs.Path)
	}
	return nil
}

// GetInfo current backend info
func GetInfo() error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Info()
	}
	return nil
}

// ListDir list files and dirs
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
		return WgetURL(url, saveas, reqHeader, thread, thunk, start, end, nil)
	}
	// 腾讯云 cos
	return nil
}

// WgetURL download a url file
func WgetURL(url string, saveas string, reqHeader http.Header, thread int32, thunk int64, start int64, end int64, proxy *http.Transport) error {
	var (
		startstr    string
		thunkstr    string
		showsizestr string
	)
	file, cstart, err := utilgo.GetContinue(saveas)
	defer file.Close()
	if err != nil {
		return err
	}
	if start == -1 {
		start = cstart
	}
	loader := fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(filepath.Base(file.Name())+" ", " ", nil, nil), proxy, nil)
	body, _, total, filesize, thread, err := loader.Load(start, end, nil)
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
	if total > 0 && filesize > 0 {
		showsizestr = fmt.Sprintf(",大小%s/%s(%d/%d)", utilgo.ByteFormat(uint64(total)), utilgo.ByteFormat(uint64(filesize)), total, filesize)
	}
	util.Log.Printf("下载中,线程%d%s%s%s", thread, thunkstr, showsizestr, startstr)
	n, err := io.Copy(file, body)
	if err != nil {
		return err
	}
	util.Log.Printf("\n下载完毕,%d%s", n, utilgo.BoolString(total > 0, fmt.Sprintf("/%d", total), ""))
	return nil
}

// Play play a backend file
func Play(filePath string, saveas string, reqHeader http.Header, thread int32, thunk int64, start int64, end int64, stdout bool) error {
	if config.IsPcs() {
		url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(filePath)
		return PlayURL(url, saveas, reqHeader, thread, thunk, start, end, stdout, nil)
	}
	return nil
}

// PlayURL play a url media
func PlayURL(url string, saveas string, reqHeader http.Header, thread int32, thunk int64, start int64, end int64, stdout bool, proxy *http.Transport) error {
	var (
		loader      *fastload.Fastloader
		file        io.WriteCloser
		err         error
		cstart      int64
		player      bool
		startstr    string
		thunkstr    string
		showsizestr string
	)
	if stdout {
		if start == -1 {
			start = utilgo.HasFileSize(saveas)
		}
		file = os.Stdout
		loader = fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(filepath.Base(saveas)+" ", " ", nil, os.Stderr), proxy, nil)
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
		loader = fastload.NewLoader(url, thread, thunk, reqHeader, utilgo.ProgressBar(filepath.Base(saveas)+" ", " ", hook, nil), proxy, nil)
	}
	defer file.Close()
	body, _, total, filesize, thread, err := loader.Load(start, end, nil)
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
	if total > 0 && filesize > 0 {
		showsizestr = fmt.Sprintf(",大小%s/%s(%d/%d)", utilgo.ByteFormat(uint64(total)), utilgo.ByteFormat(uint64(filesize)), total, filesize)
	}
	util.Log.Printf("下载中,线程%d%s%s%s", thread, thunkstr, showsizestr, startstr)
	n, err := io.Copy(file, body)
	if err != nil {
		return err
	}
	util.Log.Printf("\n下载完毕,%d%s", n, utilgo.BoolString(total > 0, fmt.Sprintf("/%d", total), ""))
	return nil
}

// GetFileInfo print file info
func GetFileInfo(filePath string, dlink bool) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).FileInfo(filePath, dlink)
	}
	return nil
}

// Put upload data to backend
func Put(savePath string, overwrite bool, file *os.File) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Put(savePath, overwrite, file)
	}
	return nil
}

// PutFile upload files to backend
func PutFile(filePath string, savePath string, fileSize uint64, ondup string) {

}

// PutFileRapid upload file
func PutFileRapid(filePath string, savePath string, fileSize uint64, ondup string, md5Str string) {

}

// Mkdir create dir
func Mkdir(path string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Mkdir(path)
	}
	return nil
}

// DeleteFile delete files
func DeleteFile(fileName string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Rm(fileName)
	}
	return nil
}

// MoveFile move file
func MoveFile(source string, target string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Mv(source, target)
	}
	return nil
}

// CopyFile copy files
func CopyFile(source string, target string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Cp(source, target)
	}
	return nil
}

// SearchFile search files
func SearchFile(fileName string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Search(fileName)
	}
	return nil
}

// Empty clear
func Empty() error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).Clear()
	}
	return nil
}

// GetTaskList print task list
func GetTaskList() error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).TaskList()
	}
	return nil
}

// AddTask add a task
func AddTask(savePath string, sourceURL string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).TaskAdd(savePath, sourceURL)
	}
	return nil
}

// RemoveTask remove a task
func RemoveTask(id string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).TaskRemove(id)
	}
	return nil
}

// GetTaskInfo print one task info
func GetTaskInfo(ids string) error {
	if config.IsPcs() {
		return baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).TaskInfo(ids)
	}
	return nil
}
