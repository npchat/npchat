package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/intob/npchat/cfg"
	"github.com/intob/npchat/server"
)

func main() {
	cfg.InitViper()

	ws := server.NewWebServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	<-c
	fmt.Printf("\r\nshutting down...\r\n")
	ws.Shutdown()
}
