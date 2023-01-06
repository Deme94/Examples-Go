package postgres

import (
	"API-REST/services/conf"
	"API-REST/services/database/postgres/models"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"

	"github.com/arthurkushman/buildsqlx"
)

func Setup() error {
	// Read conf
	postgresURI := conf.Env.GetString("POSTGRES_URI")
	postgresDB := conf.Env.GetString("POSTGRES_DB")
	// Connect DB
	postgresClient := buildsqlx.NewConnection("postgres", postgresURI+postgresDB)
	// Build DB
	db := buildsqlx.NewDb(postgresClient)
	// Check DB connection
	err := db.Sql().Ping()
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
