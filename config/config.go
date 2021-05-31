package config

import (
	"encoding/json"
	"os"
	"runtime"
)

var configPath = "/etc/disk.json"

// Version and ReleaseURL
const (
	Version    = "0.1.4"
	ReleaseURL = "https://github.com/suconghou/netdisk"
)

// Appcfg config
type appcfg struct {
	Token string
	Root  string
	Path  string
}

// Cfg config the whole app
var Cfg appcfg

func init() {
	if runtime.GOOS == "windows" {
		configPath = `C:\Users\Default\disk.json`
	}
	// 即时没有配置文件,也允许运行
	loadConfig()
}

func loadConfig() error {
	strJSON, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(strJSON), &Cfg)
}

func (Cfg *appcfg) Save() error {
	strJSON, err := json.Marshal(Cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, strJSON, 0777)
}
