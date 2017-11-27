package tools

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/suconghou/utilgo"

	kcp "github.com/xtaci/kcp-go"
)

type counter struct {
	startTime time.Time
	origin    io.ReadWriter
	totalr    uint64
	totalw    uint64
	readed    uint64
	writed    uint64
}

func progress(totalr uint64, readed uint64, totalw uint64, writed uint64, duration float64) {
	writer := os.Stderr
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
	fmt.Fprintf(writer, "\r\033[2K\r发送:%.2fKB/s|%s%s 接收:%.2fKB/s|%s%s", sends, utilgo.ByteFormat(writed), s1, receives, utilgo.ByteFormat(readed), s2)
}

func (c *counter) Read(p []byte) (int, error) {
	n, err := c.origin.Read(p)
	if err != nil {
		return n, err
	}
	c.readed += uint64(n)
	progress(c.totalr, c.readed, c.totalw, c.writed, time.Since(c.startTime).Seconds())
	return n, err
}

func (c *counter) Write(p []byte) (int, error) {
	n, err := c.origin.Write(p)
	if err != nil {
		return n, err
	}
	c.writed += uint64(n)
	progress(c.totalr, c.readed, c.totalw, c.writed, time.Since(c.startTime).Seconds())
	return n, err
}

// NcTCP is tcp nc like but with progress bar
func NcTCP(address string, port int, serve bool, prog bool, act string, arg string) error {
	var co io.ReadWriter
	if serve {
		l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
		if err != nil {
			return err
		}
		defer l.Close()
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		if prog {
			info, err := os.Stdin.Stat()
			if err != nil {
				return err
			}
			co = &counter{origin: conn, startTime: time.Now(), totalw: uint64(info.Size())}
		} else {
			co = conn
		}
		p1die := make(chan bool)
		go func(c io.ReadWriter) {
			io.Copy(c, os.Stdin)
			close(p1die)
		}(co)
		p2die := make(chan bool)
		go func(c io.ReadWriter) {
			io.Copy(os.Stdout, c)
			close(p2die)
		}(co)
		select {
		case <-p1die:
			<-p2die
		case <-p2die:
			<-p1die
		}
		return nil
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return err
	}
	if prog {
		info, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		co = &counter{origin: conn, startTime: time.Now(), totalw: uint64(info.Size())}
	} else {
		co = conn
	}
	p1die := make(chan bool)
	go func(c io.ReadWriter) {
		io.Copy(c, os.Stdin)
		close(p1die)
	}(co)
	p2die := make(chan bool)
	go func(c io.ReadWriter) {
		io.Copy(os.Stdout, c)
		close(p2die)
	}(co)
	select {
	case <-p1die:
		<-p2die
	case <-p2die:
		<-p1die
	}
	return nil
}

// Nc1 like but use kcp to transfer data
func Nc1(address string, port int, act string, arg string) error {

	const (
		dataShard    = 10
		parityShard  = 3
		noDelay      = 1
		interval     = 20
		resend       = 2
		noCongestion = 1
	)
	if act == "" {
		lis, err := kcp.ListenWithOptions(fmt.Sprintf("%s:%d", address, port), nil, dataShard, parityShard)
		if err != nil {
			return err
		}
		conn, err := lis.AcceptKCP()
		if err != nil {
			return err
		}
		conn.SetStreamMode(true)
		conn.SetWriteDelay(true)
		conn.SetNoDelay(noDelay, interval, resend, noCongestion)
		conn.SetWindowSize(1024, 1024)
		conn.SetMtu(1350)
		conn.SetACKNoDelay(true)
		if err != nil {
			return err
		}

	} else {
		conn, err := kcp.DialWithOptions(fmt.Sprintf("%s:%d", address, port), nil, dataShard, parityShard)
		if err != nil {
			return err
		}
		conn.SetStreamMode(true)
		conn.SetWriteDelay(true)
		conn.SetNoDelay(noDelay, interval, resend, noCongestion)
		conn.SetWindowSize(1024, 1024)
		conn.SetMtu(1350)
		conn.SetACKNoDelay(true)
	}
	return nil
}
