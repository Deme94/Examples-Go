package payloads

type InsertRequest struct {
	Resource  string `json:"resource" validate:"required"`
	Operation string `json:"operation" validate:"required"`
}

type UpdateRequest struct {
	Resource  string `json:"resource" validate:"required"`
	Operation string `json:"operation" validate:"required"`
}
