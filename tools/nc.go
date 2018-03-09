package tools

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/suconghou/utilgo"
)

type tcpCounter struct {
	*net.TCPConn
	startTime time.Time
	totalr    uint64
	totalw    uint64
	readed    uint64
	writed    uint64
}

func progress(totalr uint64, readed uint64, totalw uint64, writed uint64, duration float64) {
	sends := float64(writed) / 1024 / duration
	receives := float64(readed) / 1024 / duration
	var (
		s1 string
		s2 string
	)
	if totalw > 0 {
		s1 = fmt.Sprintf("/%s/%d%%", utilgo.ByteFormat(totalw), 100*writed/totalw)
	}
	if totalr > 0 {
		s2 = fmt.Sprintf("/%s/%d%%", utilgo.ByteFormat(totalr), 100*readed/totalr)
	}
	fmt.Fprintf(os.Stderr, "\r\033[2K\r发送:%.2fKB/s|%s%s 接收:%.2fKB/s|%s%s", sends, utilgo.ByteFormat(writed), s1, receives, utilgo.ByteFormat(readed), s2)
}

func (c *tcpCounter) Read(p []byte) (int, error) {
	n, err := c.TCPConn.Read(p)
	if err != nil {
		return n, err
	}
	c.readed += uint64(n)
	progress(c.totalr, c.readed, c.totalw, c.writed, time.Since(c.startTime).Seconds())
	return n, err
}

func (c *tcpCounter) Write(p []byte) (int, error) {
	n, err := c.TCPConn.Write(p)
	if err != nil {
		return n, err
	}
	c.writed += uint64(n)
	progress(c.totalr, c.readed, c.totalw, c.writed, time.Since(c.startTime).Seconds())
	return n, err
}

// NcMain start
func NcMain() error {
	progress := utilgo.HasFlag("-v")
	if len(os.Args) > 3 && os.Args[2] == "-l" && utilgo.IsPort(os.Args[3]) {
		return ncServer(os.Args[3], progress)
	} else if len(os.Args) > 2 && utilgo.IsIPPort(os.Args[2]) {
		return ncClient(os.Args[2], progress)
	}
	return errArgs
}

func ncServer(port string, progress bool) error {
	l, err := net.Listen(tcp, ":"+port)
	if err != nil {
		return err
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	if progress {
		c := &tcpCounter{TCPConn: conn.(*net.TCPConn), startTime: time.Now(), totalw: uint64(123455)}
		go func() { io.Copy(c, os.Stdin); c.CloseWrite(); wg.Done() }()
		go func() { io.Copy(os.Stdout, c); c.CloseRead(); wg.Done() }()
	} else {
		c := conn.(*net.TCPConn)
		go func() { io.Copy(c, os.Stdin); c.CloseWrite(); wg.Done() }()
		go func() { io.Copy(os.Stdout, c); c.CloseRead(); wg.Done() }()
	}
	wg.Wait()
	return nil
}

func ncClient(addr string, progress bool) error {
	conn, err := net.Dial(tcp, addr)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	if progress {
		c := &tcpCounter{TCPConn: conn.(*net.TCPConn), startTime: time.Now(), totalw: uint64(123455)}
		go func() { io.Copy(c, os.Stdin); c.CloseWrite(); wg.Done() }()
		go func() { io.Copy(os.Stdout, c); c.CloseRead(); wg.Done() }()
	} else {
		c := conn.(*net.TCPConn)
		go func() { io.Copy(c, os.Stdin); c.CloseWrite(); wg.Done() }()
		go func() { io.Copy(os.Stdout, c); c.CloseRead(); wg.Done() }()
	}
	wg.Wait()
	return nil
}
