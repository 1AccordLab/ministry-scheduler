package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ministry-scheduler/internal/domain"
	"ministry-scheduler/internal/usecase"
)

type mockUserRepository struct {
	users  map[int64]*domain.User
	nextID int64
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) GetByID(_ context.Context, id int64) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) Create(_ context.Context, user *domain.User) (*domain.User, error) {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return user, nil
}

func (m *mockUserRepository) Update(_ context.Context, user *domain.User) (*domain.User, error) {
	if _, exists := m.users[user.ID]; !exists {
		return nil, domain.ErrUserNotFound
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *mockUserRepository) Delete(_ context.Context, id int64) error {
	if _, exists := m.users[id]; !exists {
		return domain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) List(_ context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	count := 0
	for _, user := range m.users {
		if count >= offset && len(users) < limit {
			users = append(users, user)
		}
		count++
	}
	return users, nil
}

func TestUserUsecase_CreateUser(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	req := &domain.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	user, err := uc.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, user.Name)
	}

	if user.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, user.Email)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserUsecase_CreateUserDuplicate(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	req := &domain.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	_, err := uc.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error on first create, got %v", err)
	}

	_, err = uc.CreateUser(context.Background(), req)
	if !errors.Is(err, domain.ErrUserExists) {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	originalUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.users[1] = originalUser

	user, err := uc.GetUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID != originalUser.ID {
		t.Errorf("Expected ID %d, got %d", originalUser.ID, user.ID)
	}

	if user.Name != originalUser.Name {
		t.Errorf("Expected name %s, got %s", originalUser.Name, user.Name)
	}
}

func TestUserUsecase_GetUserNotFound(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	_, err := uc.GetUser(context.Background(), 999)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	originalUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.users[1] = originalUser

	newName := "Jane Doe"
	req := &domain.UpdateUserRequest{
		Name: &newName,
	}

	updatedUser, err := uc.UpdateUser(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedUser.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, updatedUser.Name)
	}

	if updatedUser.Email != originalUser.Email {
		t.Errorf("Expected email to remain %s, got %s", originalUser.Email, updatedUser.Email)
	}
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	repo := newMockUserRepository()
	uc := usecase.NewUserUsecase(repo)

	user := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.users[1] = user

	err := uc.DeleteUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = uc.GetUser(context.Background(), 1)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("Expected user to be deleted, but got %v", err)
	}
}
