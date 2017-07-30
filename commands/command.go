package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"netdisk/config"
	"netdisk/layers/fslayer"
	"netdisk/util"
	"os"
	"path/filepath"
	"strconv"

	"netdisk/layers/netlayer"

	"github.com/suconghou/utilgo"
	kcp "github.com/xtaci/kcp-go"
)

// Use choose a backend
func Use() {
	if len(os.Args) == 3 {
		err := config.Use(os.Args[2])
		if err != nil {
			util.Log.Printf("%v", err)
		}
	} else {
		Pwd()
	}
}

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
			util.Log.Printf("%v", err)
		}
	} else {
		util.Log.Print("Usage:disk cd newpath")
	}
}

// Pwd print current work dir
func Pwd() {
	err := fslayer.Pwd()
	if err != nil {
		util.Log.Printf("%v", err)
	}
}

// Mv move file from the backend
func Mv() {
	if len(os.Args) == 4 {
		err := fslayer.MoveFile(os.Args[2], os.Args[3])
		if err != nil {
			util.Log.Printf("%v", err)
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
			util.Log.Printf("%v", err)
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
			util.Log.Printf("%v", err)
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
			util.Log.Printf("%v", err)
		}
	} else {
		util.Log.Print("Usage:disk rm filepath")
	}
}

// Get do a simple download
func Get() {
	if len(os.Args) >= 3 && !utilgo.IsURL(os.Args[2]) {
		reqHeader := netlayer.ParseCookieUaRefer()
		thread, thunk, start, end := netlayer.ParseThreadThunkStartEnd(8, 2097152, -1, 0)
		saveas, err := utilgo.GetStorePath(os.Args[2])
		if err != nil {
			util.Log.Printf("%v", err)
			return
		}
		err = fslayer.Get(os.Args[2], saveas, reqHeader, thread, thunk, start, end)
		if err != nil {
			util.Log.Printf("%v", err)
		}
	} else {
		util.Log.Print("Usage:disk get filepath")
	}
}

// Put upload file to the backend
func Put() {
	if len(os.Args) >= 3 {

	} else {
		util.Log.Print("Usage:disk put filepath")
	}
}

// Wget url like wget
func Wget() {
	if len(os.Args) >= 3 && utilgo.IsURL(os.Args[2]) {
		var (
			reqHeader                 = netlayer.ParseCookieUaRefer()
			thread, thunk, start, end = netlayer.ParseThreadThunkStartEnd(8, 2097152, -1, 0)
			saveas, err               = utilgo.GetStorePath(os.Args[2])
		)
		if err != nil {
			util.Log.Printf("%v", err)
			return
		}
		err = fslayer.WgetURL(os.Args[2], saveas, reqHeader, thread, thunk, start, end)
		if err != nil {
			util.Log.Printf("%v", err)
		}
	} else {
		util.Log.Print("Usage:disk wget url")
	}
}

// Play play a url or file(pcs file)
func Play() {
	if len(os.Args) >= 3 {
		var (
			saveas                    string
			err                       error
			stdout                    = utilgo.HasFlag("--stdout")
			reqHeader                 = netlayer.ParseCookieUaRefer()
			thread, thunk, start, end = netlayer.ParseThreadThunkStartEnd(8, 2097152, -1, 0)
		)
		if !stdout {
			saveas, err = utilgo.GetStorePath(os.Args[2])
			if err != nil {
				util.Log.Printf("%v", err)
				return
			}
		} else {
			util.Log.SetOutput(os.Stderr)
		}
		util.Log.Print("Playing " + saveas)
		if utilgo.IsURL(os.Args[2]) {
			err = fslayer.PlayURL(os.Args[2], saveas, reqHeader, thread, thunk, start, end, stdout)
			if err != nil {
				util.Log.Printf("%v", err)
			}
		} else {
			err = fslayer.Play(os.Args[2], saveas, reqHeader, thread, thunk, start, end, stdout)
			if err != nil {
				util.Log.Printf("%v", err)
			}
		}
	} else {
		util.Log.Print("Usage:disk play filepath/url")
	}
}

// Sync files
func Sync() {

}

// Info print the backend info or file info
func Info() {
	if len(os.Args) >= 3 {
		err := fslayer.GetFileInfo(os.Args[2], utilgo.HasFlag("--link"))
		if err != nil {
			util.Log.Printf("%v", err)
		}
	} else {
		err := fslayer.GetInfo()
		if err != nil {
			util.Log.Printf("%v", err)
		}
	}
}

// Hash print the md5
func Hash() {
	if len(os.Args) == 3 {
		var filePath = (os.Args[2])
		util.PrintMd5(filePath)
	} else {
		fmt.Println("Usage:disk hash file")
	}
}

// Help print the help message
func Help() {
	fmt.Println(os.Args[0] + " ls info mv cp get put wget play rm mkdir pwd hash config empty search task ")
}

// Config set or get the app config
func Config() {
	if (len(os.Args) == 2) || (os.Args[2] == "list") {
		config.ConfigList()
	} else if len(os.Args) == 3 && os.Args[2] == "get" {
	} else if len(os.Args) == 4 && os.Args[2] == "set" {
	} else if len(os.Args) == 4 && os.Args[2] == "setapp" {
	} else {
		fmt.Println("Usage:disk config list/get/set/setapp")
	}
}

// Task list current backend task
func Task() {
	if (len(os.Args) == 2) || (os.Args[2] == "list") {
		fslayer.GetTaskList()
	} else if len(os.Args) == 5 && os.Args[2] == "add" {
		fslayer.AddTask((os.Args[3]), os.Args[4])
	} else if len(os.Args) == 4 && os.Args[2] == "remove" {
		fslayer.RemoveTask(os.Args[3])
	} else if len(os.Args) == 4 && os.Args[2] == "info" {
		fslayer.GetTaskInfo(os.Args[3])
	} else {
		fmt.Println("Usage:disk task list/add/info/remove")
	}
}

// Search form the backend
func Search() {
	if len(os.Args) == 3 {
		err := fslayer.SearchFile(os.Args[2])
		if err != nil {
			util.Log.Printf("%v", err)
		}
	}
}

// Empty clear cache data
func Empty() {
	fslayer.Empty()
}

func servePreCheck(callback func(int, string) error) {
	var (
		port int
		root string
	)
	var ferr flag.ErrorHandling
	var CommandLine = flag.NewFlagSet(os.Args[0], ferr)
	CommandLine.IntVar(&port, "p", 6060, "http server port")
	CommandLine.StringVar(&root, "d", "./", "root document dir")
	err := CommandLine.Parse(os.Args[2:])
	if err != nil {
		os.Exit(2)
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var doc string
	if filepath.IsAbs(root) {
		doc = root
	} else {
		doc = filepath.Join(pwd, root)
	}
	if _, err := os.Stat(doc); err == nil {
		fmt.Println("Server listening on port " + strconv.Itoa(port))
		fmt.Println("Document root " + doc)
		err = callback(port, doc)
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println(doc + " not exists")
	}
}

// Serve start a http server
func Serve() {
	servePreCheck(func(port int, doc string) error {
		return http.ListenAndServe(":"+strconv.Itoa(port), http.FileServer(http.Dir(doc)))
	})
}

// KcpServe start an kcp server
func KcpServe() {
	if len(os.Args) == 3 {
		kcpconn, err := kcp.DialWithOptions(os.Args[2], nil, 10, 3)
		if err != nil {
			log.Fatal(err)
		}
		kcpconn.SetStreamMode(true)
		kcpconn.SetWriteDelay(true)
		kcpconn.SetNoDelay(1, 10, 2, 1)
		kcpconn.SetWindowSize(1024, 1024)
		kcpconn.SetMtu(1350)
		kcpconn.SetACKNoDelay(true)
		if err := kcpconn.SetDSCP(0); err != nil {
			log.Println("SetDSCP:", err)
		}
		if err := kcpconn.SetReadBuffer(4194304); err != nil {
			log.Println("SetReadBuffer:", err)
		}
		if err := kcpconn.SetWriteBuffer(4194304); err != nil {
			log.Println("SetWriteBuffer:", err)
		}
		if _, err := io.Copy(os.Stdout, kcpconn); err != nil {
			log.Fatal(err)
		}
		kcpconn.Close()

	} else {
		servePreCheck(func(port int, doc string) error {
			lis, err := kcp.ListenWithOptions("0.0.0.0:"+strconv.Itoa(port), nil, 10, 3)
			if err != nil {
				return err
			}
			for {
				if conn, err := lis.AcceptKCP(); err == nil {
					fmt.Println(conn)
					conn.SetStreamMode(true)
					conn.SetWriteDelay(true)
					conn.SetNoDelay(1, 10, 2, 1)
					conn.SetMtu(1350)
					conn.SetWindowSize(1024, 1024)
					conn.SetACKNoDelay(true)
					if err := lis.SetDSCP(0); err != nil {
						log.Println("SetDSCP:", err)
					}
					if err := lis.SetReadBuffer(4194304); err != nil {
						log.Println("SetReadBuffer:", err)
					}
					if err := lis.SetWriteBuffer(4194304); err != nil {
						log.Println("SetWriteBuffer:", err)
					}
					go handleClient(conn, doc)
				} else {
					log.Printf("Error %+v", err)
				}
			}

		})
	}
}

func handleClient(conn io.ReadWriteCloser, file string) {
	f, err := os.OpenFile(file, os.O_RDWR, 0755)
	fmt.Println(f)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(conn, f); err != nil {
		log.Fatal(err)
	}
	conn.Close()
}

// Usage print help message
func Usage() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println(os.Args[0] + " version: disk/" + config.Version + "\n" + config.ReleaseURL)
	} else {
		Help()
	}
}
