package main

import (
	"marketplace_websocket/pkg/server"
)

func main() {
	app := &server.App{}
	app.Initialize()
	address := ":3000"
	app.Run(address)
}
