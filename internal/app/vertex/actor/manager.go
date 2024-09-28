package actor

import (
	"context"
	"fmt"
	"github.com/evernetproto/evernet/internal/app/vertex/node"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Manager struct {
	dataStore   *DataStore
	nodeManager *node.Manager
}

func NewManager(dataStore *DataStore, nodeManager *node.Manager) *Manager {
	return &Manager{dataStore: dataStore, nodeManager: nodeManager}
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
