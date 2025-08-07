package users_test

import (
	"errors"
	"testing"

	"ministry-scheduler/internal/features/users"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    users.User
		wantErr error
	}{
		{
			name: "valid user",
			user: users.User{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			user: users.User{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: users.ErrEmptyName,
		},
		{
			name: "invalid email format",
			user: users.User{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: users.ErrInvalidEmail,
		},
		{
			name: "empty email",
			user: users.User{
				Name:  "John Doe",
				Email: "",
			},
			wantErr: users.ErrInvalidEmail,
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
		req     users.CreateUserRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: users.CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			req: users.CreateUserRequest{
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: users.ErrEmptyName,
		},
		{
			name: "invalid email",
			req: users.CreateUserRequest{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: users.ErrInvalidEmail,
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
