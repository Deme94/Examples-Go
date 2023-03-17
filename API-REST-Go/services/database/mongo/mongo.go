package mongo

import (
	"API-REST/services/conf"
	"API-REST/services/database/mongo/models"
	"context"

	_ "github.com/lib/pq"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Setup() error {
	// Read conf
	mongoURI := conf.Env.GetString("MONGO_URI")
	mongoDB := conf.Env.GetString("MONGO_DB")
	// Connect DB
	mongodbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	// Build DB
	db := mongodbClient.Database(mongoDB)
	// Check DB connection
	err = mongodbClient.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	// Build Models
	models.Build(db)

	return nil
}

func Init() error {
	// Insert conf data
	// TODO
	return nil
}
