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

var urlReg = regexp.MustCompile(`^(?i:https?)://[[:print:]]+$`)

// Log is a global logger
var Log = log.New(os.Stdout, "", 0)

// Debug log to stderr
var Debug = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

// MakeSocksProxy return socks proxy Transport
func MakeSocksProxy(str string, tlsCfg *tls.Config) (*http.Transport, error) {
	dialer, err := proxy.SOCKS5("tcp", str, nil, proxy.Direct)
	if err == nil {
		if tlsCfg != nil {
			return &http.Transport{Dial: dialer.Dial, TLSClientConfig: tlsCfg}, nil
		}
		return &http.Transport{Dial: dialer.Dial}, nil
	}
	return nil, err
}

// MakeHTTPProxy return http proxy Transport
func MakeHTTPProxy(str string, tlsCfg *tls.Config) (*http.Transport, error) {
	urli := url.URL{}
	urlproxy, err := urli.Parse(str)
	if err == nil {
		if tlsCfg != nil {
			return &http.Transport{Proxy: http.ProxyURL(urlproxy), TLSClientConfig: tlsCfg}, nil
		}
		return &http.Transport{Proxy: http.ProxyURL(urlproxy)}, nil
	}
	return nil, err
}

// GetTLSConfig return if skipVerify
func GetTLSConfig() *tls.Config {
	var (
		tlsCfg     *tls.Config
		skipVerify = utilgo.HasFlag("--no-check-certificate")
	)
	if skipVerify {
		tlsCfg = &tls.Config{InsecureSkipVerify: true}
	}
	return tlsCfg
}

// GetProxy return http.Transport or nil
func GetProxy() (*http.Transport, error) {
	tlsCfg := GetTLSConfig()
	if str, err := utilgo.GetParam("--socks"); err == nil {
		return MakeSocksProxy(str, tlsCfg)
	}
	if str, err := utilgo.GetParam("--proxy"); err == nil {
		return MakeHTTPProxy(str, tlsCfg)
	}
	if tlsCfg != nil {
		return &http.Transport{TLSClientConfig: tlsCfg}, nil
	}
	return nil, nil
}

// IsURL if the given string is an url
func IsURL(url string) bool {
	return urlReg.MatchString(url)
}

// GetMirrors return url mirrors
func GetMirrors() map[string]int {
	found := false
	var mirrors = map[string]int{}
	for _, item := range os.Args {
		if !found {
			if item == "--mirrors" {
				found = true
			}
		} else if IsURL(item) {
			mirrors[item] = 1
		}
	}
	if len(mirrors) > 0 {
		return mirrors
	}
	return nil
}
