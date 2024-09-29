package messaging

type Inbox struct {
	Identifier     string `json:"identifier" db:"identifier"`
	DisplayName    string `json:"display_name" db:"display_name"`
	NodeIdentifier string `json:"node_identifier" db:"node_identifier"`
	ActorAddress   string `json:"actor_address" db:"actor_address"`
	CreatedAt      int64  `json:"created_at" db:"created_at"`
	UpdatedAt      int64  `json:"updated_at" db:"updated_at"`
}

type Outbox struct {
	Identifier     string `json:"identifier" db:"identifier"`
	DisplayName    string `json:"display_name" db:"display_name"`
	NodeIdentifier string `json:"node_identifier" db:"node_identifier"`
	ActorAddress   string `json:"actor_address" db:"actor_address"`
	CreatedAt      int64  `json:"created_at" db:"created_at"`
	UpdatedAt      int64  `json:"updated_at" db:"updated_at"`
}
