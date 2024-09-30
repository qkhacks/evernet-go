package messaging

import (
	"context"
	"fmt"
	"time"
)

type OutboxManager struct {
	dataStore *OutboxDataStore
}

func NewOutboxManager(dataStore *OutboxDataStore) *OutboxManager {
	return &OutboxManager{dataStore: dataStore}
}

func (m *OutboxManager) Create(ctx context.Context, request *OutboxCreationRequest, actorAddress string, nodeIdentifier string) (*Outbox, error) {
	identifierExists, err := m.dataStore.ExistsByIdentifierAndNodeIdentifier(ctx, request.Identifier, nodeIdentifier)

	if err != nil {
		return nil, err
	}

	if identifierExists {
		return nil, fmt.Errorf("outbox %s already exists", request.Identifier)
	}

	outbox := &Outbox{
		Identifier:     request.Identifier,
		DisplayName:    request.DisplayName,
		NodeIdentifier: nodeIdentifier,
		ActorAddress:   actorAddress,
		CreatedAt:      time.Now().UnixNano(),
		UpdatedAt:      time.Now().UnixNano(),
	}

	return m.dataStore.Insert(ctx, outbox)
}

func (m *OutboxManager) List(ctx context.Context, actorAddress string, nodeIdentifier string, page int64, size int64) ([]*Outbox, error) {
	return m.dataStore.FindByActorAddressAndNodeIdentifier(ctx, actorAddress, nodeIdentifier, page, size)
}
