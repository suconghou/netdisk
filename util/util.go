package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

func ByteFormat(bytes int64) string {
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

func DateFormat(times int64) string {
	var str string
	if time.Unix(times, 0).Format("06/01/02") == time.Now().Format("06/01/02") {
		str = time.Unix(times, 0).Format("15:04:05")
	} else {
		str = time.Unix(times, 0).Format("06/01/02")
	}
	return str
}

func BoolString(b bool, s, s1 string) string {
	if b {
		return s
	}
	return s1
}

func PrintMd5(filePath string) {
	file, inerr := os.Open(filePath)
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		fmt.Printf("%s   %x\n",filePath, md5h.Sum([]byte(""))) //md5
	} else {
		fmt.Println(inerr)
		os.Exit(1)
	}
}

func Bar(vl int, width int) string {

	return fmt.Sprintf("%s %s", strings.Repeat("█", vl/(100/width)), strings.Repeat(" ", width-vl/(100/width)))
}
