package models

import (
	"API-REST/services/database/postgres/models/feature"
	"API-REST/services/database/postgres/models/permission"
	"API-REST/services/database/postgres/models/role"
	"API-REST/services/database/postgres/models/user"

	"github.com/arthurkushman/buildsqlx"
)

var (
	User       *user.Model
	Role       *role.Model
	Permission *permission.Model
	Feature    *feature.Model
	// ...
)

func Build(postgresDB *buildsqlx.DB) error {
	var err error

	User, err = user.New(postgresDB)
	if err != nil {
		return err
	}
	Role, err = role.New(postgresDB)
	if err != nil {
		return err
	}
	Permission, err = permission.New(postgresDB)
	if err != nil {
		return err
	}
	Feature, err = feature.New(postgresDB)
	if err != nil {
		return err
	}
	// ...

	return nil
}
