package user

import (
	"fmt"
	"time"

	"github.com/pascaldekloe/jwt"

	"API-REST/services/conf"
	"API-REST/services/database/models/user"
)

// CONTROLLER ***************************************************************
type Controller struct {
	Model *user.Model
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (Controller) generateJwtToken(subject string, secret string) ([]byte, error) {
	domain := conf.Domain
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

func (c *Controller) CheckRole(id int) (string, error) {
	return c.Model.GetRole(id)
}

// ...
