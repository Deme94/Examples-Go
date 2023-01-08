package payloads

type GetAllResponse struct {
	Permissions []*GetResponse `json:"permissions"`
}
type GetResponse struct {
	ID        int    `json:"id"`
	Resource  string `json:"resource"`
	Operation string `json:"operation"`
}

type OkResponse struct {
	OK bool `json:"ok"`
}
