package postgres

import (
	"API-REST/services/conf"
	"API-REST/services/database/postgres/models"
	"database/sql"
	"errors"
	"fmt"
	"os"
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

	// Create DB and tables from postgres.sql file
	err := createDBTables()
	if err != nil {
		return err
	}

	// Insert default data into tables from conf.yml
	err = insertConfData()
	if err != nil {
		return err
	}

	return nil
}

func createDBTables() error {
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

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sqlQuery := string(bytes)

	// Open default DB
	db, err := sql.Open("postgres", postgresURI+postgresDefaultDB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create new DB
	_, err = db.Exec(`CREATE DATABASE ` + postgresDB)
	if err != nil {
		if !strings.Contains(fmt.Sprint(err), "already exists") {
			return err
		}
	}

	// Open new DB
	db, err = sql.Open("postgres", postgresURI+postgresDB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create Tables
	_, err = db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

func insertConfData() error {
	// Read conf
	postgresURI := conf.Env.GetString("POSTGRES_URI")
	postgresDB := conf.Env.GetString("POSTGRES_DB")

	permissions := conf.Conf.GetStringMapStringSlice("permissions")
	roles := conf.Conf.GetStringMap("roles")

	// Open DB
	db, err := sql.Open("postgres", postgresURI+postgresDB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert permissions
	insertPermissions := "INSERT INTO permissions (resource, operation) VALUES "
	for resource, operations := range permissions {
		for _, operation := range operations {
			insertPermissions += "('" + resource + "', '" + operation + "'),"
		}
	}
	insertPermissions = strings.TrimSuffix(insertPermissions, ",") + ";"

	_, err = db.Exec(insertPermissions)
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
