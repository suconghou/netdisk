package util

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/suconghou/utilgo"
	"golang.org/x/net/proxy"
)

var urlReg = regexp.MustCompile(`^(?i:https?)://[[:print:]]{1,500}$`)

// Log is a global logger
var Log = log.New(os.Stdout, "", 0)

// Debug log to stderr
var Debug = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

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

// IsURL if the given string is an url
func IsURL(url string) bool {
	return urlReg.MatchString(url)
}
