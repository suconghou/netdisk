package config

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	s := `Disk: a remote disk for sync files and cloud storage
Usage: disk Command [Options...]
Commands:
    cd                      Change current directory.
    connect                 Connect to the server.
    get                     Download file from the server.
    info                    Get the system information.
    list                    Get the task information.
    ls                      List files in current directory.
    mkdir                   Create a new directory.
    put                     Upload file to the server.
    rm                      Delete file from the server.
    server                  Run in server mode.
    sync                    Sync the files from server.
    wget                    Download the file in the server.
Options:
    -r, -root=<path>        Root path of the site. Default is current working directory.
    -p, -port=<port>        HTTP port. Default is 8080.
        -404=<path>         Path of a custom 404 file, relative to Root. Example: /404.html.
    -g, -gzip=<bool>        Turn on or off gzip compression. Default value is true (means turn on).
    -a, -auth=<user:pass>   Turn on digest auth and set username and password (separate by colon).
                            After turn on digest auth, all the page require authentication.
        -401=<path>         Path of a custom 401 file, relative to Root. Example: /401.html.
                            If authentication fails and 401 file is set,
                            the file content will be sent to the client.
                            If use with -make-cert, will generate a certificate to the path.
        -key=<path>         Load a file as a private key.
Other options:
    -v, -version            Show version information.
    -h, -help               Show help message.
Author:
    suconghou
    <https://github.com/suconghou>
    <http://blog.suconghou.cn>
`
	fmt.Printf(s)
	os.Exit(0)

}

func LoadConfig() {

	var key string
	var port uint
	flag.StringVar(&key, "k", "", "Password")
	flag.UintVar(&port, "p", 0, "HTTP port")
	flag.UintVar(&port, "port", 0, "HTTP port")
	flag.Parse()

}
