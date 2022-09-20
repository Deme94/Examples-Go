package main

import (
	"log"
	"os"

	"github.com/arthurkushman/buildsqlx"
	_ "github.com/lib/pq"
)

// Database Postgres SQL
var postgresDatasource = os.Getenv("POSTGRES_DATASOURCE")

type databases struct {
	postgresql *buildsqlx.DB
	//mongodb    *mongoDbType
	// ...
}

func (s *server) connectDatabases() *databases {
	dbs := databases{
		postgresql: buildsqlx.NewDb(buildsqlx.NewConnection("postgres", postgresDatasource)),
	}
	// Check connections
	err := dbs.postgresql.Sql().Ping()
	if err != nil {
		log.Println("postgres", postgresDatasource)
		log.Fatalf("Database connection failed -> %v", err)
	}

	return &dbs
}
