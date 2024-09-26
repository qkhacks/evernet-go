package admin

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Manager struct {
	dataStore *DataStore
}

func NewManager(dataStore *DataStore) *Manager {
	return &Manager{dataStore: dataStore}
}

func (m *Manager) Init(ctx context.Context, request *InitRequest) (*Admin, error) {

	exists, err := m.dataStore.Exists(ctx)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("not allowed")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	admin := &Admin{
		Identifier: request.Identifier,
		Password:   string(hashedPassword),
		Creator:    "",
		CreatedAt:  time.Now().UnixNano(),
		UpdatedAt:  time.Now().UnixNano(),
	}

	return m.dataStore.Insert(ctx, admin)
}
