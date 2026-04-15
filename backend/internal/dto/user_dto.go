// internal/api/dto/user_dto.go
package dto

import (
	"test_cdek/backend/internal/model"
	"time"
)

// RegisterRequest - DTO для регистрации
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest - DTO для входа
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateEmailRequest - DTO для обновления email
type UpdateEmailRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdatePasswordRequest - DTO для обновления пароля
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UserResponse - DTO для ответа API (без пароля)
type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
}

// LoginResponse - DTO для ответа при входе
type LoginResponse struct {
	User      UserResponse `json:"user"`
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// ToUserResponse - конвертация из model в response DTO
func ToUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
	}
}

// ToModel - конвертация из регистрационного DTO в model
func (r *RegisterRequest) ToModel() (*model.User, error) {
	return model.NewUser(r.Email, r.Password)
}

// ProfileResponse - DTO для профиля (без чувствительных данных)
type ProfileResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"registered_at"`
	EmailVerified bool      `json:"email_verified"`
}

// ToProfileResponse - конвертация в профиль
func ToProfileResponse(user *model.User) *ProfileResponse {
	return &ProfileResponse{
		ID:            user.ID,
		Email:         user.Email,
		CreatedAt:     user.CreatedAt,
		EmailVerified: user.EmailVerified,
	}
}
