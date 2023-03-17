package storage

import (
	"API-REST/services/storage/gcs"
	"API-REST/services/storage/local"
)

var Local *local.Storage
var GCS *gcs.Storage

func SetupLocal() error {
	return local.Setup(Local)
}
func SetupGCS() error {
	return gcs.Setup(GCS)
}
