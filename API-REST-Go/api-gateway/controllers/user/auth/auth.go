package auth

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"

	"API-REST/services/conf"
	"API-REST/services/database/postgres/models/permission"
	"API-REST/services/database/postgres/models/user"

	pswd "github.com/sethvargo/go-password/password"
)

type Controller struct {
	Validate *validator.Validate
	Model    *user.Model
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (Controller) generateJwtToken(subject string, secret string) ([]byte, error) {
	domain := conf.Env.GetString("DOMAIN")
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(subject)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = domain
	claims.Audiences = []string{domain}

	token, err := claims.HMACSign(jwt.HS256, []byte(secret))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (Controller) generateRandomPassword() string {
	password, _ := pswd.Generate(8, 4, 4, true, true)
	return password
}

func (Controller) compareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (c *Controller) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (c *Controller) GetRoles(userID int) ([]string, error) {
	user, err := c.Model.Get(userID)

	var roleNames []string
	for _, userRole := range user.Roles {
		roleNames = append(roleNames, userRole.Name)
	}
	return roleNames, err
}

func (c *Controller) HasPermission(userID int, resource string, operation string) (bool, error) {
	return c.Model.HasPermission(userID, &permission.Permission{Resource: resource, Operation: operation})
}

// ...
