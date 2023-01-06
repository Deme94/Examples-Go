package payloads

type GetAllResponse struct {
	Users []*GetAllUserResponse `json:"users"`
}
type GetAllUserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginResponse struct {
	Id    int    `json:"user_id"`
	Token string `json:"token"`
}

type OkResponse struct {
	OK bool `json:"ok"`
}
