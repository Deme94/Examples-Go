package role

import (
	"API-REST/services/database/postgres/models/role"

	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Validate *validator.Validate
	Model    *role.Model
}
