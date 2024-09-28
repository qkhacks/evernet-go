package actor

type SignUpRequest struct {
	Identifier  string `json:"identifier" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Type        string `json:"type" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type TokenRequest struct {
	Identifier        string `json:"identifier" binding:"required"`
	Password          string `json:"password" binding:"required"`
	TargetNodeAddress string `json:"target_node_address"`
}

type PasswordChangeRequest struct {
	Password string `json:"password" binding:"required"`
}
