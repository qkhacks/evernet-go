package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Manager struct {
	dataStore    *DataStore
	authenicator *Authenticator
}

func NewManager(dataStore *DataStore, authenticator *Authenticator) *Manager {
	return &Manager{
		dataStore:    dataStore,
		authenicator: authenticator,
	}
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

func (m *Manager) GetToken(ctx context.Context, request *TokenRequest) (*TokenResponse, error) {

	admin, err := m.dataStore.FindByIdentifier(ctx, request.Identifier)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("invalid identifier and password combination")
	}

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password))

	if err != nil {
		return nil, fmt.Errorf("invalid identifier and password combination")
	}

	token, err := m.authenicator.GenerateToken(admin.Identifier)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: token}, nil
}

func (m *Manager) Get(ctx context.Context, identifier string) (*Admin, error) {

	admin, err := m.dataStore.FindByIdentifier(ctx, identifier)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("admin not found")
	}

	if err != nil {
		return nil, err
	}

	return admin, nil
}
