package storage

import (
	"API-REST/services/storage/gcs"
	"API-REST/services/storage/local"
)

var Local *local.Storage
var GCS *gcs.Storage

func Setup() error {
	err := local.Setup(Local)
	if err != nil {
		return err
	}
	err = gcs.Setup(GCS)
	if err != nil {
		return err
	}

	return nil
}
