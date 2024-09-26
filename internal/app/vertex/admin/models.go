package admin

type Admin struct {
	Identifier string `json:"identifier" db:"identifier"`
	Password   string `json:"-" db:"password"`
	Creator    string `json:"creator" db:"creator"`
	CreatedAt  int64  `json:"created_at" db:"created_at"`
	UpdatedAt  int64  `json:"updated_at" db:"updated_at"`
}
