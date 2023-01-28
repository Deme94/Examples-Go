package role

import (
	"API-REST/services/database/postgres/models/permission"
)

type Role struct {
	ID          int                      `json:"id"`
	Name        string                   `json:"name"`
	Permissions []*permission.Permission `json:"permissions,omitempty"`
	// ...
}
