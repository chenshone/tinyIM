package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tinyIM/api"
	"tinyIM/logic"
)

func main() {
	var module string
	flag.StringVar(&module, "module", "", "assign run module")
	flag.Parse()
	fmt.Printf("start run %s module\n", module)

	switch module {
	case "logic":
		logic.New().Run()
	case "api":
		api.New().Run()
	default:
		fmt.Printf("exiting...\nmodule param error!\n")
		return
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Printf("%s module exit", module)
}
