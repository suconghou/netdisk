package baidudisk

import (
	"fmt"
	"io"
	"path"

	"github.com/bitly/go-simplejson"
)

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

func (bc *Bclient) Ls(p string) {

}

func (bc *Bclient) ApiLs(p string) (*simplejson.Json, error) {
	body, err := httpGet(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "list", bc.token, path.Join(bc.root, p)))
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Cd() {

}

func (bc *Bclient) ApiCd(p string) (*simplejson.Json, error) {
	bc.path = p
	return bc.ApiLs(p)
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

func (bc *Bclient) Mv(source string, target string) {

}

func (bc *Bclient) ApiMv(source string, target string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&from=%s&to=%s", bc.apiURL, "move", bc.token, path.Join(bc.root, source), path.Join(bc.root, target)), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Cp() {

}

func (bc *Bclient) ApiCp(source string, target string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&from=%s&to=%s", bc.apiURL, "copy", bc.token, path.Join(bc.root, source), path.Join(bc.root, target)), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Rm(file string) {

}

func (bc *Bclient) ApiRm(file string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "delete", bc.token, path.Join(bc.root, file)), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Get(file string) (io.ReadCloser, error) {
	return bc.ApiGet(file)
}

func (bc *Bclient) ApiGet(file string) (io.ReadCloser, error) {
	resp, err := httpGetResp(bc.GetDownloadURL(file))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (bc *Bclient) GetDownloadURL(file string) string {
	return fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "download", bc.token, file)
}

func (bc *Bclient) Put() {

}

func (bc *Bclient) ApiPut(savePath string, overwrite bool) (*simplejson.Json, error) {
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

func (bc *Bclient) ApiRapidPut(savePath string, fileSize uint64, md5Str string, sliceMd5 string, contentCrc32 string, overwrite bool) (*simplejson.Json, error) {
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

func (bc *Bclient) Info() {

}

func (bc *Bclient) ApiInfo() (*simplejson.Json, error) {
	body, err := httpGet(fmt.Sprintf("%s?method=%s&access_token=%s", bc.infoURL, "info", bc.token))
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Fileinfo() {

}

func (bc *Bclient) ApiFileinfo(file string) (*simplejson.Json, error) {
	body, err := httpGet(fmt.Sprintf("%s?method=%s&access_token=%s&path=%s", bc.apiURL, "meta", bc.token, path.Join(bc.root, file)))
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Search() {

}

func (bc *Bclient) ApiSearch(name string) (*simplejson.Json, error) {
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

func (bc *Bclient) ApiTasklist() (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&status=255&app_id=250528&need_task_info=1", bc.taskURL, "list_task", bc.token), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Taskinfo() {

}

func (bc *Bclient) ApiTaskinfo(ids string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&task_ids=%s&app_id=250528", bc.taskURL, "query_task", bc.token, ids), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Taskremove() {

}

func (bc *Bclient) ApiTaskremove(id string) (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&task_id=%s&app_id=250528", bc.taskURL, "cancel_task", bc.token, id), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}

func (bc *Bclient) Clear() {

}

func (bc *Bclient) ApiClear() (*simplejson.Json, error) {
	body, err := httpPost(fmt.Sprintf("%s?method=%s&access_token=%s&type=%s", bc.apiURL, "delete", bc.token, "recycle"), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(body)
}
