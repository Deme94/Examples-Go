package main

import (
	"context"
	"log"
	"os"

	"github.com/arthurkushman/buildsqlx"
	_ "github.com/lib/pq"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database Postgres SQL
var postgresURI = os.Getenv("POSTGRES_DATASOURCE")

type databases struct {
	postgresql *buildsqlx.DB
	mongodb    *mongo.Database
	// ...
}

func (s *server) connectDatabases() *databases {
	mongodbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println("postgres", postgresURI)
		log.Fatalf("Database connection failed -> %v", err)
	}

	// Build dbs
	mongodb := mongodbClient.Database("test")
	dbs := databases{
		postgresql: buildsqlx.NewDb(buildsqlx.NewConnection("postgres", postgresURI)),
		mongodb:    mongodb,
	}
	// Check connections
	err = dbs.postgresql.Sql().Ping()
	if err != nil {
		log.Println("postgres", postgresURI)
		log.Fatalf("Database connection failed -> %v", err)
	}
	err = mongodbClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("postgres", postgresURI)
		log.Fatalf("Database connection failed -> %v", err)
	}

	return &dbs
}
