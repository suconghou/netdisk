package commands

import (
	"config"
	"fmt"
	"layers/fslayer"
	"os"
	"path/filepath"
	"strings"
	"util"
)

func Use() {
	fmt.Println("select default remote connect")
	fmt.Println(config.Cfg)
}

func Ls() {
	var path string
	if len(os.Args) == 3 {
		path = absNoRoot(os.Args[2])
	} else {
		path = config.Cfg.Path
	}
	fslayer.ListDir(path)
}

func Cd() {
	if len(os.Args) == 3 {
		config.Cfg.Path = absNoRoot(os.Args[2])
		ret := fslayer.ListDir(config.Cfg.Path)
		if ret {
			config.SaveConfig()
		}
	} else {
		fmt.Println("change dir ")
	}
}

func Pwd() {
	fmt.Println(util.DiskName(config.Cfg.Disk) + config.Cfg.Root + "  ➜  " + config.Cfg.Path)
}

func Mkdir() {

}

func Rm() {
	fmt.Println("rm files")
}

func Get() {
	if len(os.Args) >= 3 {
		var dist string = ""
		if len(os.Args) >= 4 {
			dist = filepath.Clean(os.Args[3])
		}
		var path = absPath(os.Args[2])
		ok, size, hash := fslayer.GetFileInfo(path, false)
		if ok {
			fslayer.Get(path, dist, size, hash)
		}
	} else {
		fmt.Println("Usage get filepath saveas")
	}
}

func Put() {

}

func Wget() {
	if len(os.Args) >= 3 {
		var dist string = ""
		if len(os.Args) >= 4 {
			dist = filepath.Clean(os.Args[3])
		}
		var path = absPath(os.Args[2])
		ok, size, hash := fslayer.GetFileInfo(path, false)
		if ok {
			fslayer.Wget(path, dist, size, hash)
		}
	} else {
		fmt.Println("Usage get filepath saveas")
	}
}

func Sync() {

}

func Info() {
	if len(os.Args) == 2 {
		fslayer.GetInfo()
	} else {
		fslayer.GetFileInfo(absPath(os.Args[2]), false)
	}
}

func Hash() {
	if len(os.Args) == 3 {
		var filepath string = absLocalPath(os.Args[2])
		util.PrintMd5(filepath)
	} else {
		fmt.Println("hash file")
	}
}

func Play() {
	if len(os.Args) >= 3 {
		var dist string = ""
		if len(os.Args) >= 4 && (!strings.Contains(os.Args[3], "-")) {
			dist = filepath.Clean(os.Args[3])
		}
		var path = absPath(os.Args[2])
		ok, size, hash := fslayer.GetFileInfo(path, false)
		if ok {
			fslayer.GetPlayStream(path, dist, size, hash)
		}
	} else {
		fmt.Println("Usage get filepath saveas")
	}
}

func Help() {

}

func Config() {
	if (len(os.Args) == 2) || (os.Args[2] == "list") {
		config.ConfigList()
	} else if len(os.Args) == 4 && os.Args[2] == "get" {
		config.ConfigGet(os.Args[3])
	} else if len(os.Args) == 5 && os.Args[2] == "set" {
		config.ConfigSet(os.Args[3], os.Args[4])
	} else {
		config.Error()
	}
}

func Usage() {

	fmt.Println("ls cd rm get put hash ")

}

func Daemon() {
	fmt.Println("daemon start")
}

func absPath(path string) string {
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		if !strings.HasPrefix(path, config.Cfg.Root) {
			path = filepath.Join(config.Cfg.Root, "."+path)
		}
	} else {
		path = filepath.Join(config.Cfg.Root, config.Cfg.Path, path)
	}
	return path
}

func absNoRoot(path string) string {
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		path = path
	} else {
		path = filepath.Join(config.Cfg.Path, path)
	}
	return path
}

func absLocalPath(path string) string {
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		path = path
	} else {
		dir, _ := os.Getwd()
		path = filepath.Join(dir, filepath.Base(path))
	}
	return path
}
