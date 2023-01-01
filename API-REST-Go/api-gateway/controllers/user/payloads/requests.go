package payloads

import "mime/multipart"

type UpdateRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhotoBase64 string `json:"photo_base64"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePhotoRequest struct {
	PhotoBase64 string `json:"photo_base64"`
}

type UpdateCVRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
