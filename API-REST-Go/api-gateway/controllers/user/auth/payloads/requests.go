package payloads

import "mime/multipart"

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type UpdateRequest struct {
	Nick      string `json:"nick"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
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
	PhotoBase64 string `json:"photo_base64" validate:"required"`
}

type UpdateCVRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
}
