package database

import (
	"API-REST/services/conf"
	"API-REST/services/database/models"
	"context"

	"github.com/arthurkushman/buildsqlx"
	_ "github.com/lib/pq"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PostgresDB *buildsqlx.DB
var MongoDB *mongo.Database

// ...

func Setup() error {
	// Read conf
	postgresURI := conf.Env.GetString("POSTGRES_URI")
	mongoURI := conf.Env.GetString("MONGO_URI")
	// Connect DBs
	postgresClient := buildsqlx.NewConnection("postgres", postgresURI)
	mongodbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	// Build DBs
	PostgresDB = buildsqlx.NewDb(postgresClient)
	MongoDB = mongodbClient.Database(conf.Env.GetString("MONGO_DB"))
	// Check DB connections
	err = PostgresDB.Sql().Ping()
	if err != nil {
		return err
	}
	err = mongodbClient.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	// Build Models
	models.Build(PostgresDB, MongoDB)

	return nil
}
