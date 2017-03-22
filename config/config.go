package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

var ConfigPath string = "/etc/disk.json"

var Version string = "0.1.4"

var ReleaseUrl = "https://github.com/suconghou/netdisk"

type Config struct {
	Token string
	Root  string
	Path  string
	Disk  string
}

var Cfg Config

func LoadConfig() Config {

	if runtime.GOOS == "windows" {
		ConfigPath = `C:\Users\Default\disk.json`
	}
	strJson, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		Cfg.Token = "token"
		Cfg.Root = "/apps/suconghou"
		Cfg.Path = ""
		Cfg.Disk = "pcs"
		if os.IsNotExist(err) {
			SaveConfig()
			return Cfg
		} else if os.IsPermission(err) {
			fmt.Println(err)
			os.Exit(1)
		} else {
			panic(err)
		}
		return Cfg
	} else {
		cfg := &Config{}
		err = json.Unmarshal([]byte(strJson), &cfg)
		return *cfg
	}

}

func (cfg *Config) getToken() string {
	return cfg.Token
}

func (cfg *Config) setToken(token string) {
	cfg.Token = token
}

func (cfg *Config) setApp(app string) {
	cfg.Root = "/apps/" + app
}

func ConfigSet(value string) {
	Cfg.setToken(value)
	SaveConfig()
}

func ConfigSetApp(value string) {
	Cfg.setApp(value)
	SaveConfig()
}

func ConfigGet() {
	fmt.Println("Token:" + Cfg.getToken())
}

func ConfigList() {
	fmt.Println(fmt.Sprintf("Root:%s\nPath:%s\nToken:%s", Cfg.Root, Cfg.Path, Cfg.Token))
}

func SaveConfig() {
	strJson, err := json.Marshal(Cfg)
	if err != nil {
		panic(err)
	} else {
		err := ioutil.WriteFile(ConfigPath, strJson, 0666)
		if err != nil {
			if os.IsPermission(err) {
				fmt.Println(err)
				os.Exit(1)
			} else if os.IsExist(err) {
				fmt.Println(err)
				os.Exit(2)
			} else {
				panic(err)
			}
		}
	}
}
