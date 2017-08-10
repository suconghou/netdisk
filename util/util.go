package util

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/suconghou/utilgo"
	"golang.org/x/net/proxy"
)

// Log is a global logger
var Log = log.New(os.Stdout, "", 0)

// Debug log to stderr
var Debug = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

// GetCrc32AndMd5 ...
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

// MakeSocksProxy return socks proxy Transport
func MakeSocksProxy(str string) (*http.Transport, error) {
	dialer, err := proxy.SOCKS5("tcp", str, nil, proxy.Direct)
	if err == nil {
		return &http.Transport{Dial: dialer.Dial}, nil
	}
	return nil, err
}

// MakeHTTPProxy return http proxy Transport
func MakeHTTPProxy(str string) (*http.Transport, error) {
	urli := url.URL{}
	urlproxy, err := urli.Parse(str)
	if err == nil {
		return &http.Transport{
			Proxy:           http.ProxyURL(urlproxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}, nil
	}
	return nil, err
}

// GetProxy return http.Transport or nil
func GetProxy() (*http.Transport, error) {
	if str, err := utilgo.GetParam("--socks"); err == nil {
		return MakeSocksProxy(str)
	}
	if str, err := utilgo.GetParam("--proxy"); err == nil {
		return MakeHTTPProxy(str)
	}
	return nil, nil
}
