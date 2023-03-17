package database

import (
	"API-REST/services/database/mongo"
	"API-REST/services/database/postgres"

	_ "github.com/lib/pq"
)

func SetupPostgres() error {
	return postgres.Setup()
}
func SetupMongo() error {
	return mongo.Setup()
}

func SetupPostgresDockertest() error {
	return postgres.SetupDockertest()
}

// Create DB and tables
func Init() error {
	return postgres.Init()
}
