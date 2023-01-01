package payloads

import "API-REST/services/database/models/user"

type GetAllResponse struct {
	Users []*user.User `json:"users"`
}

type LoginResponse struct {
	Id    int    `json:"user_id"`
	Token string `json:"token"`
}

type OkResponse struct {
	OK bool `json:"ok"`
}
