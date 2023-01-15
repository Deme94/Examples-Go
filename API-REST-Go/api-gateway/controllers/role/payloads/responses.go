package payloads

type GetAllResponse struct {
	Roles []*GetResponse `json:"roles"`
}
type GetResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OkResponse struct {
	OK bool `json:"ok"`
}
