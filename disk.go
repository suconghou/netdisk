package main

import (
	"commands"
	"config"
	"fmt"
	"os"
)

func main() {
	config.Cfg = config.LoadConfig()
	if len(os.Args) > 1 {
		cli()
	} else {
		daemon()
	}
}

func cli() {
	switch os.Args[1] {
	case "use":
		commands.Use()
	case "ls":
		commands.Ls()
	case "cd":
		commands.Cd()
	case "pwd":
		commands.Pwd()
	case "mv":
		commands.Mv()
	case "mkdir":
		commands.Mkdir()
	case "rm":
		commands.Rm()
	case "get":
		commands.Get()
	case "put":
		commands.Put()
	case "wget":
		commands.Wget()
	case "sync":
		commands.Sync()
	case "info":
		commands.Info()
	case "hash":
		commands.Hash()
	case "play":
		commands.Play()
	case "help":
		commands.Help()
	case "config":
		commands.Config()
	case "task":
		commands.Task()
	default:
		commands.Usage()
	}
}

func daemon() {
	fmt.Println("daemon")
}
