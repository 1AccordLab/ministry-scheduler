package domain_test

import (
	"errors"
	"testing"

	"ministry-scheduler/internal/domain"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    domain.User
		wantErr error
	}{
		{
			name: "valid user",
			user: domain.User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			user: domain.User{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: domain.ErrEmptyName,
		},
		{
			name: "invalid email format",
			user: domain.User{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: domain.ErrInvalidEmail,
		},
		{
			name: "empty email",
			user: domain.User{
				Name:  "John Doe",
				Email: "",
			},
			wantErr: domain.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr == nil && err != nil {
				t.Errorf("User.Validate() error = %v, wantErr nil", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     domain.CreateUserRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: domain.CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			req: domain.CreateUserRequest{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: domain.ErrEmptyName,
		},
		{
			name: "invalid email",
			req: domain.CreateUserRequest{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: domain.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == nil && err != nil {
				t.Errorf("CreateUserRequest.Validate() error = %v, wantErr nil", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateUserRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
