package admin

import (
	"fmt"
	"github.com/evernetproto/evernet/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Authenticator struct {
	jwtSigningKey []byte
	vertex        string
}

func NewAuthenticator(jwtSigningKey string, vertex string) *Authenticator {
	return &Authenticator{
		jwtSigningKey: []byte(jwtSigningKey),
		vertex:        vertex,
	}
}

const (
	BearerToken    = "Bearer"
	TokenTypeAdmin = "admin"
)

func (a *Authenticator) GenerateToken(identifier string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  identifier,
		"iss":  a.vertex,
		"aud":  a.vertex,
		"type": TokenTypeAdmin,
		"iat":  int(time.Now().Unix()),
	})

	tokenString, err := token.SignedString(a.jwtSigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Authenticator) ValidateContext(c *gin.Context) (*Admin, error) {
	tokenType, token, err := api.ExtractToken(c)

	if err != nil {
		return nil, err
	}

	switch tokenType {
	case BearerToken:
		return a.validateBearerToken(token)
	default:
		return nil, fmt.Errorf("invalid token type")
	}
}

func (a *Authenticator) validateBearerToken(tokenString string) (*Admin, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSigningKey, nil
	}, jwt.WithAudience(a.vertex), jwt.WithIssuer(a.vertex))

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identifier, ok := claims["sub"]

		if !ok {
			return nil, fmt.Errorf("invalid access token")
		}

		identifierString, ok := identifier.(string)

		if !ok {
			return nil, fmt.Errorf("invalid access token")
		}

		tokenType, ok := claims["type"]

		if !ok {
			return nil, fmt.Errorf("invalid access token")
		}

		tokenTypeString, ok := tokenType.(string)

		if !ok {
			return nil, fmt.Errorf("invalid access token")
		}

		if tokenTypeString != TokenTypeAdmin {
			return nil, fmt.Errorf("invalid access token")
		}

		return &Admin{
			Identifier: identifierString,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid access token")
	}
}
