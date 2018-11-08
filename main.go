package main

import "github.com/opencap/go-server/api"
import "github.com/joho/godotenv"

func main() {
	godotenv.Load()

	server := api.Start()
	defer server.Shutdown(nil)

	forever := make(chan bool)
	<-forever
}
