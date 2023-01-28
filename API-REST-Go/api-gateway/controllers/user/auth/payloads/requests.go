package payloads

import "mime/multipart"

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type UpdateRequest struct {
	Nick      string `json:"nick"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateRolesRequest struct {
	RoleIDs []int `json:"role_ids" binding:"required"`
}

type UpdatePhotoRequest struct {
	PhotoBase64 string `json:"photo_base64" binding:"required"`
}

type UpdateCVRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
