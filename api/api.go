package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/opencap/go-server/database"
)

type config struct {
	db                database.Database
	jwtExpirationTime time.Duration
	jwtSecret         string
	domain            string
}

func (cfg *config) initDB() error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return errors.New("No DB_URL found in env")
	}
	dbType := os.Getenv("DB_TYPE")
	if dbURL == "" {
		return errors.New("No DB_TYPE found in env")
	}

	var err error
	switch dbType {
	case "postgres":
		cfg.db, err = database.GetGormConnection(dbURL, dbType)
	default:
		err = errors.New("Invalid DB_TYPE specified in env")
	}
	if err != nil {
		return err
	}

	if os.Getenv("PLATFORM_ENV") == "test" {
		cfg.db.CreateTables(true)
	} else {
		cfg.db.CreateTables(false)
	}

	return nil
}

func (cfg *config) initAuth() error {
	jwtExpirationMinutesString := os.Getenv("JWT_EXPIRATION_MINUTES")
	jwtExpirationMinutes, err := strconv.Atoi(jwtExpirationMinutesString)
	if err != nil || jwtExpirationMinutes < 1 {
		return errors.New("JWT_EXPIRATION_MINUTES must be greater than 0")
	}
	cfg.jwtExpirationTime = time.Duration(jwtExpirationMinutes) * time.Minute

	cfg.jwtSecret = os.Getenv("JWT_SECRET")
	if len(cfg.jwtSecret) < 1 {
		return errors.New("JWT_SECRET must be greater than 0")
	}
	return nil
}

func (cfg *config) initDomain() error {
	cfg.domain = os.Getenv("DOMAIN")
	if len(cfg.domain) < 1 {
		return errors.New("DOMAIN is missing from env")
	}
	return nil
}

// Start begins serving the API
func Start() *http.Server {
	cfg := config{}
	err := cfg.initDB()
	if err != nil {
		panic(err.Error())
	}
	err = cfg.initAuth()
	if err != nil {
		panic(err.Error())
	}
	err = cfg.initDomain()
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()

	// CAP
	r.HandleFunc("/v1/addresses", cfg.getAddressHandler).Methods("GET")

	// CAMP
	r.HandleFunc("/v1/auth", cfg.postAuthHandler).Methods("POST")
	r.HandleFunc("/v1/addresses", cfg.putAddressHandler).Methods("PUT")
	r.HandleFunc("/v1/users", cfg.deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/v1/addresses/{address_type}", cfg.deleteAddressesHandler).Methods("DELETE")
	r.HandleFunc("/v1/users", cfg.postUserHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		panic("No PORT specified in the environment")
	}
	timeoutSecondsString := os.Getenv("TIMEOUT_SECONDS")
	if timeoutSecondsString == "" {
		panic("No TIMEOUT_SECONDS specified in the environment")
	}
	timeoutSeconds, err := strconv.Atoi(timeoutSecondsString)
	if err != nil {
		panic("TIMEOUT_SECONDS must be an integer")
	}

	server := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * time.Duration(timeoutSeconds),
		ReadTimeout:  time.Second * time.Duration(timeoutSeconds),
		IdleTimeout:  time.Second * time.Duration(timeoutSeconds),
		Handler:      r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err.Error())
		}
	}()

	fmt.Println("Listening for requests on " + server.Addr)
	return server
}
