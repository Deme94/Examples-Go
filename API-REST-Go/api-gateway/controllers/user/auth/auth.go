package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"

	"API-REST/api-gateway/utilities/templates"
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
func (Controller) GenerateJwtToken(subject string, secret string) ([]byte, error) {
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

func (Controller) ValidateJwtToken(token []byte, secret string) (uuid.UUID, error) {
	claims, err := jwt.HMACCheck(token, []byte(secret))
	if err != nil {
		return uuid.Nil, errors.New("unauthorized - invalid token")
	}

	if !claims.Valid(time.Now()) {
		return uuid.Nil, errors.New("unauthorized - token expired")
	}

	domain := conf.Env.GetString("DOMAIN")
	if !claims.AcceptAudience(domain) {
		return uuid.Nil, errors.New("unauthorized - invalid audience")
	}

	if claims.Issuer != domain {
		return uuid.Nil, errors.New("unauthorized - invalid issuer")
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("unauthorized - invalid claimer")
	}

	return id, nil
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

func (Controller) GenerateConfirmationEmail(token []byte) string {
	html := strings.ReplaceAll(templates.CONFIRM_EMAIL, "FRONT_DOMAIN", conf.Env.GetString("FRONT_DOMAIN"))
	html = strings.ReplaceAll(html, "FRONT_PRODUCT_NAME", conf.Env.GetString("FRONT_PRODUCT_NAME"))
	html = strings.ReplaceAll(html, "FRONT_LOGO_URL", conf.Env.GetString("FRONT_LOGO_URL"))
	html = strings.ReplaceAll(html, "CONFIRM_EMAIL_ROUTE", conf.Env.GetString("CONFIRM_EMAIL_ROUTE"))
	html = strings.ReplaceAll(html, "COMPANY_NAME", conf.Env.GetString("COMPANY_NAME"))
	html = strings.ReplaceAll(html, "COMPANY_OWNER", conf.Env.GetString("COMPANY_OWNER"))
	html = strings.ReplaceAll(html, "CONFIRM_EMAIL_TOKEN", string(token))

	return html
}

func (c *Controller) GetRoles(userID uuid.UUID) ([]string, error) {
	user, err := c.Model.Get(userID)

	var roleNames []string
	for _, userRole := range user.Roles {
		roleNames = append(roleNames, userRole.Name)
	}
	return roleNames, err
}

func (c *Controller) HasPermission(userID uuid.UUID, resource string, operation string) (bool, error) {
	return c.Model.HasPermission(userID, &permission.Permission{Resource: resource, Operation: operation})
}

func (c *Controller) HasVerifiedEmail(userID uuid.UUID) (bool, error) {
	return c.Model.HasVerifiedEmail(userID)
}

// ...
