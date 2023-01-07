package role

import (
	"API-REST/services/database/postgres/models/permission"
	"os/user"
)

type Role struct {
	ID          int                      `json:"id"`
	Name        string                   `json:"name"`
	Permissions []*permission.Permission `json:"permissions"`
	Users       []*user.User             `json:"users"`
	// ...
}
