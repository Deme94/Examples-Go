package payloads

import (
	"time"
)

type QueryParams struct {
	Any      *string `query:"any"`
	Page     *int    `query:"page"`
	PageSize *int    `query:"pageSize"`
	Deleted  *bool   `query:"deleted"`
	Banned   *bool   `query:"banned"`
	Year     *int    `query:"year"`
}

type InsertRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email,min=6,max=32"`
	Password  string `json:"password" validate:"required"`
	Nick      string `json:"nick,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

type UpdateRequest struct {
	Nick      string `json:"nick,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
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
