package user

import (
	"API-REST/api-gateway/controllers/user/auth"
	"API-REST/services/database/postgres/models/user"

	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	Model *user.Model
	Auth  *auth.Controller
}

func (c *Controller) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
