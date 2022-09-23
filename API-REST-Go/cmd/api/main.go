package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// CONSTANTS ----------------------------------------

// App version
const VERSION = "1.0.0"

var domain = os.Getenv("DOMAIN")
var p, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
var env = os.Getenv("ENVIRONMENT")
var secret = os.Getenv("SERVER_JWT")

// --------------------------------------------------

type config struct {
	port int
	env  string
	jwt  struct {
		secret string
	}
}

type server struct {
	config      config
	logger      *log.Logger
	databases   *databases
	controllers *controllers
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", p, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", env, "Server environment ("+env+")")
	flag.StringVar(&cfg.env, "domain", domain, "Host domain ("+domain+")")
	flag.Parse()

	cfg.jwt.secret = secret

	// Server server
	s := &server{
		config: cfg,
		logger: logger(),
	}
	s.databases = s.connectDatabases()
	s.controllers = s.createControllers(s.databases)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      s.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	srv.SetKeepAlivesEnabled(true)

	s.logger.Println("Started server on port", cfg.port)
	log.Println("Started server on port", cfg.port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
