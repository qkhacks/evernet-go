package actor

import (
	"github.com/evernetproto/evernet/internal/app/vertex/node"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Authenticator struct {
	vertex string
}

func NewAuthenticator(vertex string) *Authenticator {
	return &Authenticator{vertex: vertex}
}

const (
	TokenTypeActor = "actor"
)

func (a *Authenticator) GenerateToken(identifier string, node *node.Node, targetNodeAddress string) (string, error) {
	if targetNodeAddress == "" {
		targetNodeAddress = node.GetAddress(a.vertex)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub":  identifier,
		"iss":  node.GetAddress(a.vertex),
		"aud":  targetNodeAddress,
		"type": TokenTypeActor,
		"iat":  int(time.Now().Unix()),
	})

	signingPrivateKey, err := node.GetSigningPrivateKey()

	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(signingPrivateKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
