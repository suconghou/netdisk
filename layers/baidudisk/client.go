package baidudisk

import (
	"bytes"
	"fmt"
	"io"
	"path"

	"github.com/bitly/go-simplejson"
	"github.com/suconghou/utilgo"
)

const name = "百度网盘"

// Bclient is a baidudisk client
type Bclient struct {
	token     string
	root      string
	path      string
	apiURL    string
	uploadURL string
	taskURL   string
	infoURL   string
}

var taskStatusMap = map[string]string{
	"0": "下载成功",
	"1": "下载进行中",
	"2": "系统错误",
	"3": "资源不存在",
	"4": "下载超时",
	"5": "资源存在但下载失败",
	"6": "存储空间不足",
	"7": "任务已取消",
}

// NewClient return a client
func NewClient(token string, root string) *Bclient {
	return &Bclient{
		token:     token,
		root:      root,
		apiURL:    "https://pcs.baidu.com/rest/2.0/pcs/file",
		infoURL:   "https://pcs.baidu.com/rest/2.0/pcs/quota",
		uploadURL: "https://c.pcs.baidu.com/rest/2.0/pcs/file",
		taskURL:   "https://pan.baidu.com/rest/2.0/services/cloud_dl",
	}
}

// Pwd print current dir
func (bc *Bclient) Pwd(p string) error {
	fmt.Println(name + bc.root + "  ➜  " + p)
	return nil
}

// Ls print dir content in cli
func (bc *Bclient) Ls(p string) error {
	js, err := bc.APILs(p)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	var total uint64
	for i := range js.Get("list").MustArray() {
		item := js.Get("list").GetIndex(i)
		size := item.Get("size").MustUint64()
		total = total + size
		b.WriteString("\n")
		b.WriteString(utilgo.StringPadding(utilgo.DateFormat(item.Get("ctime").MustInt64()), 22))
		b.WriteString(utilgo.StringPadding(utilgo.DateFormat(item.Get("mtime").MustInt64()), 22))
		b.WriteString(utilgo.StringPadding(utilgo.ByteFormat(size), 10))
		b.WriteString(utilgo.StringPadding(item.Get("path").MustString(), 20))
	}
	fmt.Print(name + bc.root + "  ➜  " + p + " " + utilgo.ByteFormat(total))
	fmt.Println(b.String())
	return nil
}

// APILs response ls
func (bc *Bclient) APILs(p string) (*simplejson.Json, error) {
	body, err := httpGet(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "list", bc.token, path.Join(bc.root, p)))
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Cd() {

}

func (bc *Bclient) APICd(p string) (*simplejson.Json, error) {
	bc.path = p
	return bc.APILs(p)
}

func (bc *Bclient) Mkdir(p string) {

}

func (bc *Bclient) ApiMkdir(p string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "mkdir", bc.token, path.Join(bc.root, p)), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Mv(source string, target string) error {
	js, err := bc.APIMv(source, target)
	if err != nil {
		return err
	}
	errMsg := js.Get("error_msg").MustString()
	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}
	fmt.Println(name + bc.root)
	fmt.Println(source + " 已移动至 " + target)
	return nil
}

func (bc *Bclient) APIMv(source string, target string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&from=%s&to=%s", bc.apiURL, "move", bc.token, path.Join(bc.root, source), path.Join(bc.root, target)), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Cp(source string, target string) error {
	js, err := bc.APICp(source, target)
	if err != nil {
		return err
	}
	errMsg := js.Get("error_msg").MustString()
	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}
	fmt.Println(name + bc.root)
	fmt.Println(source + " 已复制至 " + target)
	return nil
}

func (bc *Bclient) APICp(source string, target string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&from=%s&to=%s", bc.apiURL, "copy", bc.token, path.Join(bc.root, source), path.Join(bc.root, target)), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Rm(file string) error {
	js, err := bc.APIRm(file)
	if err != nil {
		return err
	}
	errMsg := js.Get("error_msg").MustString()
	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}
	fmt.Println(name + bc.root)
	fmt.Println(file + " 已删除")
	return nil
}

func (bc *Bclient) APIRm(file string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "delete", bc.token, path.Join(bc.root, file)), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Get(file string) (io.ReadCloser, error) {
	return bc.APIGet(file)
}

func (bc *Bclient) APIGet(file string) (io.ReadCloser, error) {
	resp, err := httpGetResp(bc.GetDownloadURL(file))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// GetDownloadURL return download url
func (bc *Bclient) GetDownloadURL(file string) string {
	return fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "download", bc.token, path.Join(bc.root, file))
}

func (bc *Bclient) Put() {

}

func (bc *Bclient) APIPut(savePath string, overwrite bool) (*simplejson.Json, error) {
	var ondup = "newcopy"
	if overwrite {
		ondup = "overwrite"
	}
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s&ondup=%s", bc.uploadURL, "upload", bc.token, path.Join(bc.root, savePath), ondup), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) RapidPut() {

}

func (bc *Bclient) APIRapidPut(savePath string, fileSize uint64, md5Str string, sliceMd5 string, contentCrc32 string, overwrite bool) (*simplejson.Json, error) {
	var ondup = "newcopy"
	if overwrite {
		ondup = "overwrite"
	}
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s&content-length=%d&content-md5=%s&slice-md5=%s&content-crc32=%s&ondup=%s", bc.uploadURL, "rapidupload", bc.token, path.Join(bc.root, savePath), fileSize, md5Str, sliceMd5, contentCrc32, ondup), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

// Info print the disk usage
func (bc *Bclient) Info() error {
	js, err := bc.APIInfo()
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	quota := js.Get("quota").MustUint64()
	used := js.Get("used").MustUint64()
	b.WriteString(name + "\n总大小:" + utilgo.ByteFormat(quota))
	b.WriteString("\n已使用:" + utilgo.ByteFormat(used))
	b.WriteString(fmt.Sprintf("\n利用率:%.1f%%", float32(used)/float32(quota)*100))
	fmt.Println(b.String())
	return nil
}

// APIInfo response usage info
func (bc *Bclient) APIInfo() (*simplejson.Json, error) {
	body, err := httpGet(fmt.Sprintf("%s?method=%s&access_token=%s", bc.infoURL, "info", bc.token))
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

// Fileinfo print the file/dir info
func (bc *Bclient) Fileinfo(p string, dlink bool) error {
	js, err := bc.APIFileinfo(p)
	if err != nil {
		return err
	}
	errMsg := js.Get("error_msg").MustString()
	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}
	b := bytes.Buffer{}
	item := js.Get("list").GetIndex(0)
	b.WriteString(name + item.Get("path").MustString())
	b.WriteString("\n文件类型:" + utilgo.BoolString(item.Get("isdir").MustInt() == 0, "文件", "文件夹"))
	b.WriteString("\n文件大小:" + utilgo.ByteFormat(item.Get("size").MustUint64()))
	b.WriteString(fmt.Sprintf("\n文件字节:%d\n文件标识:%d", item.Get("size").MustInt64(), item.Get("fs_id").MustInt64()))
	b.WriteString("\n创建时间:" + utilgo.DateFormat(item.Get("ctime").MustInt64()))
	b.WriteString("\n修改时间:" + utilgo.DateFormat(item.Get("mtime").MustInt64()))
	blockstr := item.Get("block_list").MustString()
	if blockstr != "" {
		blocks, err := simplejson.NewJson([]byte(item.Get("block_list").MustString()))
		if err != nil {
			fmt.Println(b.String())
			return err
		}
		blocksarr := blocks.MustStringArray()
		if len(blocksarr) == 1 {
			b.WriteString("\n文件哈希:" + blocksarr[0])
		} else {
			for _, v := range blocksarr {
				b.WriteString("\n哈希块:" + v)
			}
		}
		if dlink {
			b.WriteString("\n下载地址:" + bc.GetDownloadURL(p))
		}
	}
	fmt.Println(b.String())
	return nil
}

// APIFileinfo response info
func (bc *Bclient) APIFileinfo(file string) (*simplejson.Json, error) {
	body, err := httpGet(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "meta", bc.token, path.Join(bc.root, file)))
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Search(fileName string) error {
	js, err := bc.APISearch(fileName)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	var total uint64
	for i := range js.Get("list").MustArray() {
		item := js.Get("list").GetIndex(i)
		size := item.Get("size").MustUint64()
		total = total + size
		b.WriteString("\n")
		b.WriteString(utilgo.StringPadding(utilgo.DateFormat(item.Get("ctime").MustInt64()), 22))
		b.WriteString(utilgo.StringPadding(utilgo.DateFormat(item.Get("mtime").MustInt64()), 22))
		b.WriteString(utilgo.StringPadding(utilgo.ByteFormat(size), 10))
		b.WriteString(utilgo.StringPadding(item.Get("path").MustString(), 20))
	}
	fmt.Print(name + bc.root + "  ➜  搜索[" + fileName + "] " + utilgo.ByteFormat(total))
	fmt.Println(b.String())
	return nil
}

func (bc *Bclient) APISearch(name string) (*simplejson.Json, error) {
	body, err := httpGet(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s&wd=%s&re=%s", bc.apiURL, "search", bc.token, bc.root, name, "1"))
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Taskadd() {

}

func (bc *Bclient) ApiTaskadd(savePath string, sourceURL string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&save_path=%s/&source_url=%s&app_id=250528", bc.taskURL, "add_task", bc.token, savePath, sourceURL), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Tasklist() {

}

func (bc *Bclient) APITasklist() (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&status=255&app_id=250528&need_task_info=1", bc.taskURL, "list_task", bc.token), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Taskinfo() {

}

func (bc *Bclient) APITaskinfo(ids string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&task_ids=%s&app_id=250528", bc.taskURL, "query_task", bc.token, ids), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Taskremove() {

}

func (bc *Bclient) APITaskremove(id string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&task_id=%s&app_id=250528", bc.taskURL, "cancel_task", bc.token, id), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Clear() {

}

func (bc *Bclient) APIClear() (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&type=%s", bc.apiURL, "delete", bc.token, "recycle"), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func showTaskStatus(status string) string {
	if v, ok := taskStatusMap[status]; ok {
		return v
	}
	return "未知状态"
}
