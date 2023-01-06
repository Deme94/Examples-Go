package database

import (
	"API-REST/services/database/mongo"
	"API-REST/services/database/postgres"

	_ "github.com/lib/pq"
)

func Setup() error {
	err := postgres.Setup()
	if err != nil {
		return err
	}
	err = mongo.Setup()
	if err != nil {
		return err
	}

	return nil
}

// Create DB and tables
func Init() error {
	return postgres.Init()
}
