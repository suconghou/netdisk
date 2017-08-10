package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"netdisk/commands"
	"netdisk/route"
	"netdisk/util"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/suconghou/utilgo"
)

var (
	startTime = time.Now()
	port      int
	doc       string
)

var sysStatus struct {
	Uptime       string
	GoVersion    string
	Hostname     string
	MemAllocated uint64
	MemTotal     uint64
	MemSys       uint64
	NumGoroutine int
	CPUNum       int
	Pid          int
}

func main() {
	if len(os.Args) > 1 {
		cli()
	} else {
		err := daemon()
		if err != nil {
			util.Log.Print(err)
		}
	}
}

func cli() {
	switch os.Args[1] {
	case "use":
		commands.Use()
	case "ls":
		commands.Ls()
	case "cd":
		commands.Cd()
	case "pwd":
		commands.Pwd()
	case "cp":
		commands.Cp()
	case "mv":
		commands.Mv()
	case "mkdir":
		commands.Mkdir()
	case "rm":
		commands.Rm()
	case "get":
		commands.Get()
	case "put":
		commands.Put()
	case "wget":
		commands.Wget()
	case "sync":
		commands.Sync()
	case "info":
		commands.Info()
	case "hash":
		commands.Hash("")
	case "md5", "md5sum":
		commands.Hash("md5")
	case "crc32":
		commands.Hash("crc32")
	case "play":
		commands.Play()
	case "help":
		commands.Help()
	case "config":
		commands.Config()
	case "task":
		commands.Task()
	case "search":
		commands.Search()
	case "empty":
		commands.Empty()
	case "serve":
		commands.Serve()
	case "proxy":
		commands.Proxy()
	case "reverse":
		commands.HTTPProxy()
	case "nc":
		commands.Nc()
	default:
		commands.Usage()
	}
}

func daemon() error {
	flag.IntVar(&port, "p", 6060, "give me a port number")
	flag.StringVar(&doc, "d", "./", "document root dir")
	flag.Parse()
	pwd, err := utilgo.PathMustHave(doc)
	if err != nil {
		return err
	}
	http.HandleFunc("/status", status)
	http.HandleFunc("/", routeMatch)
	util.Log.Printf("Starting up on port %d\nDocument root %s", port, pwd)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func status(w http.ResponseWriter, r *http.Request) {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	sysStatus.Uptime = time.Since(startTime).String()
	sysStatus.NumGoroutine = runtime.NumGoroutine()
	sysStatus.MemAllocated = memStat.Alloc
	sysStatus.MemTotal = memStat.TotalAlloc
	sysStatus.MemSys = memStat.Sys
	sysStatus.CPUNum = runtime.NumCPU()
	sysStatus.GoVersion = runtime.Version()
	sysStatus.Hostname, _ = os.Hostname()
	sysStatus.Pid = os.Getpid()
	if bs, err := json.Marshal(&sysStatus); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
	} else {
		utilgo.JSONPut(w, bs, true, 60)
	}
}

func routeMatch(w http.ResponseWriter, r *http.Request) {
	found := false
	for _, p := range route.RoutePath {
		if p.Reg.MatchString(r.URL.Path) {
			found = true
			p.Handler(w, r, p.Reg.FindStringSubmatch(r.URL.Path))
			break
		}
	}
	if !found {
		fallback(w, r)
	}
}

func fallback(w http.ResponseWriter, r *http.Request) {
	files := []string{"index.html"}
	if r.URL.Path != "/" {
		files = []string{r.URL.Path, path.Join(r.URL.Path, "index.html")}
	}
	if !tryFiles(files, w, r) {
		http.NotFound(w, r)
	}
}

func tryFiles(files []string, w http.ResponseWriter, r *http.Request) bool {
	for _, file := range files {
		realpath := filepath.Join(doc, file)
		if f, err := os.Stat(realpath); err == nil {
			if f.Mode().IsRegular() {
				http.ServeFile(w, r, realpath)
				return true
			}
		}
	}
	return false
}
