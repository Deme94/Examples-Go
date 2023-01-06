package role

import (
	"API-REST/services/database/postgres/models/permission"
	"errors"
	"os/user"

	"github.com/arthurkushman/buildsqlx"
)

// MAIN STRUCT
type Role struct {
	ID          int                      `json:"id"`
	Name        string                   `json:"name"`
	Permissions []*permission.Permission `json:"permissions"`
	Users       []*user.User             `json:"users"`
	// ...
}

// DB MODEL ****************************************************************
type Model struct {
	Db *buildsqlx.DB
}

func New(db *buildsqlx.DB) (*Model, error) {
	exists, err := db.HasTable("public", "roles")
	if err != nil {
		return nil, err
	}
	if !exists {
		err = errors.New("table roles doesn't exist in db")
		return nil, err
	}
	return &Model{Db: db}, nil
}
