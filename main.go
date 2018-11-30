package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/opencap/go-server/configure"

	"github.com/joho/godotenv"
	"github.com/opencap/go-server/api"
)

func main() {
	godotenv.Load()

	openPort := flag.String("openport", "", "Open the PORT from your router to this device")
	closePort := flag.String("closeport", "", "Close the PORT from your router to this device")
	getIP := flag.Bool("getip", false, "Print out the public IP address of this machine")
	setupDatabase := flag.Bool("setupdatabase", false, "Setup the database for the first time")
	flag.Parse()

	if *openPort != "" {
		err := configure.OpenPort(*openPort)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Port %v forwarded successfully", *openPort)
		os.Exit(0)
	}

	if *closePort != "" {
		err := configure.ClosePort(*closePort)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("Port %v closed successfully", *closePort)
		os.Exit(0)
	}

	if *getIP {
		ip, err := configure.GetPublicIP()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("Your IP address is: " + ip)
		os.Exit(0)
	}

	if *setupDatabase {
		cfg := api.Config{}
		if err := cfg.InitDB(); err != nil {
			log.Fatal(err.Error())
		}
		if err := cfg.SetupDB(); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("Database has been setup correctly")
		os.Exit(0)
	}

	server := api.Start()
	defer server.Shutdown(nil)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
