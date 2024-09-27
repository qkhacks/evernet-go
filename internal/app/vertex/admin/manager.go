package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sethvargo/go-password/password"
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

func (m *Manager) ChangePassword(ctx context.Context, identifier string, request *PasswordChangeRequest) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	err = m.dataStore.UpdatePasswordByIdentifier(ctx, string(hashedPassword), identifier)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("admin %s not found", identifier)
	}

	return err
}

func (m *Manager) Add(ctx context.Context, request *AdditionRequest, creator string) (*AdditionResponse, error) {
	identifierExists, err := m.dataStore.ExistsByIdentifier(ctx, request.Identifier)

	if err != nil {
		return nil, err
	}

	if identifierExists {
		return nil, fmt.Errorf("admin %s already exists", request.Identifier)
	}

	newPassword, err := password.Generate(16, 4, 2, false, false)

	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &Admin{
		Identifier: request.Identifier,
		Password:   string(hashedPassword),
		Creator:    creator,
		CreatedAt:  time.Now().UnixNano(),
		UpdatedAt:  time.Now().UnixNano(),
	}

	admin, err = m.dataStore.Insert(ctx, admin)

	if err != nil {
		return nil, err
	}

	return &AdditionResponse{
		Password: newPassword,
		Admin:    admin,
	}, nil
}

func (m *Manager) Delete(ctx context.Context, identifier string) error {
	err := m.dataStore.DeleteByIdentifier(ctx, identifier)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("admin %s not found", identifier)
	}

	return err
}

func (m *Manager) List(ctx context.Context, page int64, size int64) ([]*Admin, error) {
	return m.dataStore.FindAll(ctx, page, size)
}

func (m *Manager) ResetPassword(ctx context.Context, identifier string) (*PasswordResponse, error) {
	newPassword, err := password.Generate(16, 4, 2, false, false)
	if err != nil {
		return nil, err
	}

	err = m.ChangePassword(ctx, identifier, &PasswordChangeRequest{
		Password: newPassword,
	})

	if err != nil {
		return nil, err
	}

	return &PasswordResponse{
		Password: newPassword,
	}, nil
}
