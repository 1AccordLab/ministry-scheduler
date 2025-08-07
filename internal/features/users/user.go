package users

import (
	"context"
	"errors"
	"regexp"
	"time"
)

// User represents the domain entity for user management.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the data needed to create a new user.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest represents the data needed to update a user.
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// User domain errors.
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidUserData = errors.New("invalid user data")
	ErrEmptyName       = errors.New("name cannot be empty")
	ErrInvalidEmail    = errors.New("invalid email format")
)

// Validation constants.
const (
	MinNameLength  = 1
	MaxNameLength  = 100
	MinEmailLength = 5
	MaxEmailLength = 254
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Repository defines the interface for user data persistence.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

// Service defines the business logic interface for users.
type Service interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
}

// Validate checks if the user data is valid.
func (u *User) Validate() error {
	if err := validateName(u.Name); err != nil {
		return err
	}
	if err := validateEmail(u.Email); err != nil {
		return err
	}
	return nil
}

func (req *CreateUserRequest) Validate() error {
	if err := validateName(req.Name); err != nil {
		return err
	}
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	return nil
}

func validateName(name string) error {
	if len(name) < MinNameLength {
		return ErrEmptyName
	}
	if len(name) > MaxNameLength {
		return errors.New("name is too long")
	}
	return nil
}

func validateEmail(email string) error {
	if len(email) < MinEmailLength {
		return ErrInvalidEmail
	}
	if len(email) > MaxEmailLength {
		return errors.New("email is too long")
	}
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}
