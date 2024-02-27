package main

import (
	"os"
	"os/signal"
	"syscall"

	"workspace/fileserver"
	"workspace/idp"
	"workspace/proxy"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	go proxy.StartProxy()

	go idp.StartIDPServer()

	go fileserver.StartFileHostingServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
