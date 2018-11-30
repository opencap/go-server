package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	opencap "github.com/opencap/go-opencap"
	"github.com/opencap/go-server/database"
	"golang.org/x/crypto/acme/autocert"
)

// Config represents the configuration of this API
type Config struct {
	db                 database.Database
	jwtExpirationTime  time.Duration
	jwtSecret          string
	createUserPassword string
	domainName         string
}

// InitDB get a connection to the database
func (cfg *Config) InitDB() error {
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
	case "sqlite3":
		cfg.db, err = database.GetGormConnection(dbURL, dbType)
	case "mssql":
		cfg.db, err = database.GetGormConnection(dbURL, dbType)
	case "mysql":
		cfg.db, err = database.GetGormConnection(dbURL, dbType)
	default:
		err = errors.New("Invalid DB_TYPE specified in env")
	}
	if err != nil {
		return err
	}

	return nil
}

// SetupDB setup the database for the first time
func (cfg *Config) SetupDB() error {
	if os.Getenv("PLATFORM_ENV") == "prod" {
		return cfg.db.CreateTables(false)
	}
	return cfg.db.CreateTables(true)
}

func (cfg *Config) intJWTConfig() error {
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

func (cfg *Config) initAuthPassword() error {
	cfg.createUserPassword = os.Getenv("CREATE_USER_PASSWORD")
	const minLength = 8
	if len(cfg.createUserPassword) < minLength {
		return errors.New("CREATE_USER_PASSWORD must be longer than " + string(minLength) + " characters")
	}
	return nil
}

func (cfg *Config) initDomainName() error {
	cfg.domainName = os.Getenv("DOMAIN_NAME")
	if !opencap.ValidateDomain(cfg.domainName) {
		return errors.New("Invalid DOMAIN_NAME in .env")
	}
	return nil
}

// Start begins serving the API
func Start() *http.Server {
	cfg := Config{}
	if err := cfg.InitDB(); err != nil {
		log.Fatal(err.Error())
	}
	if !cfg.db.HasTables() {
		log.Fatal("Database not setup. Please use the \"--setupdatabase\" option")
	}
	if err := cfg.intJWTConfig(); err != nil {
		log.Fatal(err.Error())
	}
	if err := cfg.initAuthPassword(); err != nil {
		log.Fatal(err.Error())
	}
	if err := cfg.initDomainName(); err != nil {
		log.Fatal(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/addresses", cfg.getAddressHandler).Methods("GET")
	r.HandleFunc("/v1/auth", cfg.postAuthHandler).Methods("POST")
	r.HandleFunc("/v1/addresses", cfg.putAddressHandler).Methods("PUT")
	r.HandleFunc("/v1/users", cfg.deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/v1/addresses/{address_type}", cfg.deleteAddressesHandler).Methods("DELETE")
	r.HandleFunc("/v1/users", cfg.postUserHandler).Methods("POST")

	if os.Getenv("PLATFORM_ENV") == "prod" {
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := cfg.domainName
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache("certs"),
		}

		http.Handle("/", r)
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		go http.Serve(certManager.Listener(), nil)
		fmt.Println("Production OpenCAP server started successfully")
		return &http.Server{} // Don't use a server struct for production
	}

	testPort := os.Getenv("TEST_PORT")
	if testPort == "" {
		log.Fatal("No PORT specified in the environment")
	}
	server := http.Server{
		Addr:    ":" + testPort,
		Handler: r,
	}
	go server.ListenAndServe()

	fmt.Println("Listening for requests on localhost:" + testPort)
	return &server
}
