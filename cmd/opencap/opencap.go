package main

import (
	"fmt"
	"github.com/opencap/opencap/internal/app/opencap"
	"github.com/opencap/opencap/internal/pkg/config"
	"github.com/opencap/opencap/internal/pkg/database"
	"log"
	"os"
)

const configPath = "opencap.json"

func main() {
	var (
		server = &opencap.Server{}
		conf   config.Config
		db     database.Database
		err    error
	)

	log.Print("Reading config")
	conf, err = config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}
	server.SetConfig(conf)

	log.Print("Opening database")
	if conf.DatabaseType() == "sqlite" {
		db, err = database.NewSQLiteDatabase(conf.DatabaseDataSource())
	} else {
		err = fmt.Errorf("invalid database type")
	}

	if err != nil {
		log.Printf("Failed to open database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	server.SetDatabase(db)

	log.Printf("Running server")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
