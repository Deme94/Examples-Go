package role

import (
	"errors"

	"github.com/arthurkushman/buildsqlx"
)

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
