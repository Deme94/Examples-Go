package models

import (
	"errors"
	"strings"
	"time"

	"github.com/arthurkushman/buildsqlx"
)

type User struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Password           string    `json:"-"`
	PhotoName          string    `json:"-"`
	CVName             string    `json:"-"`
	CreatedAt          time.Time `json:"created_at"`
	LastLogin          time.Time `json:"last_login"`
	LastPasswordChange time.Time `json:"last_password_change"`
	Role               string    `json:"role"`
}

// DB MODEL ****************************************************************
type UserModel struct {
	Db *buildsqlx.DB
}

func NewUserModel(db *buildsqlx.DB) *UserModel {
	return &UserModel{db}
}

// DB QUERIES -----------------------------------------------------------
func (m *UserModel) GetAll() ([]*User, error) {
	res, err := m.Db.Table("users").Select("users.id", "name", "email", "created_at", "last_login", "last_password_change", "role").
		LeftJoin("roles", "users.role_id", "=", "roles.id").
		Get()
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
			Role:               r["role"].(string),
		}

		usrs = append(usrs, &u)
	}

	return usrs, nil
}
func (m *UserModel) Get(id int) (*User, error) {
	res, err := m.Db.Table("users").Select("name", "email", "created_at", "last_login", "last_password_change", "role").
		Where("users.id", "=", id).
		LeftJoin("roles", "users.role_id", "=", "roles.id").
		Get()
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
		Role:               r["role"].(string),
	}

	return &u, nil
}
func (m *UserModel) GetByEmailWithPassword(email string) (*User, error) {
	email = strings.ToLower(email)

	res, err := m.Db.Table("users").Select("id", "name", "password", "created_at", "last_login", "last_password_change").Where("email", "=", email).Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}
	r := res[0]

	// Atributes can be nil
	lastLogin := time.Time{}
	ll := r["last_login"]
	if ll != nil {
		lastLogin = ll.(time.Time)
	}

	u := User{
		ID:                 int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
		Name:               r["name"].(string),
		Email:              email,
		Password:           r["password"].(string),
		CreatedAt:          r["created_at"].(time.Time),
		LastLogin:          lastLogin,
		LastPasswordChange: r["last_password_change"].(time.Time),
	}

	return &u, nil
}
func (m *UserModel) GetByNameWithPassword(name string) (*User, error) {
	res, err := m.Db.Table("users").Select("id", "email", "password", "created_at", "last_login", "last_password_change").Where("name", "=", name).Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}
	r := res[0]

	// Atributes can be nil
	lastLogin := time.Time{}
	ll := r["last_login"]
	if ll != nil {
		lastLogin = ll.(time.Time)
	}

	u := User{
		ID:                 int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
		Name:               name,
		Email:              r["email"].(string),
		Password:           r["password"].(string),
		CreatedAt:          r["created_at"].(time.Time),
		LastLogin:          lastLogin,
		LastPasswordChange: r["last_password_change"].(time.Time),
	}

	return &u, nil
}
func (m *UserModel) GetPhoto(id int) (string, error) {
	res, err := m.Db.Table("users").Select("photo_name").Where("id", "=", id).Get()
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("user not found")
	}

	return res[0]["photo_name"].(string), nil
}
func (m *UserModel) GetCV(id int) (string, error) {
	res, err := m.Db.Table("users").Select("cv_name").Where("id", "=", id).Get()
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("user not found")
	}

	return res[0]["cv_name"].(string), nil
}
func (m *UserModel) GetRole(id int) (string, error) {
	res, err := m.Db.Table("users").Select("role").Where("users.id", "=", id).LeftJoin("roles", "users.role_id", "=", "roles.id").Get() // hay que hacer join
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("user not found")
	}

	return res[0]["role"].(string), nil
}
func (m *UserModel) Insert(u *User) error {
	email := strings.ToLower(u.Email)

	err := m.Db.Table("users").Insert(map[string]interface{}{"name": u.Name, "email": email, "password": u.Password, "created_at": "NOW()", "last_password_change": "NOW()"})
	if err != nil {
		return err
	}
	return nil
}
func (m *UserModel) Update(u *User) error {
	_, err := m.Db.Table("users").Where("id", "=", u.ID).Update(map[string]interface{}{"name": u.Name, "email": u.Email, "password": u.Password, "last_password_change": "NOW()"})
	return err
}
func (m *UserModel) UpdatePassword(id int, password string) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Update(map[string]interface{}{"password": password, "last_password_change": "NOW()"})
	return err
}
func (m *UserModel) UpdatePhoto(id int, photoName string) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Update(map[string]interface{}{"photo_name": photoName})
	return err
}
func (m *UserModel) UpdateCV(id int, cvName string) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Update(map[string]interface{}{"cv_name": cvName})
	return err
}
func (m *UserModel) Delete(id int) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Delete()
	return err
}

// ...
