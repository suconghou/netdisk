package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"netdisk/util"
	"strconv"
	"strings"

	"github.com/suconghou/fastload/fastload"
)

type proxy struct{}

// Pipe response stream
func Pipe(w http.ResponseWriter, r *http.Request, match []string) {
	var url string
	if match[1] == "" {
		url = fmt.Sprintf("http:/%s", match[0])
	} else {
		url = strings.Replace(strings.TrimPrefix(match[0], "/"), ":/", "://", 1)
	}
	if r.URL.RawQuery != "" {
		url = url + "?" + r.URL.RawQuery
	}
	_, err := fastload.Pipe(w, r, url, func(out http.Header, res *http.Response) int {
		out.Del("Set-Cookie")
		return res.StatusCode
	}, 3600, nil)
	if err != nil {
		util.Log.Printf("pipe %s error:%s", url, err)
	}
}

// Proxy is a http_proxy https_proxy socks5 server
func Proxy(client net.Conn) error {
	defer client.Close()
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		return err
	}
	fmt.Println(b[0], string(b[:]), n)
	if b[0] == 0x05 { //Socks5协议
		fmt.Println("socks5")
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			return err
		}
		defer server.Close()
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		p1die := make(chan struct{})
		go func() { io.Copy(server, client); close(p1die) }()
		p2die := make(chan struct{})
		go func() { io.Copy(client, server); close(p2die) }()
		select {
		case <-p1die:
		case <-p2die:
		}

	} else if b[0] == 67 { // CONNECT https_proxy

		// if req.Method == "CONNECT" {
		// 	client.Write([]byte("HTTP/1.1 200 Connection established\r\n"))
		// } else {

		// }
		// http.ReadResponse
	} else { // http_proxy
		r := io.MultiReader(bytes.NewReader(b[:]), bufio.NewReader(client))
		req, err := http.ReadRequest(bufio.NewReader(r))
		if err != nil {
			return err
		}
		fmt.Print(req.Header)
		// _, err = fastload.Pipe(w, req, req.RequestURI, nil)
		// if err != nil {
		// 	util.Log.Printf("%s proxy error:%s", url, err)
		// }
	}

	return nil
}

func (p proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		url = r.RequestURI
		err error
	)
	util.Log.Print(r.RequestURI, r.Host, r.Method, r.Host, r.Header, r.URL)
	if r.Method == "CONNECT" {

		w.WriteHeader(200)
		return
	}
	_, err = fastload.Pipe(w, r, url, nil, 3600, nil)
	if err != nil {
		util.Log.Printf("%s proxy error:%s", url, err)
	}
}
