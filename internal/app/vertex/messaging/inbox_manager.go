package messaging

import (
	"context"
	"database/sql"
	"errors"
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

func (m *InboxManager) List(ctx context.Context, actorAddress string, nodeIdentifier string, page int64, size int64) ([]*Inbox, error) {
	return m.dataStore.FindByActorAddressAndNodeIdentifier(ctx, actorAddress, nodeIdentifier, page, size)
}

func (m *InboxManager) Get(ctx context.Context, identifier string, actorAddress string, nodeIdentifier string) (*Inbox, error) {
	inbox, err := m.dataStore.FindByIdentifierAndActorAddressAndNodeIdentifier(ctx, identifier, actorAddress, nodeIdentifier)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("inbox %s not found", identifier)
	}

	if err != nil {
		return nil, err
	}

	return inbox, nil
}

func (m *InboxManager) Update(ctx context.Context, identifier string, request *InboxUpdateRequest, actorAddress string, nodeIdentifier string) error {
	err := m.dataStore.UpdateDisplayNameByIdentifierAndActorAddressAndNodeIdentifier(ctx, request.DisplayName, identifier, actorAddress, nodeIdentifier)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("inbox %s not found", identifier)
	}

	return err
}

func (m *InboxManager) Delete(ctx context.Context, identifier string, actorAddress string, nodeIdentifier string) error {
	err := m.dataStore.DeleteByIdentifierAndActorAddressAndNodeIdentifier(ctx, identifier, actorAddress, nodeIdentifier)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("inbox %s not found", identifier)
	}

	return err
}
