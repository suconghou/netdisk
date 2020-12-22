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

var client *baidudisk.Bclient

func init() {
	client = baidudisk.NewClient(config.Cfg.Token, config.Cfg.Root)
}

// Pwd print current path
func Pwd() error {
	return client.Pwd(config.Cfg.Path)
}

// GetInfo current backend info
func GetInfo() error {
	return client.Info()
}

// ListDir list files and dirs
func ListDir(filePath string, keep bool) error {
	if filePath == "" {
		filePath = config.Cfg.Path
	}
	err := client.Ls(filePath)
	if keep && err == nil && filePath != config.Cfg.Path {
		config.Cfg.Path = filePath
		config.Cfg.Save()
	}
	return err
}

// Get file form backend
func Get(filePath string, saveas string, transport *http.Transport) error {
	url := client.GetDownloadURL(filePath)
	return WgetURL(url, saveas, transport)
}

// WgetURL download a url file
func WgetURL(url string, saveas string, transport *http.Transport) error {
	file, fstart, err := utilgo.GetContinue(saveas)
	if err != nil {
		return err
	}
	defer file.Close()
	return fastloader.Load(file, map[string]int{url: 1}, 8, 1048576, fstart, 0, nil, transport, os.Stdout, nil)
}

// Play play a backend file
func Play(filePath string, saveas string, stdout bool, transport *http.Transport) error {
	url := client.GetDownloadURL(filePath)
	return PlayURL(url, saveas, stdout, transport)
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
	return fastloader.Load(file, map[string]int{url: 1}, 8, 1048576, fstart, 0, nil, transport, writer, hook)
}

// GetFileInfo print file info
func GetFileInfo(filePath string, dlink bool) error {
	return client.FileInfo(filePath, dlink)
}

// Put upload data to backend
func Put(savePath string, overwrite bool, file *os.File) error {
	return client.Put(savePath, overwrite, file)
}

// PutFile upload files to backend
func PutFile(filePath string, savePath string, fileSize uint64, ondup string) {

}

// PutFileRapid upload file
func PutFileRapid(filePath string, savePath string, fileSize uint64, ondup string, md5Str string) {

}

// Mkdir create dir
func Mkdir(path string) error {
	return client.Mkdir(path)
}

// DeleteFile delete files
func DeleteFile(fileName string) error {
	return client.Rm(fileName)
}

// MoveFile move file
func MoveFile(source string, target string) error {
	return client.Mv(source, target)
}

// CopyFile copy files
func CopyFile(source string, target string) error {
	return client.Cp(source, target)
}

// SearchFile search files
func SearchFile(fileName string) error {
	return client.Search(fileName)
}

// Empty clear
func Empty() error {
	return client.Clear()
}

// GetTaskList print task list
func GetTaskList() error {
	return client.TaskList()
}

// AddTask add a task
func AddTask(savePath string, sourceURL string) error {
	return client.TaskAdd(savePath, sourceURL)
}

// RemoveTask remove a task
func RemoveTask(id string) error {
	return client.TaskRemove(id)
}

// GetTaskInfo print one task info
func GetTaskInfo(ids string) error {
	return client.TaskInfo(ids)
}
