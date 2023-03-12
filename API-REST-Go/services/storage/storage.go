package storage

import (
	"API-REST/services/storage/local"
)

var Local *local.Storage

func Setup() error {
	err := local.Setup(Local)
	if err != nil {
		return err
	}

	return nil
}
