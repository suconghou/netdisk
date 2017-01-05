package fslayer

import (
	"bytes"
	"config"
	"encoding/json"
	"fmt"
	"io"
	"layers/netlayer"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"
	"util"
)

type PcsRequest struct {
	Request_id uint64
}

type PcsRequestError struct {
	Request_id uint64
	Error_code int
	Error_msg  string
}

type PcsInfo struct {
	Quota      uint64
	Used       uint64
	Request_id uint64
}

type FileItem struct {
	Fs_id uint64
	Path  string
	Ctime uint64
	Mtime uint64
	Size  uint64
	IsDir int
}

type PcsFileList struct {
	List       []FileItem
	Request_id uint64
}

type PcsFileMeta struct {
	Fs_id       uint64
	Path        string
	Ctime       uint64
	Mtime       uint64
	Size        uint64
	Isdir       int
	Ifhassubdir int
	Block_list  string
	Filenum     int
}

type PcsFileMetaList struct {
	List       []PcsFileMeta
	Request_id uint64
}

type UpFileInfo struct {
	Fs_id      uint64
	Path       string
	Ctime      uint64
	Mtime      uint64
	Md5        string
	Size       uint64
	Request_id uint64
}

type AddTaskRet struct {
	Task_id    uint64
	Request_id uint64
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
func ListDir(filePath string) bool {
	absPath := path.Join(config.Cfg.Root, filePath)
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "list", config.Cfg.Token, absPath)
	str := netlayer.Get(url)
	info := &PcsFileList{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		b := bytes.Buffer{}
		var total uint64 = 0
		for _, item := range info.List {
			total = total + item.Size
			b.WriteString("\n")
			b.WriteString(util.StringPad(util.DateFormat(item.Ctime), 10))
			b.WriteString(util.StringPad(util.DateFormat(item.Mtime), 10))
			b.WriteString(util.StringPad(util.ByteFormat(item.Size), 10))
			b.WriteString(util.StringPad(fmt.Sprintf("%s@%d", util.BoolString(item.IsDir == 0, "f", "d"), item.Fs_id), 20))
			b.WriteString(util.StringPad(item.Path, 20))
		}
		fmt.Print(util.DiskName(config.Cfg.Disk) + config.Cfg.Root + "  ➜  " + filePath + " " + util.ByteFormat(total))
		fmt.Println(b.String())
	}
	return true
}

func Get(filePath string, dist string, size uint64, hash string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath)
	netlayer.Download(url, dist, size, hash)
}

func Wget(filePath string, dist string, size uint64, hash string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath)
	netlayer.WgetDownload(url, dist, size, hash)
}

func WgetUrl(url string, dist string) {
	size := netlayer.GetUrlInfo(url)
	if size > 1 {
		fmt.Println(dist + " " + util.ByteFormat(size))
		var hash string = ""
		netlayer.WgetDownload(url, dist, size, hash)
	}
}

func PlayUrl(url string, dist string) {
	size := netlayer.GetUrlInfo(url)
	if size > 1 {
		fmt.Println(dist + " " + util.ByteFormat(size))
		var hash string = ""
		netlayer.PlayStream(url, dist, size, hash)
	}
}

func GetPlayStream(filePath string, dist string, size uint64, hash string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath)
	netlayer.PlayStream(url, dist, size, hash)
}

func GetFileInfo(path string, noprint bool) (bool, uint64, string) {
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
		b.WriteString("\n文件标识:" + strconv.FormatUint(item.Fs_id, 10))
		b.WriteString("\n创建时间:" + time.Unix(int64(item.Ctime), 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n修改时间:" + time.Unix(int64(item.Mtime), 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n类型:" + util.BoolString(item.Isdir == 0, "文件", "文件夹"))
		b.WriteString("\n子目录:" + util.BoolString(item.Ifhassubdir == 0, "无", "包含"))
		var Block_list []string
		var md5 string = " "
		if item.Block_list == "" {
			b.WriteString("\n包含文件:" + strconv.Itoa(item.Filenum))
		} else {

			if err := json.Unmarshal([]byte(item.Block_list), &Block_list); err == nil {
				if len(Block_list) == 1 {
					md5 = Block_list[0]
					b.WriteString("\n文件哈希:" + Block_list[0])
				} else if len(Block_list) > 1 {
					for _, v := range Block_list {
						b.WriteString("\n哈希块:" + v)
					}
				}
			} else {
				fmt.Println(err)
			}

		}

		if !noprint {
			fmt.Println(b.String())
		}
		return true, item.Size, md5
	}
	return false, 0, ""
}

func PutFile(filePath string, savePath string, ondup string) {
	url := fmt.Sprintf("https://c.pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s&ondup=%s", "upload", config.Cfg.Token, savePath, ondup)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Println("error writing to buffer")
		panic(err)
	}
	//打开文件句柄操作
	fh, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening file")
		panic(err)
	}
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		panic(err)
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	str := netlayer.Post(url, contentType, bodyBuf)
	info := &UpFileInfo{}
	err = json.Unmarshal(str, &info)
	if err != nil {
		panic(err)
	}
	b := bytes.Buffer{}
	b.WriteString("已保存为:" + info.Path)
	b.WriteString("\n文件大小:" + util.ByteFormat(info.Size))
	b.WriteString("\n文件标识:" + strconv.FormatUint(info.Fs_id, 10))
	b.WriteString("\n创建时间:" + time.Unix(int64(info.Ctime), 0).Format("2006/01/02 15:04:05"))
	b.WriteString("\n修改时间:" + time.Unix(int64(info.Mtime), 0).Format("2006/01/02 15:04:05"))
	b.WriteString("\n文件哈希:" + info.Md5)
	fmt.Println(b.String())
}

func Mkdir(path string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "mkdir", config.Cfg.Token, path)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	infoerr := &PcsRequestError{}
	if err := json.Unmarshal(str, &infoerr); err == nil {
		if infoerr.Error_msg == "" {
			fmt.Println("已创建 " + path)
		} else {
			fmt.Println(path + " " + infoerr.Error_msg)
		}
	} else {
		panic(err)
	}
}

func Put() {
	fmt.Println("put")
}

func DeleteFile(path string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "delete", config.Cfg.Token, path)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &PcsRequest{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		panic(err)
	}
	fmt.Println(util.BoolString(info.Request_id > 1, "已删除", "未删除"))
}

func MoveFile(source string, target string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&from=%s&to=%s", "move", config.Cfg.Token, source, target)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &PcsRequest{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		panic(err)
	}
	fmt.Println(util.BoolString(info.Request_id > 1, "已移动", "未移动"))
}

func GetTaskList() {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl?method=%s&access_token=%s", "list_task", config.Cfg.Token)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	os.Stdout.Write(str)
}

func AddTask(savePath string, sourceUrl string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/services/cloud_dl?method=%s&access_token=%s&save_path=%s&source_url=%s", "add_task", config.Cfg.Token, savePath, sourceUrl)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &AddTaskRet{}
	if err := json.Unmarshal(str, &info); err == nil {
		fmt.Println("任务ID:" + strconv.FormatUint(info.Task_id, 10))
	} else {
		fmt.Println(err)
	}
}

func RemoveTask(id string) {

}

func GetTaskInfo(ids string) {

}
