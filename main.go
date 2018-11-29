package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/opencap/go-server/api"
	portforward "github.com/opencap/go-server/port-forward"
)

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	port, err := strconv.Atoi(portString)
	if err != nil || port < 1 {
		fmt.Println("PORT must be greater than 0")
		os.Exit(1)
	}

	ip, err := portforward.Open(uint16(port))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Opened port " + string(port) + " with ip " + ip)
	}
	defer portforward.Close(uint16(port))

	server := api.Start()
	defer server.Shutdown(nil)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
