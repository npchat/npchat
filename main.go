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
	// handle config
	cfg.InitViper()

	// start server
	ws := server.NewWebServer()

	// wait for OS signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	<-c

	// exit gracefully
	fmt.Printf("\r\nshutting down...\r\n")
	ws.Shutdown()
}
