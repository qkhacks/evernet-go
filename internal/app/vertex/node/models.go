package node

type Node struct {
	Identifier        string `json:"identifier" db:"identifier"`
	DisplayName       string `json:"display_name" db:"display_name"`
	SigningPrivateKey string `json:"-" db:"signing_private_key"`
	SigningPublicKey  string `json:"signing_public_key" db:"signing_public_key"`
	Creator           string `json:"creator" db:"creator"`
	CreatedAt         int64  `json:"created_at" db:"created_at"`
	UpdatedAt         int64  `json:"updated_at" db:"updated_at"`
}
