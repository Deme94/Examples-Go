package feature

import (
	"API-REST/services/database/postgres/models/feature"

	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Validate *validator.Validate
	Model    *feature.Model
}
