package payloads

import "github.com/google/uuid"

type GetAllResponse struct {
	Users []*GetResponse `json:"users,omitempty"`
}
type GetResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}
