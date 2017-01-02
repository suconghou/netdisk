package config

import (
	"encoding/json"
	_ "flag"
	"fmt"
	"io/ioutil"
	"os"
)

var ConfigPath string = "/etc/disk.json"

type Config struct {
	Token string
	Root  string
	Path  string
	Disk  string
}

var Cfg Config

func LoadConfig() Config {
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

func ConfigSet(key string, value string) {
	fmt.Println(key)
	fmt.Println(value)
}

func ConfigGet(key string) {

}

func ConfigList() {

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

func Error() {

}
