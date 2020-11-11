package tools

import (
	"io"
	"net/http"
	"os"

	"github.com/suconghou/libtspeed"
)

// SpeedTest for http speed test
func SpeedTest(input string, chunk uint, timeout uint, transport *http.Transport) error {
	var (
		r   io.Reader
		err error
	)
	if input == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(input)
		if err != nil {
			return err
		}
	}
	return libtspeed.Run(r, chunk, timeout, transport)
}

// SpeedTestWithHost for http speed test
func SpeedTestWithHost(input string, host string, path string, https bool, chunk uint, timeout uint, transport *http.Transport) error {
	var (
		r   io.Reader
		err error
	)
	if input == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(input)
		if err != nil {
			return err
		}
	}
	return libtspeed.RunHost(r, host, path, https, chunk, timeout, transport)
}
