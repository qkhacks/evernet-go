package messaging

type InboxCreationRequest struct {
	Identifier  string `json:"identifier" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type InboxUpdateRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
}
