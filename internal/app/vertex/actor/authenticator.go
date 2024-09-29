package actor

import (
	"context"
	"fmt"
	"github.com/evernetproto/evernet/internal/app/vertex/node"
	"github.com/evernetproto/evernet/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type AuthenticatedActor struct {
	Identifier           string
	Address              string
	SourceNodeIdentifier string
	SourceVertex         string
	SourceNodeAddress    string
	TargetNodeIdentifier string
	TargetVertex         string
	TargetNodeAddress    string
	IsLocal              bool
}

type Authenticator struct {
	vertex            string
	nodeManager       *node.Manager
	remoteNodeManager *node.RemoteManager
}

func NewAuthenticator(vertex string, nodeManager *node.Manager, remoteNodeManager *node.RemoteManager) *Authenticator {
	return &Authenticator{vertex: vertex, nodeManager: nodeManager, remoteNodeManager: remoteNodeManager}
}

const (
	TokenTypeActor = "actor"
	BearerToken    = "Bearer"
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

func (a *Authenticator) ValidateContext(ctx context.Context, c *gin.Context) (*AuthenticatedActor, error) {
	tokenType, token, err := api.ExtractToken(c)

	if err != nil {
		return nil, err
	}

	switch tokenType {
	case BearerToken:
		return a.validateBearerToken(ctx, token)
	default:
		return nil, fmt.Errorf("invalid token type")
	}
}

func (a *Authenticator) validateBearerToken(ctx context.Context, tokenString string) (*AuthenticatedActor, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		issuer, ok := claims["iss"]

		if !ok {
			return nil, fmt.Errorf("invalid token issuer")
		}

		issuerString, ok := issuer.(string)

		if !ok {
			return nil, fmt.Errorf("invalid token issuer")
		}

		issuerComponents := strings.Split(issuerString, "/")

		if len(issuerComponents) != 2 {
			return nil, fmt.Errorf("invalid token issuer")
		}

		sourceVertex := issuerComponents[0]
		sourceNodeIdentifier := issuerComponents[1]

		if sourceVertex == a.vertex {
			sourceNode, err := a.nodeManager.Get(ctx, sourceNodeIdentifier)

			if err != nil {
				return nil, err
			}

			return sourceNode.GetSigningPublicKey()
		} else {
			sourceNode, err := a.remoteNodeManager.Get(ctx, sourceVertex, sourceNodeIdentifier)

			if err != nil {
				return nil, err
			}

			return sourceNode.GetSigningPublicKey()
		}
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

		if tokenTypeString != TokenTypeActor {
			return nil, fmt.Errorf("invalid access token")
		}

		issuer, ok := claims["iss"]

		if !ok {
			return nil, fmt.Errorf("invalid token issuer")
		}

		sourceNodeAddress, ok := issuer.(string)

		if !ok {
			return nil, fmt.Errorf("invalid token issuer")
		}

		issuerComponents := strings.Split(sourceNodeAddress, "/")

		if len(issuerComponents) != 2 {
			return nil, fmt.Errorf("invalid token issuer")
		}

		sourceVertex := issuerComponents[0]
		sourceNodeIdentifier := issuerComponents[1]

		aud, ok := claims["aud"]

		if !ok {
			return nil, fmt.Errorf("invalid token audience")
		}

		targetNodeAddress, ok := aud.(string)

		if !ok {
			return nil, fmt.Errorf("invalid token audience")
		}

		audienceComponents := strings.Split(targetNodeAddress, "/")

		if len(audienceComponents) != 2 {
			return nil, fmt.Errorf("invalid token audience")
		}

		targetVertex := audienceComponents[0]
		targetNodeIdentifier := audienceComponents[1]

		return &AuthenticatedActor{
			Identifier:           identifierString,
			Address:              fmt.Sprintf("%s/%s/%s", sourceVertex, sourceNodeIdentifier, identifierString),
			SourceNodeIdentifier: sourceNodeIdentifier,
			SourceVertex:         sourceVertex,
			SourceNodeAddress:    sourceNodeAddress,
			TargetNodeIdentifier: targetNodeIdentifier,
			TargetVertex:         targetVertex,
			TargetNodeAddress:    targetNodeAddress,
			IsLocal:              sourceNodeIdentifier == targetNodeIdentifier && sourceVertex == targetVertex,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid access token")
	}
}
