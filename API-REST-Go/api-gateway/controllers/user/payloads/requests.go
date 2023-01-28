package payloads

import "mime/multipart"

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type InsertRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
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
	RoleIDs []int `json:"role_ids" binding:"required"`
}

type UpdatePhotoRequest struct {
	PhotoBase64 string `json:"photo_base64" binding:"required"`
}

type UpdateCVRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
