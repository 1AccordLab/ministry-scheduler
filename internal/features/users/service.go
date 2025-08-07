package users

import (
	"context"
	"errors"
	"time"
)

// userService implements the Service interface.
type userService struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	user := &User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed to an existing one
	if req.Email != nil && *req.Email != user.Email {
		existingUser, emailErr := s.repo.GetByEmail(ctx, *req.Email)
		if emailErr != nil && !errors.Is(emailErr, ErrUserNotFound) {
			return nil, emailErr
		}
		if existingUser != nil {
			return nil, ErrUserExists
		}
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	user.UpdatedAt = time.Now()

	if validationErr := user.Validate(); validationErr != nil {
		return nil, validationErr
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

const (
	defaultListLimit = 10
	maxListLimit     = 100
)

func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]*User, error) {
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}
