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

const tcp = "tcp"

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

// NcMain start
func NcMain() error {
	if len(os.Args) > 3 && os.Args[2] == "-l" && utilgo.IsPort(os.Args[3]) {
		return ncServer(os.Args[3])
	} else if len(os.Args) > 2 && utilgo.IsIPPort(os.Args[2]) {
		return ncClient(os.Args[2])
	}
	return fmt.Errorf("args error")
}

func ncServer(port string) error {
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
	// con := &counter{origin: conn, startTime: time.Now(), totalw: uint64(10243246)}
	go func() { io.Copy(conn, os.Stdin); fmt.Fprintf(os.Stderr, "os.Stdin 2 conn done"); wg.Done() }()
	go func() { io.Copy(os.Stdout, conn); fmt.Fprintf(os.Stderr, "conn 2 os.Stdout done"); wg.Done() }()
	wg.Wait()
	conn.Close()
	return nil
}

func ncClient(addr string) error {
	conn, err := net.Dial(tcp, addr)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { io.Copy(conn, os.Stdin); fmt.Fprintf(os.Stderr, "conn 2 os.Stdin done"); wg.Done() }()
	go func() { io.Copy(os.Stdout, conn); fmt.Fprintf(os.Stderr, "os.Stdout 2 conn done"); wg.Done() }()
	wg.Wait()
	conn.Close()
	return nil
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
