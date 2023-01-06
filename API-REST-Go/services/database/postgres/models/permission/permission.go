package permission

import (
	"errors"

	"github.com/arthurkushman/buildsqlx"
)

// MAIN STRUCT
type Permission struct {
	ID        int    `json:"id"`
	Resource  string `json:"resource"`
	Operation string `json:"operation"`
	// ...
}

// DB MODEL ****************************************************************
type Model struct {
	Db *buildsqlx.DB
}

func New(db *buildsqlx.DB) (*Model, error) {
	exists, err := db.HasTable("public", "permissions")
	if err != nil {
		return nil, err
	}
	if !exists {
		err = errors.New("table permissions doesn't exist in db")
		return nil, err
	}
	return &Model{Db: db}, nil
}
