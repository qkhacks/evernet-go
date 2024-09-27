package node

import (
	"context"
	"fmt"
	"github.com/evernetproto/evernet/internal/pkg/keys"
	"time"
)

type Manager struct {
	dataStore *DataStore
}

func NewManager(dataStore *DataStore) *Manager {
	return &Manager{dataStore: dataStore}
}

func (m *Manager) Create(ctx context.Context, request *CreationRequest, creator string) (*Node, error) {
	identifierExists, err := m.dataStore.ExistsByIdentifier(ctx, request.Identifier)

	if err != nil {
		return nil, err
	}

	if identifierExists {
		return nil, fmt.Errorf("node %s already exists", request.Identifier)
	}

	signingPublicKey, signingPrivateKey, err := keys.GenerateED25519KeyPair()

	if err != nil {
		return nil, err
	}

	node := &Node{
		Identifier:        request.Identifier,
		DisplayName:       request.DisplayName,
		SigningPrivateKey: keys.ConvertED25519PrivateKeyToString(signingPrivateKey),
		SigningPublicKey:  keys.ConvertED25519PublicKeyToString(signingPublicKey),
		Creator:           creator,
		CreatedAt:         time.Now().UnixNano(),
		UpdatedAt:         time.Now().UnixNano(),
	}

	return m.dataStore.Insert(ctx, node)
}

func (m *Manager) List(ctx context.Context, page int64, size int64) ([]*Node, error) {
	return m.dataStore.FindAll(ctx, page, size)
}
