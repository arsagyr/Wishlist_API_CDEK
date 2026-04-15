package model

import (
	"errors"
	"time"

	"test_cdek/backend/pkg/crypto"
)

// User - доменная сущность пользователя
type User struct {
	ID            string     `db:"id"`
	Email         string     `db:"email"`
	PasswordHash  string     `db:"password_hash"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	LastLoginAt   *time.Time `db:"last_login_at"`
	IsActive      bool       `db:"is_active"`
	EmailVerified bool       `db:"email_verified"`
}

// NewUser - конструктор нового пользователя
func NewUser(email, plainPassword string) (*User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	if !crypto.IsPasswordStrong(plainPassword) {
		return nil, errors.New("password is not strong enough: min 8 chars, need upper, lower, number and special character")
	}

	hash, err := crypto.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Email:         email,
		PasswordHash:  hash,
		CreatedAt:     now,
		UpdatedAt:     now,
		IsActive:      true,
		EmailVerified: false,
	}, nil
}

// isValidEmail - простая валидация email
func isValidEmail(email string) bool {
	// Базовая проверка - должна содержать @ и .
	atIndex := -1
	dotIndex := -1

	for i, ch := range email {
		if ch == '@' {
			atIndex = i
		}
		if ch == '.' && atIndex != -1 && i > atIndex+1 {
			dotIndex = i
		}
	}

	return atIndex > 0 && dotIndex > atIndex+1 && dotIndex < len(email)-1
}

// UpdatePassword - обновление пароля
func (u *User) UpdatePassword(newPassword string) error {
	if !crypto.IsPasswordStrong(newPassword) {
		return errors.New("password is not strong enough")
	}

	hash, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.PasswordHash = hash
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail - обновление email
func (u *User) UpdateEmail(newEmail string) error {
	if newEmail == "" {
		return errors.New("email cannot be empty")
	}

	if !isValidEmail(newEmail) {
		return errors.New("invalid email format")
	}

	u.Email = newEmail
	u.EmailVerified = false // При смене email нужно подтверждение заново
	u.UpdatedAt = time.Now()
	return nil
}

// VerifyEmail - подтверждение email
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now()
}

// Activate - активация пользователя
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate - деактивация пользователя
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// RecordLogin - запись времени входа
func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// ValidatePassword - проверка пароля
func (u *User) ValidatePassword(password string) bool {
	return crypto.CheckPassword(password, u.PasswordHash)
}

// IsActiveUser - проверка активности
func (u *User) IsActiveUser() bool {
	return u.IsActive
}

// IsEmailConfirmed - проверка подтверждения email
func (u *User) IsEmailConfirmed() bool {
	return u.EmailVerified
}
