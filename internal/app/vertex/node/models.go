package node

import (
	"crypto/ed25519"
	"fmt"
	"github.com/evernetproto/evernet/internal/pkg/keys"
)

type Node struct {
	Identifier        string `json:"identifier" db:"identifier"`
	DisplayName       string `json:"display_name" db:"display_name"`
	SigningPrivateKey string `json:"-" db:"signing_private_key"`
	SigningPublicKey  string `json:"signing_public_key" db:"signing_public_key"`
	Creator           string `json:"creator" db:"creator"`
	CreatedAt         int64  `json:"created_at" db:"created_at"`
	UpdatedAt         int64  `json:"updated_at" db:"updated_at"`
}

func (n *Node) GetAddress(vertex string) string {
	return fmt.Sprintf("%s/%s", vertex, n.Identifier)
}

func (n *Node) GetSigningPrivateKey() (ed25519.PrivateKey, error) {
	return keys.ConvertED25519PrivateKeyFromString(n.SigningPrivateKey)
}

func (n *Node) GetSigningPublicKey() (ed25519.PublicKey, error) {
	return keys.ConvertED25519PublicKeyFromString(n.SigningPublicKey)
}
