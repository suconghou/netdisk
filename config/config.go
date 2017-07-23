package config

import (
	"encoding/json"
	"io/ioutil"
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
	Driver string
	Pcs    struct {
		Token string
		Root  string
		Path  string
	}
}

// Cfg config the whole app
var Cfg appcfg

func init() {
	if runtime.GOOS == "windows" {
		configPath = `C:\Users\Default\disk.json`
	}
	loadConfig()
}

func loadConfig() error {
	strJSON, err := ioutil.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			Cfg.Pcs.Root = "/apps/suconghou"
			Cfg.Driver = "pcs"
		} else {
			return err
		}
	}
	return json.Unmarshal([]byte(strJSON), &Cfg)
}

func (Cfg *appcfg) Save() error {
	strJSON, err := json.Marshal(Cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, strJSON, 0777)
}

// IsPcs return if its driver is pcs
func IsPcs() bool {
	if Cfg.Driver == "pcs" {
		return true
	}
	return false
}

func Use(driver string) error {
	Cfg.Save()
	return nil
}

func ConfigList() {

}
