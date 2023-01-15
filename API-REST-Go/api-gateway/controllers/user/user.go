package user

import (
	"API-REST/api-gateway/controllers/user/auth"
	"API-REST/services/database/postgres/models/user"
)

type Controller struct {
	Model *user.Model
	Auth  *auth.Controller
}
