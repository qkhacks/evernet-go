package node

type CreationRequest struct {
	Identifier  string `json:"identifier" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type UpdateRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
}
