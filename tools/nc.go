package tools

import (
	"fmt"

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
