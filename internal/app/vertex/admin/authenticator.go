package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
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
	tokenType, token, err := a.extractToken(c)

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
	})

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

func (a *Authenticator) extractToken(c *gin.Context) (string, string, error) {
	authorizationHeader := c.GetHeader("Authorization")

	if len(authorizationHeader) == 0 {
		return "", "", fmt.Errorf("authorization header is not set")
	}

	components := strings.Split(authorizationHeader, " ")

	if len(components) != 2 {
		return "", "", fmt.Errorf("invalid authorization header")
	}

	return components[0], components[1], nil
}
