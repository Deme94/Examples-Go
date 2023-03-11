package permission

import (
	"API-REST/services/database/postgres/models/permission"

	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Validate *validator.Validate
	Model    *permission.Model
}
