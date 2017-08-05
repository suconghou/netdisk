package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
)

// Log is a global logger
var Log = log.New(os.Stdout, "", 0)

// Debug log to stderr
var Debug = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func DiskName(code string) string {

	unit := map[string]string{
		"pcs":     "百度网盘",
		"dropbox": "DROPBOX",
	}
	if v, ok := unit[code]; ok {
		return v
	} else {
		return "Unknow"
	}

}

func GetCrc32AndMd5(filePath string) (string, string) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(err)
			os.Exit(1)
		} else if os.IsPermission(err) {
			fmt.Println(err)
			os.Exit(1)
		} else {
			panic(err)
		}
	} else {
		defer file.Close()
		crc32h := crc32.NewIEEE()
		data := make([]byte, 262144)
		io.Copy(crc32h, file)
		file.ReadAt(data, 0)
		crc32Str := hex.EncodeToString(crc32h.Sum(nil))
		md5Str := fmt.Sprintf("%x", md5.Sum(data))
		return crc32Str, md5Str
	}
	return "", ""
}
