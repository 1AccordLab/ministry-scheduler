package usecase

import (
	"context"
	"errors"
	"time"

	"ministry-scheduler/internal/domain"
)

type UserUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (u *UserUsecase) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UserUsecase) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingUser, err := u.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserExists
	}

	user := &domain.User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.repo.Create(ctx, user)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed to an existing one
	if req.Email != nil && *req.Email != user.Email {
		existingUser, emailErr := u.repo.GetByEmail(ctx, *req.Email)
		if emailErr != nil && !errors.Is(emailErr, domain.ErrUserNotFound) {
			return nil, emailErr
		}
		if existingUser != nil {
			return nil, domain.ErrUserExists
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

	return u.repo.Update(ctx, user)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return u.repo.Delete(ctx, id)
}

const (
	defaultListLimit = 10
	maxListLimit     = 100
)

func (u *UserUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if offset < 0 {
		offset = 0
	}

	return u.repo.List(ctx, limit, offset)
}
