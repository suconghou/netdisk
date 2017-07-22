package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"ipproxypool/util"
	"log"
	"net/http"
	"netdisk/commands"
	"netdisk/config"
	"netdisk/route"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	startTime = time.Now()
	port      string
	doc       string
	cfgPath   string
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
	config.LoadConfig()
	if len(os.Args) > 1 {
		cli()
	} else {
		daemon()
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
		commands.Hash()
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
	case "kcp":
		commands.KcpServe()
	default:
		commands.Usage()
	}
}

func daemon() {
	flag.StringVar(&port, "port", "8090", "give me a port number")
	flag.StringVar(&doc, "doc", "./", "document root dir")
	flag.StringVar(&cfgPath, "config", "/etc/disk.json", "config file path")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !filepath.IsAbs(doc) {
		doc = filepath.Join(pwd, doc)
	}
	if f, err := os.Stat(doc); err == nil {
		if !f.Mode().IsDir() {
			fmt.Println(doc + " is not directory")
			os.Exit(3)
		}
	} else {
		fmt.Println(doc + " not exists")
		os.Exit(2)
	}
	flag.Parse()
	http.HandleFunc("/status", status)
	http.HandleFunc("/", routeMatch)
	fmt.Println("Starting up on port " + port)
	fmt.Println("Document root " + doc)
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	log.Fatal(http.ListenAndServe(bind, nil))
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
		util.JsonPut(w, bs, true, 60)
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
	var files []string
	if r.URL.Path == "/" {
		files = []string{"index.html"}
	} else {
		files = []string{r.URL.Path, filepath.Join(r.URL.Path, "index.html")}
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
