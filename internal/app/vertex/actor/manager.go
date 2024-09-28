package actor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/evernetproto/evernet/internal/app/vertex/node"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Manager struct {
	dataStore     *DataStore
	nodeManager   *node.Manager
	authenticator *Authenticator
}

func NewManager(dataStore *DataStore, nodeManager *node.Manager, authenticator *Authenticator) *Manager {
	return &Manager{dataStore: dataStore, nodeManager: nodeManager, authenticator: authenticator}
}

func (m *Manager) SignUp(ctx context.Context, nodeIdentifier string, request *SignUpRequest) (*Actor, error) {
	nodeData, err := m.nodeManager.Get(ctx, nodeIdentifier)

	if err != nil {
		return nil, err
	}

	identifierExists, err := m.dataStore.ExistsByIdentifierAndNodeIdentifier(ctx, request.Identifier, nodeData.Identifier)

	if err != nil {
		return nil, err
	}

	if identifierExists {
		return nil, fmt.Errorf("actor %s already exists", request.Identifier)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	actor := &Actor{
		Identifier:     request.Identifier,
		Password:       string(hashedPassword),
		Type:           request.Type,
		DisplayName:    request.DisplayName,
		NodeIdentifier: nodeData.Identifier,
		Creator:        "",
		CreatedAt:      time.Now().UnixNano(),
		UpdatedAt:      time.Now().UnixNano(),
	}

	return m.dataStore.Insert(ctx, actor)
}

func (m *Manager) GetToken(ctx context.Context, nodeIdentifier string, request *TokenRequest) (*TokenResponse, error) {
	nodeData, err := m.nodeManager.Get(ctx, nodeIdentifier)
	if err != nil {
		return nil, err
	}

	actor, err := m.dataStore.FindByIdentifierAndNodeIdentifier(ctx, request.Identifier, nodeIdentifier)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("invalid identifier and password combination")
	}

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(actor.Password), []byte(request.Password))

	if err != nil {
		return nil, fmt.Errorf("invalid username and password combination")
	}

	token, err := m.authenticator.GenerateToken(actor.Identifier, nodeData, request.TargetNodeAddress)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: token}, nil
}

func (m *Manager) Get(ctx context.Context, identifier string, nodeIdentifier string) (*Actor, error) {
	actor, err := m.dataStore.FindByIdentifierAndNodeIdentifier(ctx, identifier, nodeIdentifier)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("actor %s not found", identifier)
	}

	if err != nil {
		return nil, err
	}

	return actor, nil
}
