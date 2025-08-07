package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ministry-scheduler/internal/features/users"
)

type mockUserRepository struct {
	userMap map[int64]*users.User
	nextID  int64
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		userMap: make(map[int64]*users.User),
		nextID:  1,
	}
}

func (m *mockUserRepository) GetByID(_ context.Context, id int64) (*users.User, error) {
	user, exists := m.userMap[id]
	if !exists {
		return nil, users.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(_ context.Context, email string) (*users.User, error) {
	for _, user := range m.userMap {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, users.ErrUserNotFound
}

func (m *mockUserRepository) Create(_ context.Context, user *users.User) (*users.User, error) {
	user.ID = m.nextID
	m.nextID++
	m.userMap[user.ID] = user
	return user, nil
}

func (m *mockUserRepository) Update(_ context.Context, user *users.User) (*users.User, error) {
	if _, exists := m.userMap[user.ID]; !exists {
		return nil, users.ErrUserNotFound
	}
	m.userMap[user.ID] = user
	return user, nil
}

func (m *mockUserRepository) Delete(_ context.Context, id int64) error {
	if _, exists := m.userMap[id]; !exists {
		return users.ErrUserNotFound
	}
	delete(m.userMap, id)
	return nil
}

func (m *mockUserRepository) List(_ context.Context, limit, offset int) ([]*users.User, error) {
	var userList []*users.User
	count := 0
	for _, user := range m.userMap {
		if count >= offset && len(userList) < limit {
			userList = append(userList, user)
		}
		count++
	}
	return userList, nil
}

func TestUserService_CreateUser(t *testing.T) {
	repo := newMockUserRepository()
	service := users.NewService(repo)

	req := &users.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	user, err := service.CreateUser(context.Background(), req)
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

func TestUserService_CreateUserDuplicate(t *testing.T) {
	repo := newMockUserRepository()
	service := users.NewService(repo)

	req := &users.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	_, err := service.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error on first create, got %v", err)
	}

	_, err = service.CreateUser(context.Background(), req)
	if !errors.Is(err, users.ErrUserExists) {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}
}

func TestUserService_GetUser(t *testing.T) {
	repo := newMockUserRepository()
	service := users.NewService(repo)

	originalUser := &users.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.userMap[1] = originalUser

	user, err := service.GetUser(context.Background(), 1)
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

func TestUserService_GetUserNotFound(t *testing.T) {
	repo := newMockUserRepository()
	service := users.NewService(repo)

	_, err := service.GetUser(context.Background(), 999)
	if !errors.Is(err, users.ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	repo := newMockUserRepository()
	service := users.NewService(repo)

	originalUser := &users.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.userMap[1] = originalUser

	newName := "Jane Doe"
	req := &users.UpdateUserRequest{
		Name: &newName,
	}

	updatedUser, err := service.UpdateUser(context.Background(), 1, req)
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

func TestUserService_DeleteUser(t *testing.T) {
	repo := newMockUserRepository()
	service := users.NewService(repo)

	user := &users.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.userMap[1] = user

	err := service.DeleteUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = service.GetUser(context.Background(), 1)
	if !errors.Is(err, users.ErrUserNotFound) {
		t.Errorf("Expected user to be deleted, but got %v", err)
	}
}
