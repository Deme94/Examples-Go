package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/arthurkushman/buildsqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"

	util "API-REST/cmd/api/utilities"
	m "API-REST/models"
)

// CONTROLLER ***************************************************************
type UserController struct {
	model     *m.UserModel
	logger    *log.Logger
	jwtSecret string
	domain    string
}

func NewUserController(db *buildsqlx.DB, logger *log.Logger, secret string, domain string) *UserController {
	c := UserController{}
	c.model = &m.UserModel{Db: db}
	c.logger = logger
	c.jwtSecret = secret
	c.domain = domain

	return &c
}

// METHODS CONTROLLER ---------------------------------------------------------------
func (c *UserController) generateJwtToken(subject string, secret string) ([]byte, error) {

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(subject)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = c.domain
	claims.Audiences = []string{c.domain}

	token, err := claims.HMACSign(jwt.HS256, []byte(secret))
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ...

// PAYLOADS (json input) ----------------------------------------------------------------
type userPayload struct {
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	LastPasswordChange time.Time `json:"last_password_change"`
}
type credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ...

// API HANDLERS ---------------------------------------------------------------
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	usrs, err := c.model.GetAll()
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, usrs, "users")
}
func (c *UserController) Get(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		c.logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(w, err)
		return
	}

	u, err := c.model.Get(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, u, "user")
}
func (c *UserController) Insert(w http.ResponseWriter, r *http.Request) {
	var creds credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	err = c.model.Insert(&m.User{Name: creds.Name, Email: creds.Email, Password: string(hashedPassword)})
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, "user created successfully", "response")

}
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		c.logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(w, err)
		return
	}

	var p userPayload

	uJson := r.PostFormValue("user")
	err = json.Unmarshal([]byte(uJson), &p)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	var u m.User
	u.ID = id
	u.Name = p.Name // cambiar por ruta del archivo creado
	u.Email = p.Email
	u.Password = p.Password
	u.LastPasswordChange = p.LastPasswordChange

	err = c.model.Update(&u)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}
	util.WriteJSON(w, http.StatusOK, ok, "OK")
}
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		c.logger.Print(errors.New("invalid id parameter"))
		util.ErrorJSON(w, err)
		return
	}

	err = c.model.Delete(id)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}
	util.WriteJSON(w, http.StatusOK, ok, "OK")
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var creds credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	u, err := c.model.GetByEmailWithPassword(creds.Email)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	hashedPassword := u.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		util.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// Generate jwt token after successful login
	token, err := c.generateJwtToken(fmt.Sprint(u.ID), c.jwtSecret)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusNotImplemented)
		return
	}

	util.WriteJSON(w, http.StatusOK, string(token), "response")
}

// ...
