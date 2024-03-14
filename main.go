package main

import "ddz/wsserver"

func main() {
	server := wsserver.NewWsServer()
	server.Start()
}
