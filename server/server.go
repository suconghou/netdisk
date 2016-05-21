package server

import (
	"io/ioutil"
	"os"
)

func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
			files = append(files, dirPth+PthSep+fi.Name())
	}
	return files, nil
}

func md5File() {

}

func fileSize() {

}
