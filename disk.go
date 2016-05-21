package main

import (
	"layers/fslayer"
	"layers/netlayer"
	"server"
	"fmt"
)

func main() {
	fslayer.Get()
	netlayer.Get()
	a,_:=server.ListDir("/data/tmp")
	fmt.Println(a)
}
