package payloads

import (
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type InsertRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email,min=6,max=32"`
	Password  string `json:"password" validate:"required"`
	Nick      string `json:"nick"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type UpdateRequest struct {
	Nick      string `json:"nick"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type UpdateRolesRequest struct {
	RoleIDs []int `json:"role_ids" validate:"required"`
}

type BanRequest struct {
	BanExpire time.Time `json:"ban_expire" validate:"required" example:"2006-01-02T00:00:00Z"`
}

type UpdatePhotoRequest struct {
	PhotoBase64 string `json:"photo_base64" validate:"required"`
}
