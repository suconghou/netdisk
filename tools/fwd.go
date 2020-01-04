package tools

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/suconghou/netdisk/util"
	"github.com/suconghou/utilgo"
)

const (
	tryTimes = 5
)

// FwdMain run port forward
func FwdMain() error {
	if len(os.Args) > 3 {
		if utilgo.IsIPPort(os.Args[2]) && utilgo.IsIPPort(os.Args[3]) {
			proto := utilgo.BoolString(utilgo.HasFlag(os.Args, "-u"), udp, tcp)
			l, err := net.Listen(proto, os.Args[2])
			if err != nil {
				return err
			}
			for {
				conn, err := l.Accept()
				if err != nil {
					return err
				}
				go streamCopy(conn, os.Args[3], proto)
			}
		}
	}
	return errArgs
}

func streamCopy(conn net.Conn, remote string, proto string) {

	var (
		c   net.Conn
		t   int
		err error
		wg  sync.WaitGroup
	)

	for {
		fmt.Print(proto, remote)
		c, err = net.Dial(proto, remote)
		t++
		if err != nil {
			util.Debug.Println(err, t)
			if t > tryTimes {
				break
			}
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	if err == nil {
		wg.Add(2)
		go func() { io.Copy(conn, c); wg.Done() }()
		go func() { io.Copy(c, conn); wg.Done() }()
		wg.Wait()
		conn.Close()
		c.Close()
	} else {
		util.Debug.Println(err)
		conn.Close()
	}

}
