package main

import (
	"github.com/opencap/opencap/internal/app/opencap"
	"github.com/opencap/opencap/internal/pkg/database"
	"log"
	"os"
)

const (
	host = ""
	port = 41145
)

func main() {
	log.Print("Opening database")
	db, err := database.NewSQLiteDatabase("file:opencap.db")
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	log.Print("Creating server")
	server := &opencap.Server{}
	server.SetDatabase(db)

	log.Printf("Listening on %s:%d", host, port)
	if err := server.Run(host, port); err != nil {
		panic(err)
	}
}
