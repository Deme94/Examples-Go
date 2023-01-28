package payloads

type GetAllResponse struct {
	Permissions []*GetResponse `json:"permissions,omitempty"`
}
type GetResponse struct {
	ID        int    `json:"id"`
	Resource  string `json:"resource"`
	Operation string `json:"operation"`
}
