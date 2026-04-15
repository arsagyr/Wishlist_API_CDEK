package crypto

import (
    "fmt"
    
    "golang.org/x/crypto/bcrypt"
)

// HashPassword - хэширование пароля
func HashPassword(password string) (string, error) {
    if password == "" {
        return "", fmt.Errorf("password cannot be empty")
    }
    
    // bcrypt.DefaultCost = 10 - хороший баланс безопасности и производительности
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    
    return string(hash), nil
}

// CheckPassword - проверка пароля
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// IsPasswordStrong - проверка сложности пароля
func IsPasswordStrong(password string) bool {
    if len(password) < 8 {
        return false
    }
    
    hasUpper := false
    hasLower := false
    hasNumber := false
    hasSpecial := false
    
    for _, char := range password {
        switch {
        case char >= 'A' && char <= 'Z':
            hasUpper = true
        case char >= 'a' && char <= 'z':
            hasLower = true
        case char >= '0' && char <= '9':
            hasNumber = true
        case char >= '!' && char <= '/' || char >= ':' && char <= '@' || char >= '[' && char <= '`' || char >= '{' && char <= '~':
            hasSpecial = true
        }
    }
    
    return hasUpper && hasLower && hasNumber && hasSpecial
}