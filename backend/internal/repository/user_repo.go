// internal/repository/user_repo.go
package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    
    "your-project/internal/domain/entity"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create - создание пользователя
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
    query := `
        INSERT INTO users (id, email, password_hash, created_at, updated_at, is_active, email_verified)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
    
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.ID,
        user.Email,
        user.PasswordHash,
        user.CreatedAt,
        user.UpdatedAt,
        user.IsActive,
        user.EmailVerified,
    ).Scan(&user.ID)
    
    if err != nil {
        // Проверка на дубликат email
        if isDuplicateKeyError(err) {
            return fmt.Errorf("user with email %s already exists", user.Email)
        }
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

// GetByID - получение пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
    query := `
        SELECT id, email, password_hash, created_at, updated_at, last_login_at, is_active, email_verified
        FROM users
        WHERE id = $1
    `
    
    var user entity.User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.LastLoginAt,
        &user.IsActive,
        &user.EmailVerified,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user with id %s not found", id)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

// GetByEmail - получение пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    query := `
        SELECT id, email, password_hash, created_at, updated_at, last_login_at, is_active, email_verified
        FROM users
        WHERE email = $1
    `
    
    var user entity.User
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.LastLoginAt,
        &user.IsActive,
        &user.EmailVerified,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user with email %s not found", email)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

// Update - обновление пользователя
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
    query := `
        UPDATE users
        SET email = $1, password_hash = $2, updated_at = $3, 
            last_login_at = $4, is_active = $5, email_verified = $6
        WHERE id = $7
        RETURNING updated_at
    `
    
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Email,
        user.PasswordHash,
        user.UpdatedAt,
        user.LastLoginAt,
        user.IsActive,
        user.EmailVerified,
        user.ID,
    ).Scan(&user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    return nil
}

// UpdateLastLogin - обновление времени последнего входа
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string, lastLoginAt *entity.Time) error {
    query := `UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3`
    
    now := time.Now()
    _, err := r.db.ExecContext(ctx, query, lastLoginAt, now, userID)
    if err != nil {
        return fmt.Errorf("failed to update last login: %w", err)
    }
    
    return nil
}

// Delete - удаление пользователя (мягкое удаление через деактивацию)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
    query := `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`
    
    result, err := r.db.ExecContext(ctx, query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user with id %s not found", id)
    }
    
    return nil
}

// HardDelete - полное удаление из БД (использовать осторожно!)
func (r *UserRepository) HardDelete(ctx context.Context, id string) error {
    query := `DELETE FROM users WHERE id = $1`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to hard delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user with id %s not found", id)
    }
    
    return nil
}

// List - получение списка пользователей с пагинацией
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
    query := `
        SELECT id, email, password_hash, created_at, updated_at, last_login_at, is_active, email_verified
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    defer rows.Close()
    
    var users []*entity.User
    for rows.Next() {
        var user entity.User
        err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.PasswordHash,
            &user.CreatedAt,
            &user.UpdatedAt,
            &user.LastLoginAt,
            &user.IsActive,
            &user.EmailVerified,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, &user)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating users: %w", err)
    }
    
    return users, nil
}

// CheckEmailExists - проверка существования email
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    
    var exists bool
    err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check email existence: %w", err)
    }
    
    return exists, nil
}

// isDuplicateKeyError - проверка ошибки дубликата (зависит от драйвера)
func isDuplicateKeyError(err error) bool {
    // Для lib/pq
    if pqErr, ok := err.(*pq.Error); ok {
        return pqErr.Code == "23505" // unique_violation
    }
    return false
}