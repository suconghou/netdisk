package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/suconghou/fastload/fastload"
	"github.com/suconghou/netdisk/util"
	"golang.org/x/net/proxy"
)

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
	_, err := fastload.Pipe(w, cleanHeader(r, xheaders), url, usecachefilter, 3600, nil)
	if err != nil {
		util.Log.Printf("pipe %s error:%s", url, err)
	}
}

// Proxy is a http_proxy and just http_proxy server
func Proxy(w http.ResponseWriter, r *http.Request) error {
	return fastload.HTTPProxy(w, r)
}

// ProxySocks is a http_proxy https_proxy socks proxy server
func ProxySocks(client net.Conn, dialer proxy.Dialer) error {
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		return err
	}
	if b[0] == 5 { // socks5
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		if err != nil {
			return err
		}
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
		server, err := dialer.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			return err
		}
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		p1die := make(chan bool)
		go func() { io.Copy(server, client); close(p1die) }()
		p2die := make(chan bool)
		go func() { io.Copy(client, server); close(p2die) }()
		select {
		case <-p1die:
		case <-p2die:
		}
		return nil
	}
	// try https_proxy and http_proxy
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		return err
	}
	if method == "CONNECT" { // https_proxy
		_, err = client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			return err
		}
		address = fmt.Sprintf("%s", hostPortURL)
	} else { // at last parse as http_proxy
		address = hostPortURL.Host
		if strings.Index(address, ":") == -1 {
			address = address + ":80"
		}
	}
	server, err := dialer.Dial("tcp", address)
	if err != nil {
		return err
	}
	if method != "CONNECT" {
		_, err = server.Write(b[:n])
		if err != nil {
			return err
		}
	}
	p1die := make(chan bool)
	go func() { io.Copy(server, client); close(p1die) }()
	p2die := make(chan bool)
	go func() { io.Copy(client, server); close(p2die) }()
	select {
	case <-p1die:
	case <-p2die:
	}
	return nil
}
