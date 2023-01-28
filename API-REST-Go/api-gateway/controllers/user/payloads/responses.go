package payloads

type GetAllResponse struct {
	Users []*GetResponse `json:"users,omitempty"`
}
type GetResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginResponse struct {
	ID    int    `json:"user_id"`
	Token string `json:"token"`
}
