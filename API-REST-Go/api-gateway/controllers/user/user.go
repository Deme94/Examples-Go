package user

import (
	"API-REST/api-gateway/controllers/user/auth"
	"API-REST/services/database/postgres/models/user"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	Validate *validator.Validate
	Model    *user.Model
	Auth     *auth.Controller
}

func (c *Controller) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
