package tools

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/suconghou/utilgo"

	kcp "github.com/xtaci/kcp-go"
)

// Nc like but use kcp to transfer data
func Nc(address string, port int, act string, arg string) error {
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
		total := 32
		count := 0
		p := make([]byte, 1)
		for {
			n, err := conn.Read(p)
			if err != nil {
				return err
			}
			count += n
			if count == total {
				break
			}
		}
		conn.SetStreamMode(true)
		conn.SetWriteDelay(true)
		conn.SetNoDelay(noDelay, interval, resend, noCongestion)
		conn.SetWindowSize(1024, 1024)
		conn.SetMtu(1350)
		conn.SetACKNoDelay(true)
		p1die := make(chan struct{})
		go func() { io.Copy(os.Stdout, conn); close(p1die) }()
		p2die := make(chan struct{})
		go func() { io.Copy(conn, os.Stdin); close(p2die) }()
		select {
		case <-p1die:
			<-p2die
		case <-p2die:
			<-p1die
		}
		conn.Close()
		lis.Close()
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
		if act == "put" {
			c := "put " + arg
			_, err := conn.Write([]byte(c + strings.Repeat(" ", 32-len(c))))
			if err != nil {
				return err
			}
			f, err := utilgo.GetOpenFile(arg)
			if err != nil {
				return err
			}
			io.Copy(conn, f)
			f.Close()
			conn.Close()
		} else if act == "get" {
			c := "get " + arg
			_, err := conn.Write([]byte(c + strings.Repeat(" ", 32-len(c))))
			if err != nil {
				return err
			}
			p, err := utilgo.GetStorePath(arg)
			if err != nil {
				return err
			}
			f, _, err := utilgo.GetContinue(p)
			io.Copy(f, conn)
		} else {
			c := "ls"
			_, err := conn.Write([]byte(c + strings.Repeat(" ", 32-len(c))))
			if err != nil {
				return err
			}
			io.Copy(os.Stdout, conn)
			conn.Close()
		}
	}
	return nil
}
