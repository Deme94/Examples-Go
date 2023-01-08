package payloads

type InsertRequest struct {
	Resource  string `json:"resource" binding:"required"`
	Operation string `json:"operation" binding:"required"`
}

type UpdateRequest struct {
	Resource  string `json:"resource" binding:"required"`
	Operation string `json:"operation" binding:"required"`
}
