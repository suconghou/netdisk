package util

import (
	"log"
	"os"
)

// Log is a global logger
var Log = log.New(os.Stdout, "", 0)

// Debug log to stderr
var Debug = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func LE(err error) {
	if err != nil {
		Log.Println(err)
	}
}
