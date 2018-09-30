package main

import "github.com/opencap/go-server/api"

func main() {
	server := api.Start()
	defer server.Shutdown(nil)

	forever := make(chan bool)
	<-forever
}
