package user

import (
	"API-REST/services/database/postgres/models/role"
	"time"
)

type User struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`

	Nick      string `json:"nick"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`

	LastLogin          *time.Time `json:"last_login,omitempty"`
	LastPasswordChange *time.Time `json:"last_password_change,omitempty"`
	VerifiedMail       bool       `json:"verified_mail,omitempty"`
	VerifiedPhone      bool       `json:"verified_phone,omitempty"`
	BanDate            *time.Time `json:"ban_date,omitempty"`
	BanExpire          *time.Time `json:"ban_expire,omitempty"`

	Roles []*role.Role `json:"roles,omitempty"`

	PhotoName string `json:"-"`
	CVName    string `json:"-"`
	// ...
}
