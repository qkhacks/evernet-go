package messaging

import (
	"context"
	"fmt"
	"time"
)

type InboxManager struct {
	dataStore *InboxDataStore
}

func NewInboxManager(dataStore *InboxDataStore) *InboxManager {
	return &InboxManager{dataStore: dataStore}
}

func (m *InboxManager) Create(ctx context.Context, request *InboxCreationRequest, actorAddress string, nodeIdentifier string) (*Inbox, error) {
	identifierExists, err := m.dataStore.ExistsByIdentifierAndNodeIdentifier(ctx, request.Identifier, nodeIdentifier)

	if err != nil {
		return nil, err
	}

	if identifierExists {
		return nil, fmt.Errorf("inbox %s already exists", request.Identifier)
	}

	inbox := &Inbox{
		Identifier:     request.Identifier,
		DisplayName:    request.DisplayName,
		NodeIdentifier: nodeIdentifier,
		ActorAddress:   actorAddress,
		CreatedAt:      time.Now().UnixNano(),
		UpdatedAt:      time.Now().UnixNano(),
	}

	return m.dataStore.Insert(ctx, inbox)
}
