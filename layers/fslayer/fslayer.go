package fslayer

import (
	"io"
	"net/http"
	"os"

	"github.com/suconghou/fastload/fastloader"
	"github.com/suconghou/netdisk/config"
	"github.com/suconghou/netdisk/layers/baidudisk"
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
func Get(filePath string, saveas string, transport *http.Transport) error {
	if config.IsPcs() {
		url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(filePath)
		return WgetURL(url, saveas, transport)
	}
	// 腾讯云 cos
	return nil
}

// WgetURL download a url file
func WgetURL(url string, saveas string, transport *http.Transport) error {
	file, fstart, err := utilgo.GetContinue(saveas)
	if err != nil {
		return err
	}
	defer file.Close()
	return fastloader.Load(file, url, fstart, transport, os.Stdout, nil)
}

// Play play a backend file
func Play(filePath string, saveas string, stdout bool, transport *http.Transport) error {
	if config.IsPcs() {
		url := baidudisk.NewClient(config.Cfg.Pcs.Token, config.Cfg.Pcs.Root).GetDownloadURL(filePath)
		return PlayURL(url, saveas, stdout, transport)
	}
	return nil
}

// PlayURL play a url media
func PlayURL(url string, saveas string, stdout bool, transport *http.Transport) error {
	var (
		file   *os.File
		err    error
		fstart int64
		player bool
		writer io.Writer = os.Stderr
		hook   func(loaded float64, speed float64, remain float64)
	)
	if stdout {
		file = os.Stdout
	} else {
		file, fstart, err = utilgo.GetContinue(saveas)
		if err != nil {
			return err
		}
		hook = func(loaded float64, speed float64, remain float64) {
			if loaded > 2 && !player {
				utilgo.CallPlayer(saveas)
				player = true
			}
		}
	}
	defer file.Close()
	return fastloader.Load(file, url, fstart, transport, writer, hook)
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
