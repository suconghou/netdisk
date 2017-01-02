package fslayer

import (
	"bytes"
	"config"
	"encoding/json"
	"fmt"
	"layers/netlayer"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"util"
)

type PcsInfo struct {
	Quota      int64
	Used       int64
	Request_id int
}

type FileItem struct {
	Fs_id int
	Path  string
	Ctime int64
	Mtime int64
	Size  int64
	IsDir int
}

type PcsFileList struct {
	List       []FileItem
	Request_id int
}

type PcsFileMeta struct {
	Fs_id       int
	Path        string
	Ctime       int64
	Mtime       int64
	Size        int64
	Isdir       int
	Ifhassubdir int
	Block_list  string
}

type PcsFileMetaList struct {
	List       []PcsFileMeta
	Request_id int
}

func GetInfo() {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/quota?method=%s&access_token=%s", "info", config.Cfg.Token)
	str := netlayer.Get(url)
	info := &PcsInfo{}
	err := json.Unmarshal([]byte(str), &info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		b := bytes.Buffer{}
		b.WriteString(util.DiskName(config.Cfg.Disk))
		b.WriteString("\n总大小:")
		b.WriteString(util.ByteFormat(info.Quota))
		b.WriteString("\n已使用:")
		b.WriteString(util.ByteFormat(info.Used))
		fmt.Println(b.String())
	}
}

//TODO order and limit
func ListDir(path string) bool {
	absPath := filepath.Join(config.Cfg.Root, path)
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "list", config.Cfg.Token, absPath)
	str := netlayer.Get(url)
	info := &PcsFileList{}
	err := json.Unmarshal([]byte(str), &info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		b := bytes.Buffer{}
		var total int64 = 0
		for _, item := range info.List {
			total = total + item.Size
			b.WriteString("\n")
			b.WriteString(util.StringPad(util.DateFormat(item.Ctime), 10))
			b.WriteString(util.StringPad(util.DateFormat(item.Mtime), 10))
			b.WriteString(util.StringPad(util.ByteFormat(item.Size), 10))
			b.WriteString(util.StringPad(util.BoolString(item.IsDir == 0, "f", "d")+"@"+strconv.Itoa(item.Fs_id), 20))
			b.WriteString(util.StringPad(item.Path, 20))
		}
		fmt.Print(util.DiskName(config.Cfg.Disk) + config.Cfg.Root + "  ➜  " + path + " " + util.ByteFormat(total))
		fmt.Println(b.String())
	}
	return true
}

func Get(filePath string, dist string, size int64, hash string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath)
	if dist == "" {
		dir, _ := os.Getwd()
		dist = filepath.Join(dir, filepath.Base(filePath))
	}
	netlayer.Download(url, dist, size, hash)
}

func Wget(filePath string, dist string, size int64, hash string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath)
	if dist == "" {
		dir, _ := os.Getwd()
		dist = filepath.Join(dir, filepath.Base(filePath))
	}
	netlayer.WgetDownload(url, dist, size, hash)
}

func GetPlay(filePath string, dist string, size int64, hash string) {
	var playType string = "M3U8_FLV_264_480"
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s&type=%s", "streaming", config.Cfg.Token, filePath, playType)
	if dist == "" {
		dir, _ := os.Getwd()
		dist = filepath.Join(dir, filepath.Base(filePath))
	}
	netlayer.DownloadPlay(url, dist, size, hash)
}

func GetPlayStream(filePath string, dist string, size int64, hash string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath)
	if dist == "" {
		dir, _ := os.Getwd()
		dist = filepath.Join(dir, filepath.Base(filePath))
	}
	netlayer.PlayStream(url, dist, size, hash)
}

func GetFileInfo(path string, noprint bool) (bool, int64, string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "meta", config.Cfg.Token, path)
	str := netlayer.Get(url)
	info := &PcsFileMetaList{}
	err := json.Unmarshal([]byte(str), &info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		if len(info.List) == 0 {
			fmt.Println(path + " 不存在")
			os.Exit(2)
		}
		item := info.List[0]
		b := bytes.Buffer{}
		b.WriteString(item.Path)
		b.WriteString("\n文件大小:" + util.ByteFormat(item.Size))
		b.WriteString("\n文件ID:" + strconv.Itoa(item.Fs_id))
		b.WriteString("\n创建时间:" + time.Unix(item.Ctime, 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n修改时间:" + time.Unix(item.Mtime, 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n类型:" + util.BoolString(item.Isdir == 0, "文件", "文件夹"))
		b.WriteString("\n子目录:" + util.BoolString(item.Ifhassubdir == 0, "无", "包含"))

		var Block_list []string
		if err := json.Unmarshal([]byte(item.Block_list), &Block_list); err == nil {
			if len(Block_list) == 1 {
				b.WriteString("\n文件哈希:" + Block_list[0])
			} else {
				for _, v := range Block_list {
					b.WriteString("\n哈希块:" + v)
				}
			}
		}
		if !noprint {
			fmt.Println(b.String())
		}
		return true, item.Size, Block_list[0]
	}
	return false, 0, ""
}

func Put() {
	fmt.Println("put")

}
