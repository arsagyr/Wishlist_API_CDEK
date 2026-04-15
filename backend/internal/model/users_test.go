package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid", "test@example.com", false},
		{"invalid email", "invalid-email", true},
		{"empty email", "", true},
		{"weak password", "test@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var password string
			if tt.name == "weak password" {
				password = "weak"
			} else {
				password = "Password123!"
			}
			_, err := NewUser(tt.email, password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	user := &User{
		Email:        "test@example.com",
		PasswordHash: "$2a$10$abcdefghijklmnopqrstuvwxyz12345678901234567890",
	}

	assert.False(t, user.ValidatePassword("wrong"))
}

func TestUser_IsActiveUser(t *testing.T) {
	tests := []struct {
		name     string
		isActive bool
		expected bool
	}{
		{"active", true, true},
		{"inactive", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{IsActive: tt.isActive}
			assert.Equal(t, tt.expected, user.IsActiveUser())
		})
	}
}
