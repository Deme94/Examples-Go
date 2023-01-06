package models

import (
	"API-REST/services/database/postgres/models/user"

	"github.com/arthurkushman/buildsqlx"
)

var (
	User *user.Model
	// ...
)

func Build(postgresDB *buildsqlx.DB) error {
	var err error

	User, err = user.New(postgresDB)
	if err != nil {
		return err
	}
	// ...

	return nil
}
