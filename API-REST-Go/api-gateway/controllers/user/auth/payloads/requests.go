package payloads

import "mime/multipart"

type LoginRequest struct {
	Username string `json:"username" validate:"required_without=Email"`
	Email    string `json:"email" validate:"required_without=Username"`
	Password string `json:"password" validate:"required"`
}

type UpdateRequest struct {
	Nick      string `json:"nick,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type UpdateRolesRequest struct {
	RoleIDs []int `json:"role_ids" validate:"required"`
}

type UpdatePhotoRequest struct {
	PhotoBase64 string `json:"photo_base64" validate:"required,base64"`
}

type UpdateCVRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
}
