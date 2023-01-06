package mongo

import (
	"API-REST/services/conf"
	"API-REST/services/database/mongo/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"

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

// Create DB and tables
func Init() error {
	// Read conf
	postgresURI := conf.Env.GetString("POSTGRES_URI")
	postgresDefaultDB := conf.Env.GetString("POSTGRES_DEFAULT_DB")
	postgresDB := conf.Env.GetString("POSTGRES_DB")

	// Read sql file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("no caller information")
	}
	dir := path.Dir(filename)
	path := dir + "/" + "postgres.sql"

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	sqlQuery := string(bytes)

	// Create DB
	db, err := sql.Open("postgres", postgresURI+postgresDefaultDB)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE DATABASE ` + postgresDB)
	if err != nil {
		if !strings.Contains(fmt.Sprint(err), "already exists") {
			return err
		}
	}

	db, err = sql.Open("postgres", postgresURI+postgresDB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create Tables
	_, err = db.Exec(sqlQuery)

	return err
}
