package actor

type SignUpRequest struct {
	Identifier  string `json:"identifier" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Type        string `json:"type" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}
