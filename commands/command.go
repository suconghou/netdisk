package commands

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/suconghou/netdisk/config"
	"github.com/suconghou/netdisk/layers/fslayer"
	"github.com/suconghou/netdisk/middleware"
	"github.com/suconghou/netdisk/tools"
	"github.com/suconghou/netdisk/util"
	"github.com/suconghou/utilgo"
	"golang.org/x/net/proxy"
)

// Ls list files
func Ls() {
	var dir string
	if len(os.Args) >= 3 {
		dir = os.Args[2]
	}
	err := fslayer.ListDir(dir, false)
	if err != nil {
		util.Log.Printf("%v", err)
	}
}

// Cd enter dir and list files
func Cd() {
	var dir string
	if len(os.Args) == 3 {
		dir = os.Args[2]
		err := fslayer.ListDir(dir, true)
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		util.Log.Print("Usage:disk cd newpath")
	}
}

// Pwd print current work dir
func Pwd() {
	err := fslayer.Pwd()
	if err != nil {
		util.Log.Print(err)
	}
}

// Mv move file from the backend
func Mv() {
	if len(os.Args) == 4 {
		err := fslayer.MoveFile(os.Args[2], os.Args[3])
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		util.Log.Print("Usage:disk mv path newpath")
	}
}

// Cp copy file from the backend
func Cp() {
	if len(os.Args) == 4 {
		err := fslayer.CopyFile(os.Args[2], os.Args[3])
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		util.Log.Print("Usage:disk cp path newpath")
	}
}

// Mkdir mkdir to the backend
func Mkdir() {
	if len(os.Args) == 3 {
		err := fslayer.Mkdir(os.Args[2])
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		util.Log.Print("Usage:disk mkdir path")
	}
}

// Rm delete file from the backend
func Rm() {
	if len(os.Args) == 3 {
		err := fslayer.DeleteFile(os.Args[2])
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		util.Log.Print("Usage:disk rm filepath")
	}
}

// Get do a simple download
func Get() {
	if len(os.Args) >= 3 && !utilgo.IsURL(os.Args[2], true) {
		saveas, err := utilgo.GetStorePath(os.Args[2])
		if err != nil {
			util.Log.Print(err)
			return
		}
		transport, err := util.GetProxy()
		if err != nil {
			util.Log.Print(err)
			return
		}
		err = fslayer.Get(os.Args[2], saveas, transport)
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		util.Log.Print("Usage:disk get filepath")
	}
}

// Put upload file to the backend
func Put() {
	var (
		file      = os.Stdin
		overwrite = true
		fileName  = "-"
		err       error
	)
	if len(os.Args) >= 3 {
		fileName = os.Args[2]
		if fileName != "-" {
			file, err = utilgo.GetOpenFile(fileName)
		}
		overwrite = utilgo.HasFlag(os.Args, "-f")
		saveName, _ := utilgo.GetParam(os.Args, "-f")
		if saveName != "" {
			fileName = saveName
		}
	}
	if err == nil {
		defer file.Close()
		err = fslayer.Put(fileName, overwrite, file)
	}
	if err != nil {
		util.Log.Print(err)
	}
}

// Wget url like wget
func Wget() {
	if len(os.Args) >= 3 && utilgo.IsURL(os.Args[2], true) {
		var (
			saveas, err = utilgo.GetStorePath(os.Args[2])
		)
		if err != nil {
			util.Log.Print(err)
			return
		}
		transport, err := util.GetProxy()
		if err != nil {
			util.Log.Print(err)
			return
		}
		err = fslayer.WgetURL(os.Args[2], saveas, transport)
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		util.Log.Print("Usage:disk wget url")
	}
}

// Play play a url or file(pcs file)
func Play() {
	if len(os.Args) >= 3 {
		var (
			saveas string
			err    error
			stdout = utilgo.HasFlag(os.Args, "--stdout")
		)
		if stdout {
			util.Log.SetOutput(os.Stderr)
		}
		saveas, err = utilgo.GetStorePath(os.Args[2])
		if err != nil {
			util.Log.Print(err)
			return
		}
		transport, err := util.GetProxy()
		if err != nil {
			util.Log.Print(err)
			return
		}
		util.Log.Print("Playing " + saveas)
		if utilgo.IsURL(os.Args[2], true) {
			err = fslayer.PlayURL(os.Args[2], saveas, stdout, transport)
			if err != nil {
				util.Log.Print(err)
			}
		} else {
			err = fslayer.Play(os.Args[2], saveas, stdout, transport)
			if err != nil {
				util.Log.Print(err)
			}
		}
	} else {
		util.Log.Print("Usage:disk play filepath/url")
	}
}

// Info print the backend info or file info
func Info() {
	if len(os.Args) >= 3 {
		err := fslayer.GetFileInfo(os.Args[2], utilgo.HasFlag(os.Args, "--link"))
		if err != nil {
			util.Log.Print(err)
		}
	} else {
		err := fslayer.GetInfo()
		if err != nil {
			util.Log.Print(err)
		}
	}
}

// Hash print the sha1sum sha256sum
func Hash(t string) {
	var (
		l    = len(os.Args)
		err  error
		file *os.File
		x    []byte
	)
	if l >= 3 {
		for i := 2; i < l; i++ {
			file, err = utilgo.GetOpenFile(os.Args[i])
			if err == nil {
				x, err = utilgo.GetFileHash(file, t)
				file.Close()
				if err == nil {
					util.Log.Printf("%x  %s", x, filepath.Base(file.Name()))
				}
			}
			if err != nil {
				util.Log.Print(err)
			}
		}
	} else {
		util.Log.Print("Usage:disk hash file")
	}
}

// Help print the help message
func Help() {
	util.Log.Print(os.Args[0] + " ls info mv cp get put wget play rm mkdir pwd hash config empty search task ")
}

// Task list current backend task
func Task() {
	var err error
	if (len(os.Args) == 2) || (os.Args[2] == "list") {
		err = fslayer.GetTaskList()
	} else if len(os.Args) == 5 && os.Args[2] == "add" {
		err = fslayer.AddTask((os.Args[3]), os.Args[4])
	} else if len(os.Args) == 4 && os.Args[2] == "remove" {
		err = fslayer.RemoveTask(os.Args[3])
	} else if len(os.Args) == 4 && os.Args[2] == "info" {
		err = fslayer.GetTaskInfo(os.Args[3])
	} else {
		util.Log.Print("Usage:disk task list/add/info/remove")
	}
	if err != nil {
		util.Log.Print(err)
	}
}

// Search form the backend
func Search() {
	if len(os.Args) == 3 {
		err := fslayer.SearchFile(os.Args[2])
		if err != nil {
			util.Log.Print(err)
		}
	}
}

// Empty clear cache data
func Empty() {
	fslayer.Empty()
}

// Serve start a http file server
func Serve() {
	var (
		port        int
		root        string
		print       bool
		ferr        flag.ErrorHandling
		CommandLine = flag.NewFlagSet(os.Args[1], ferr)
	)
	CommandLine.IntVar(&port, "p", 6060, "listen port")
	CommandLine.StringVar(&root, "d", "./", "document root")
	CommandLine.BoolVar(&print, "l", false, "print address")
	err := CommandLine.Parse(os.Args[2:])
	if err == nil {
		root, err = utilgo.PathMustHave(root)
		if err == nil {
			var (
				addr string
				ip   *net.IPNet
			)
			if print {
				ip, err = utilgo.GetCurIpv4()
				if err == nil {
					addr = fmt.Sprintf("http://%s:%d", ip.IP.String(), port)
				}
			} else {
				addr = fmt.Sprintf("port %d", port)
			}
			util.Log.Printf("Starting up on %s\nDocument root %s", addr, root)
			err = http.ListenAndServe(":"+strconv.Itoa(port), http.FileServer(http.Dir(root)))
		}
	}
	if err != nil {
		util.Log.Print(err)
	}
}

// Proxy enable a socks proxy server
func Proxy() {
	var (
		port        int
		socks       string
		ferr        flag.ErrorHandling
		l           net.Listener
		CommandLine = flag.NewFlagSet(os.Args[1], ferr)
		dialer      = proxy.FromEnvironment()
		d           proxy.Dialer
	)
	CommandLine.IntVar(&port, "p", 8123, "listen port")
	CommandLine.StringVar(&socks, "socks", "", "socks proxy")
	err := CommandLine.Parse(os.Args[2:])
	if err == nil {
		if socks != "" {
			d, err = proxy.SOCKS5("tcp", socks, nil, proxy.Direct)
			if err == nil {
				dialer = d
			}
		}
		if err != nil {
			util.Log.Print(err)
			return
		}
		util.Log.Printf("Starting up on port %d", port)
		l, err = net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			for {
				client, err := l.Accept()
				if err == nil {
					go func() {
						defer func() {
							if err := recover(); err != nil {
								util.Log.Print(err)
							}
						}()
						err := middleware.ProxySocks(client, dialer)
						if err != nil && err != io.EOF {
							util.Log.Print(err)
						}
					}()
				} else {
					util.Log.Print(err)
				}
			}
		}
	}
	if err != nil {
		util.Log.Print(err)
	}
}

// HTTPProxy is a http reverse proxy like nginx but can use given upstream
func HTTPProxy() {
	var (
		port        int
		url         string
		proxy       string
		socks       string
		header      string
		ferr        flag.ErrorHandling
		CommandLine = flag.NewFlagSet(os.Args[1], ferr)
		transport   *http.Transport
	)
	CommandLine.IntVar(&port, "p", 8123, "listen port")
	CommandLine.StringVar(&url, "u", "http://127.0.0.1:8080", "reverse url")
	CommandLine.StringVar(&proxy, "proxy", "", "http proxy")
	CommandLine.StringVar(&socks, "socks", "", "socks proxy")
	CommandLine.StringVar(&header, "header", "", "allow headers")
	err := CommandLine.Parse(os.Args[2:])
	if err == nil {
		if socks != "" {
			transport, err = util.MakeSocksProxy(socks, util.GetTLSConfig())
		} else if proxy != "" {
			transport, err = util.MakeHTTPProxy(proxy, util.GetTLSConfig())
		}
		if err == nil {
			err = tools.HTTPProxy(port, url, transport, header)
		}
	}
	if err != nil {
		util.Log.Print(err)
	}
}

// Nc like but have progress bar
func Nc() {
	err := tools.NcMain()
	if err != nil {
		util.Debug.Printf("nc error:%s", err)
	}
}

// Fwd port forward
func Fwd() {
	err := tools.FwdMain()
	if err != nil {
		util.Debug.Printf("fwd error:%s", err)
	}
}

// Network test http speed
func Network() {
	var (
		chunk       uint
		timeout     uint
		input       string
		proxy       string
		socks       string
		host        string
		path        string
		https       bool
		ferr        flag.ErrorHandling
		CommandLine = flag.NewFlagSet(os.Args[1], ferr)
		transport   *http.Transport
	)
	CommandLine.UintVar(&chunk, "s", 256, "chunk size")
	CommandLine.UintVar(&timeout, "t", 15, "timeout")
	CommandLine.StringVar(&input, "i", "-", "input file")
	CommandLine.StringVar(&proxy, "proxy", "", "http proxy")
	CommandLine.StringVar(&socks, "socks", "", "socks proxy")
	CommandLine.StringVar(&host, "host", "", "http host")
	CommandLine.StringVar(&path, "path", "", "http path")
	CommandLine.BoolVar(&https, "https", false, "use https")

	err := CommandLine.Parse(os.Args[2:])
	if err == nil {
		if socks != "" {
			transport, err = util.MakeSocksProxy(socks, util.GetTLSConfig())
		} else if proxy != "" {
			transport, err = util.MakeHTTPProxy(proxy, util.GetTLSConfig())
		}
		if err == nil {
			if host == "" {
				err = tools.SpeedTest(input, chunk, timeout, transport)
			} else {
				err = tools.SpeedTestWithHost(input, host, path, https, chunk, timeout, transport)
			}
		}
	}
	if err != nil {
		util.Log.Print(err)
	}

}

// Usage print help message
func Usage() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		util.Log.Print(os.Args[0] + " version: disk/" + config.Version + "\n" + config.ReleaseURL)
	} else {
		Help()
	}
}
