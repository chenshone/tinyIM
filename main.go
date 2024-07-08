package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tinyIM/api"
)

func main() {
	var module string
	flag.StringVar(&module, "module", "", "assign run module")
	flag.Parse()
	fmt.Printf("start run %s module\n", module)

	switch module {
	case "api":
		api.New().Run()
	default:
		fmt.Printf("exiting...\nmodule param error!\n")
		return
	}
	fmt.Printf("run %s module done!", module)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
