package payloads

type InsertRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdatePermissionsRequest struct {
	PermissionIDs []int `json:"permission_ids" validate:"required"`
}
