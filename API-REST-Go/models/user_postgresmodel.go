package models

import (
	"errors"
	"strings"
	"time"

	"github.com/arthurkushman/buildsqlx"
)

type User struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string
	CreatedAt          time.Time `json:"created_at"`
	LastLogin          time.Time `json:"last_login"`
	LastPasswordChange time.Time `json:"last_password_change"`
}

// DB MODEL ****************************************************************
type UserModel struct {
	Db *buildsqlx.DB
}

func NewUserModel(db *buildsqlx.DB) *UserModel {
	return &UserModel{db}
}

// DB QUERIES -----------------------------------------------------------
func (t *UserModel) GetAll() ([]*User, error) {
	res, err := t.Db.Table("users").Select("id", "name", "email", "created_at", "last_login", "last_password_change").Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("model users is empty")
	}

	var usrs []*User
	for _, r := range res {
		// Atributes can be nil
		lastLogin := time.Time{}
		ll := r["last_login"]
		if ll != nil {
			lastLogin = ll.(time.Time)
		}

		u := User{
			ID:                 int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
			Name:               r["name"].(string),
			Email:              r["email"].(string),
			CreatedAt:          r["created_at"].(time.Time),
			LastLogin:          lastLogin,
			LastPasswordChange: r["last_password_change"].(time.Time),
		}

		usrs = append(usrs, &u)
	}

	return usrs, nil
}
func (t *UserModel) Get(id int) (*User, error) {
	res, err := t.Db.Table("users").Select("name", "email", "created_at", "last_login", "last_password_change").Where("id", "=", id).Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}

	r := res[0]

	// Attributes can be nil
	lastLogin := time.Time{}
	ll := r["last_login"]
	if ll != nil {
		lastLogin = ll.(time.Time)
	}

	u := User{
		ID:                 id, // la DB devuelve interface{} y se hace cast a int
		Name:               r["name"].(string),
		Email:              r["email"].(string),
		CreatedAt:          r["created_at"].(time.Time),
		LastLogin:          lastLogin,
		LastPasswordChange: r["last_password_change"].(time.Time),
	}

	return &u, nil
}
func (t *UserModel) Insert(u *User) error {
	email := strings.ToLower(u.Email)

	err := t.Db.Table("users").Insert(map[string]interface{}{"name": u.Name, "email": email, "password": u.Password, "created_at": "NOW()", "last_password_change": "NOW()"})
	if err != nil {
		return err
	}
	return nil
}
func (t *UserModel) Update(u *User) error {
	_, err := t.Db.Table("users").Where("id", "=", u.ID).Update(map[string]interface{}{"name": u.Name, "email": u.Email, "password": u.Password, "last_password_change": "NOW()"})
	return err
}
func (t *UserModel) Delete(id int) error {
	_, err := t.Db.Table("users").Where("id", "=", id).Delete()
	return err
}

func (t *UserModel) GetByEmailWithPassword(email string) (*User, error) {
	email = strings.ToLower(email)

	res, err := t.Db.Table("users").Select("id", "name", "password", "created_at", "last_login", "last_password_change").Where("email", "=", email).Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}
	r := res[0]

	u := User{
		ID:                 r["id"].(int), // la DB devuelve interface{} y se hace cast a int
		Name:               r["name"].(string),
		Email:              email,
		Password:           r["password"].(string),
		CreatedAt:          r["created_at"].(time.Time),
		LastLogin:          r["last_login"].(time.Time),
		LastPasswordChange: r["last_password_change"].(time.Time),
	}

	return &u, nil
}
func (t *UserModel) UpdatePassword(id int, password string) error {
	return nil
}

// ...
