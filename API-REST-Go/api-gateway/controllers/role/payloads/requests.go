package payloads

type InsertRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdatePermissionsRequest struct {
	PermissionIDs []int `json:"permission_ids" binding:"required"`
}
