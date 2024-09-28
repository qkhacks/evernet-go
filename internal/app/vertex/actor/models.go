package actor

type Actor struct {
	Identifier     string `json:"id" db:"id"`
	Password       string `json:"-" db:"password"`
	Type           string `json:"type" db:"type"`
	DisplayName    string `json:"display_name" db:"display_name"`
	NodeIdentifier string `json:"node_identifier" db:"node_identifier"`
	Creator        string `json:"creator" db:"creator"`
	CreatedAt      int64  `json:"created_at" db:"created_at"`
	UpdatedAt      int64  `json:"updated_at" db:"updated_at"`
}
