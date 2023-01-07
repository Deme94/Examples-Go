package permission

type Permission struct {
	ID        int    `json:"id"`
	Resource  string `json:"resource"`
	Operation string `json:"operation"`
	// ...
}
