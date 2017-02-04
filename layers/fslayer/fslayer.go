package fslayer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/suconghou/fastload/fastload"
	"io"
	"mime/multipart"
	"netdisk/config"
	"netdisk/layers/netlayer"
	"netdisk/util"
	"os"
	"path"
	"strconv"
	"time"
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
	Task_id        uint64
	Request_id     uint64
	Rapid_download uint8
}

type TaskItem struct {
	Task_id     string
	Od_type     string
	Source_url  string
	Save_path   string
	Rate_limit  string
	Timeout     string
	Status      string
	Create_time string
	Task_name   string
}

type TaskDetailFile struct {
	File_name string
	File_size string
}

type TaskDetail struct {
	Status        string
	File_size     string
	Finished_size string
	Create_time   string
	Start_time    string
	Finish_time   string
	Save_path     string
	Source_url    string
	Task_name     string
	Od_type       string
	File_list     []TaskDetailFile
	Rresult       int
}

type TaskDetailList struct {
	Task_info  map[string]TaskDetail
	Request_id uint64
}

type TaskList struct {
	Task_info  []TaskItem
	Total      uint32
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
	size := fastload.GetUrlInfo(url)
	if size > 1 {
		fmt.Println(dist + " " + util.ByteFormat(size))
		var hash string = ""
		netlayer.WgetDownload(url, dist, size, hash)
	} else {
		fmt.Println("远程文件大小未知")
	}
}

func PlayUrl(url string, dist string, stdout bool) {
	size := fastload.GetUrlInfo(url)
	if size > 1 {
		fmt.Println(dist + " " + util.ByteFormat(size))
		var hash string = ""
		netlayer.PlayStream(url, dist, size, hash, stdout)
	}
}

func GetPlayStream(filePath string, dist string, size uint64, hash string, stdout bool) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath)
	netlayer.PlayStream(url, dist, size, hash, stdout)
}

func GetFileInfo(filePath string, noprint bool) (bool, uint64, string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "meta", config.Cfg.Token, filePath)
	str := netlayer.Get(url)
	info := &PcsFileMetaList{}
	err := json.Unmarshal([]byte(str), &info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		if len(info.List) == 0 {
			os.Stderr.Write([]byte(filePath + " 不存在\n"))
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
				panic(err)
			}
		}
		if !noprint {
			fmt.Println(b.String())
			if len(os.Args) == 4 && os.Args[3] == "-i" {
				fmt.Println(fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s", "download", config.Cfg.Token, filePath))
			}
		}
		return true, item.Size, md5
	}
	return false, 0, ""
}

func PutFile(filePath string, savePath string, fileSize uint64, ondup string) {
	startTime := time.Now()

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
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening file")
		panic(err)
	}
	//iocopy
	counter := &netlayer.PutWriteCounter{}
	counter.Size = fileSize

	reader := io.TeeReader(file, fileWriter)
	_, err = io.Copy(counter, reader)
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
	} else {
		endTime := time.Since(startTime).Seconds()
		speed := float64(fileSize/1024) / endTime
		b := bytes.Buffer{}
		b.WriteString("\n已保存为:" + info.Path)
		b.WriteString("\n文件大小:" + util.ByteFormat(info.Size))
		b.WriteString("\n文件标识:" + strconv.FormatUint(info.Fs_id, 10))
		b.WriteString("\n创建时间:" + time.Unix(int64(info.Ctime), 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n修改时间:" + time.Unix(int64(info.Mtime), 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n文件哈希:" + info.Md5)
		b.WriteString("\n耗时:" + fmt.Sprintf("%.1fs,速度:%.2fKB/s", endTime, speed))
		fmt.Println(b.String())
	}

}

func PutFileRapid(filePath string, savePath string, fileSize uint64, ondup string, md5Str string) {
	startTime := time.Now()
	contentCrc32, sliceMd5 := util.GetCrc32AndMd5(filePath)
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s&content-length=%d&content-md5=%s&slice-md5=%s&content-crc32=%s&ondup=%s", "rapidupload", config.Cfg.Token, savePath, fileSize, md5Str, sliceMd5, contentCrc32, ondup)
	str := netlayer.PostWait(url, "application/x-www-form-urlencoded", nil)
	info := &UpFileInfo{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		panic(err)
	} else if info.Fs_id > 1 {
		endTime := time.Since(startTime).Seconds()
		speed := float64(fileSize/1024) / endTime
		b := bytes.Buffer{}
		b.WriteString("已保存为:" + info.Path)
		b.WriteString("\n文件大小:" + util.ByteFormat(info.Size))
		b.WriteString("\n文件标识:" + strconv.FormatUint(info.Fs_id, 10))
		b.WriteString("\n创建时间:" + time.Unix(int64(info.Ctime), 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n修改时间:" + time.Unix(int64(info.Mtime), 0).Format("2006/01/02 15:04:05"))
		b.WriteString("\n文件哈希:" + info.Md5)
		b.WriteString("\n耗时:" + fmt.Sprintf("%.1fs,秒传,速度:%.2fKB/s", endTime, speed))
		fmt.Println(b.String())
	} else {
		PutFile(filePath, savePath, fileSize, ondup)
	}
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

func CopyFile(source string, target string) {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&from=%s&to=%s", "copy", config.Cfg.Token, source, target)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &PcsRequest{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		panic(err)
	}
	fmt.Println(util.BoolString(info.Request_id > 1, "已复制", "未复制"))
}

func SearchFile(fileName string) {
	var re int = 1
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&path=%s&wd=%s&re=%s", "search", config.Cfg.Token, config.Cfg.Root, fileName, re)
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
		fmt.Print(util.DiskName(config.Cfg.Disk) + config.Cfg.Root + "  ➜  " + fileName + " " + util.ByteFormat(total))
		fmt.Println(b.String())
	}

}

func Empty() {
	url := fmt.Sprintf("https://pcs.baidu.com/rest/2.0/pcs/file?method=%s&access_token=%s&type=%s", "delete", config.Cfg.Token, "recycle")
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &PcsRequest{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		panic(err)
	}
	fmt.Println(util.BoolString(info.Request_id > 1, "已清空回收站", "未清空回收站"))
}

func GetTaskList() {
	url := fmt.Sprintf("https://pan.baidu.com/rest/2.0/services/cloud_dl?method=%s&access_token=%s&status=255&app_id=250528&need_task_info=1", "list_task", config.Cfg.Token)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &TaskList{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		os.Stdout.Write(str)
		panic(err)
	} else {
		b := bytes.Buffer{}
		for _, item := range info.Task_info {
			create_time, _ := strconv.Atoi(item.Create_time)
			b.WriteString(fmt.Sprintf("\n%s@%s %s", item.Task_id, item.Task_name, util.DateFormat(uint64(create_time))))
			b.WriteString(fmt.Sprintf("\n%s ➜ %s\n", item.Source_url, item.Save_path))
		}
		fmt.Print(util.DiskName(config.Cfg.Disk) + config.Cfg.Root + fmt.Sprintf("  ➜  离线任务: %d个任务", info.Total))
		fmt.Println(b.String())
	}
}

func AddTask(savePath string, sourceUrl string) {
	url := fmt.Sprintf("https://pan.baidu.com/rest/2.0/services/cloud_dl?method=%s&access_token=%s&save_path=%s&source_url=%s&app_id=250528", "add_task", config.Cfg.Token, savePath, sourceUrl)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &AddTaskRet{}
	if err := json.Unmarshal(str, &info); err == nil {
		fmt.Println("任务ID:" + strconv.FormatUint(info.Task_id, 10))
		if info.Rapid_download == 1 {
			fmt.Println("离线已秒杀:" + savePath)
			GetFileInfo(savePath, false)
		}
	} else {
		fmt.Println(err)
	}
}

func RemoveTask(id string) {

}

func GetTaskInfo(ids string) {
	url := fmt.Sprintf("https://pan.baidu.com/rest/2.0/services/cloud_dl?method=%s&access_token=%s&task_ids=%s&app_id=250528", "query_task", config.Cfg.Token, ids)
	str := netlayer.Post(url, "application/x-www-form-urlencoded", nil)
	info := &TaskDetailList{}
	err := json.Unmarshal(str, &info)
	if err != nil {
		os.Stdout.Write(str)
		panic(err)
	} else {
		b := bytes.Buffer{}
		for id, item := range info.Task_info {
			create_time, _ := strconv.Atoi(item.Create_time)
			start_time, _ := strconv.Atoi(item.Start_time)
			file_size, _ := strconv.Atoi(item.File_size)
			finished_size, _ := strconv.Atoi(item.Finished_size)
			b.WriteString(fmt.Sprintf("文件:%s@%s %s\n", id, item.Task_name, util.DateFormat(uint64(create_time))))
			b.WriteString(fmt.Sprintf("大小:%s\n", util.ByteFormat(uint64(file_size))))
			b.WriteString(fmt.Sprintf("开始下载时间:%s\n", util.DateFormat(uint64(start_time))))
			b.WriteString(fmt.Sprintf("已下载:%s\n", util.ByteFormat(uint64(finished_size))))
			b.WriteString(fmt.Sprintf("原地址:%s\n", item.Source_url))
			b.WriteString(fmt.Sprintf("存储为:%s\n", item.Save_path))
		}
		fmt.Print(util.DiskName(config.Cfg.Disk) + config.Cfg.Root + "  ➜  任务详情: \n")
		fmt.Println(b.String())
	}

}
