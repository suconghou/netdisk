package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/suconghou/utilgo"
	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
)

var (
	publicDir = "/Public"
	client    *dropy.Client
)

type fileInfo struct {
	Name    string
	Size    int64
	IsDir   bool
	Path    string
	ModTime time.Time
}

type fileInfoList struct {
	List  []fileInfo
	Total int
}

func init() {
	token := os.Getenv("DROPBOX_ACCESS_TOKEN")
	client = dropy.New(dropbox.New(dropbox.NewConfig(token)))
}

// Files serve dropbox files
func Files(w http.ResponseWriter, r *http.Request, match []string) {
	err := func() error {
		filePath := path.Join(publicDir, match[0])
		file, err := client.Stat(filePath)
		if err != nil {
			return err
		}
		if file.IsDir() {
			fileList, err := client.List(filePath)
			if err != nil {
				return err
			}
			infoList := make([]fileInfo, 0)
			for _, item := range fileList {
				info := fileInfo{item.Name(), item.Size(), item.IsDir(), path.Join(match[0], item.Name()), item.ModTime()}
				infoList = append(infoList, info)
			}
			res := fileInfoList{Total: len(fileList), List: infoList}
			bs, err := json.Marshal(&res)
			if err != nil {
				return err
			}
			utilgo.JSONPut(w, bs, true, 3600)
		} else {
			h := w.Header()
			utilgo.CrossShare(h, nil, "")
			utilgo.UseHTTPCache(h, 86400)
			h.Set("Content-Length", strconv.Itoa(int(file.Size())))
			_, err := io.Copy(w, client.Open(filePath))
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Imgs use some other disk backend
func Imgs(w http.ResponseWriter, r *http.Request, match []string) {

}
