package postgres

import (
	"API-REST/services/conf"
	"API-REST/services/database/postgres/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/arthurkushman/buildsqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func Setup() error {
	// Read conf
	username := conf.Env.GetString("POSTGRES_DB_USERNAME")
	password := conf.Env.GetString("POSTGRES_DB_PASSWORD")
	host := conf.Env.GetString("POSTGRES_DB_HOST")
	port := conf.Env.GetString("POSTGRES_DB_PORT")
	postgresDB := conf.Env.GetString("POSTGRES_DB")

	postgresURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username,
		password,
		host,
		port,
		postgresDB,
	)
	// Connect DB
	postgresClient := buildsqlx.NewConnection("postgres", postgresURI)
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

func SetupDockertest() error {
	// Create docker db
	db, err := createDBDockertest()
	if err != nil {
		return err
	}
	// Create tables
	err = createTables(db)
	if err != nil {
		return err
	}
	// Insert conf values
	err = insertConfData(db)
	if err != nil {
		return err
	}

	// Build sqlx client
	postgresClient := buildsqlx.NewConnectionFromDb(db)
	dbSqlx := buildsqlx.NewDb(postgresClient)
	err = dbSqlx.Sql().Ping()
	if err != nil {
		return err
	}

	// Build Models
	models.Build(dbSqlx)

	return nil
}

// Create DB and tables
func Init() error {

	// Create db if not exists
	db, err := createDB()
	if err != nil {
		return err
	}
	defer db.Close()
	// Create tables from postgres.sql file
	err = createTables(db)
	if err != nil {
		return err
	}
	// Insert default data into tables from conf.yml
	err = insertConfData(db)
	if err != nil {
		return err
	}

	return nil
}

func createDB() (*sql.DB, error) {
	// Read conf
	username := conf.Env.GetString("POSTGRES_DB_USERNAME")
	password := conf.Env.GetString("POSTGRES_DB_PASSWORD")
	host := conf.Env.GetString("POSTGRES_DB_HOST")
	port := conf.Env.GetString("POSTGRES_DB_PORT")
	postgresDefaultDB := conf.Env.GetString("POSTGRES_DEFAULT_DB")
	postgresDB := conf.Env.GetString("POSTGRES_DB")

	// Open default DB
	db, err := sql.Open("postgres", "postgres://"+username+":"+password+"@"+host+":"+port+"/"+postgresDefaultDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Create new DB
	_, err = db.Exec(`CREATE DATABASE ` + postgresDB)
	if err != nil {
		if !strings.Contains(fmt.Sprint(err), "already exists") {
			return nil, err
		}
	}

	// Open new DB
	db, err = sql.Open("postgres", "postgres://"+username+":"+password+"@"+host+":"+port+"/"+postgresDB)
	if err != nil {
		return nil, err
	}

	return db, nil
}
func createDBDockertest() (*sql.DB, error) {
	// Read conf
	username := conf.Env.GetString("POSTGRES_DB_USERNAME")
	password := conf.Env.GetString("POSTGRES_DB_PASSWORD")
	port := conf.Env.GetString("POSTGRES_DB_PORT")
	postgresDB := conf.Env.GetString("POSTGRES_DB")

	// DOCKERTEST
	var db *sql.DB
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + username,
			"POSTGRES_DB=" + postgresDB,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %s", err)
	}
	resource.Expire(30)
	pool.MaxWait = 30 * time.Second

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		uri := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			username, password, resource.GetHostPort(port+"/tcp"), postgresDB)
		db, err = sql.Open("postgres", uri)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to database: %s", err)
	}

	// Destroy db after 30 seconds
	go func() {
		time.Sleep(time.Second * 30)
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}()
	return db, nil
}
func createTables(db *sql.DB) error {
	// Read sql file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("no caller information")
	}
	dir := path.Dir(filename)
	path := dir + "/" + "postgres.sql"

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sqlQuery := string(bytes)

	// Create Tables
	_, err = db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}
func insertConfData(db *sql.DB) error {
	// Read conf
	permissions := conf.Conf.GetStringMapStringSlice("permissions")
	roles := conf.Conf.GetStringMap("roles")

	// Insert permissions
	insertPermissions := "INSERT INTO permissions (resource, operation) VALUES "
	for resource, operations := range permissions {
		for _, operation := range operations {
			insertPermissions += "('" + resource + "', '" + operation + "'),"
		}
	}
	insertPermissions = strings.TrimSuffix(insertPermissions, ",") + ";"

	_, err := db.Exec(insertPermissions)
	if err != nil {
		return err
	}

	// Insert roles and assign permissions
	insertRoles := "INSERT INTO roles (name) VALUES "
	insertRolesPermissions := ""
	for role, permissions := range roles {
		insertRoles += "('" + role + "'),"
		for resource, operations := range permissions.(map[string]interface{}) {
			for _, operation := range operations.([]interface{}) {
				insertRolesPermissions += "INSERT INTO roles_permissions (role_id, permission_id) " +
					"SELECT roles.id, permissions.id FROM roles " +
					"LEFT JOIN permissions " +
					"ON permissions.resource = '" + resource + "' AND operation = '" + operation.(string) + "' " +
					"WHERE roles.name = '" + role + "';"
			}
		}
	}
	insertRoles = strings.TrimSuffix(insertRoles, ",") + ";"

	_, err = db.Exec(insertRoles)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertRolesPermissions)
	if err != nil {
		return err
	}
	return nil
}
