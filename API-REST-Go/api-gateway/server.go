package api

import (
	"log"
	"net/http"
	"time"

	"API-REST/api-gateway/controllers"
	"API-REST/services/conf"
	"API-REST/services/logger"

	_ "github.com/lib/pq"
)

func Setup() error {
	controllers.Build()

	port := conf.Env.GetString("PORT")
	srv := &http.Server{
		Addr:         conf.Env.GetString("HOST") + ":" + port,
		Handler:      Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	srv.SetKeepAlivesEnabled(true)

	logger.Logger.Println("Starting server on port", port)
	log.Println("Starting server on port", port)

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
