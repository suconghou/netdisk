package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

// Log is a global logger
var Log = log.New(os.Stdout, "", 0)

// Debug log to stderr
var Debug = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func ByteFormat(bytes uint64) string {
	unit := [...]string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	if bytes >= 1024 {
		e := math.Floor(math.Log(float64(bytes)) / math.Log(float64(1024)))
		return fmt.Sprintf("%.2f%s", float64(bytes)/math.Pow(1024, math.Floor(e)), unit[int(e)])
	}
	return fmt.Sprintf("%d%s", bytes, unit[0])
}

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

func StringPad(str string, le int) string {
	l := le - len(str)
	if l > 0 {
		for i := 0; i < l; i++ {
			str = str + " "
		}
	}
	return str
}

func DateFormat(times uint64) string {
	t := int64(times)
	var str string
	if time.Unix(t, 0).Format("06/01/02") == time.Now().Format("06/01/02") {
		str = time.Unix(t, 0).Format("15:04:05")
	} else {
		str = time.Unix(t, 0).Format("06/01/02")
	}
	return str
}

func DateS(times int64) string {
	return time.Unix(times, 0).Format("2006/01/02 15:04:05")
}

func PrintMd5(filePath string) {
	file, err := os.Open(filePath)
	if err == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		fmt.Printf("%s   %x\n", filePath, md5h.Sum([]byte(""))) //md5
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}

func FileOk(filePath string) (uint64, string) {
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
		stat, err := os.Stat(filePath)
		if err != nil {
			panic(err)
		} else {
			fileSize := uint64(stat.Size())
			md5h := md5.New()
			io.Copy(md5h, file)
			md5Str := hex.EncodeToString(md5h.Sum([]byte("")))
			fmt.Printf("%s  %s  %s \n", filePath, md5Str, ByteFormat(fileSize)) //md5
			return fileSize, md5Str
		}
	}
	return 0, ""
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

func JSONPut(w http.ResponseWriter, bs []byte, httpCache bool, cacheTime uint32) {
	CrossShare(w)
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if httpCache {
		UseHTTPCache(w, cacheTime)
	}
	w.Write(bs)
}

func UseHTTPCache(w http.ResponseWriter, cacheTime uint32) {
	w.Header().Set("Expires", time.Now().Add(time.Second*time.Duration(cacheTime)).Format(http.TimeFormat))
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", cacheTime))
}

func CrossShare(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Content-Length, Accept, Accept-Encoding")
}
