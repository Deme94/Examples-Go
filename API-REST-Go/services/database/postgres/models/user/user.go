package user

import (
	"API-REST/services/database/postgres/models/role"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`

	Nick      string `json:"nick"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`

	LastLogin          time.Time `json:"last_login"`
	LastPasswordChange time.Time `json:"last_password_change"`
	VerifiedMail       bool      `json:"verified_mail"`
	VerifiedPhone      bool      `json:"verified_phone"`
	BanDate            time.Time `json:"ban_date"`
	BanExpire          time.Time `json:"ban_expire"`

	Roles []*role.Role `json:"roles"`

	PhotoName string `json:"-"`
	CVName    string `json:"-"`
	// ...
}
