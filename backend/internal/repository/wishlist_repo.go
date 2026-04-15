package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"test_cdek/backend/internal/model"
)

type WishlistRepository struct {
	db *sql.DB
}

func NewWishlistRepository(db *sql.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Create(ctx context.Context, w *model.Wishlist) error {
	query := `
        INSERT INTO wishlists (id, user_id, title, description, event_date, public_token, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		w.ID,
		w.UserID,
		w.Title,
		w.Description,
		w.EventDate,
		w.PublicToken,
		w.CreatedAt,
		w.UpdatedAt,
	).Scan(&w.ID)

	if err != nil {
		return fmt.Errorf("failed to create wishlist: %w", err)
	}

	return nil
}

func (r *WishlistRepository) GetByID(ctx context.Context, id string) (*model.Wishlist, error) {
	query := `
        SELECT id, user_id, title, description, event_date, public_token, created_at, updated_at
        FROM wishlists
        WHERE id = $1
    `

	var w model.Wishlist
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID,
		&w.UserID,
		&w.Title,
		&w.Description,
		&w.EventDate,
		&w.PublicToken,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("wishlist with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	return &w, nil
}

func (r *WishlistRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Wishlist, error) {
	query := `
        SELECT id, user_id, title, description, event_date, public_token, created_at, updated_at
        FROM wishlists
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list wishlists: %w", err)
	}
	defer rows.Close()

	var wishlists []*model.Wishlist
	for rows.Next() {
		var w model.Wishlist
		err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Title,
			&w.Description,
			&w.EventDate,
			&w.PublicToken,
			&w.CreatedAt,
			&w.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan wishlist: %w", err)
		}
		wishlists = append(wishlists, &w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating wishlists: %w", err)
	}

	return wishlists, nil
}

func (r *WishlistRepository) GetByPublicToken(ctx context.Context, token string) (*model.Wishlist, error) {
	query := `
        SELECT id, user_id, title, description, event_date, public_token, created_at, updated_at
        FROM wishlists
        WHERE public_token = $1
    `

	var w model.Wishlist
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&w.ID,
		&w.UserID,
		&w.Title,
		&w.Description,
		&w.EventDate,
		&w.PublicToken,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("wishlist not found")
		}
		return nil, fmt.Errorf("failed to get wishlist by token: %w", err)
	}

	return &w, nil
}

func (r *WishlistRepository) Update(ctx context.Context, w *model.Wishlist) error {
	query := `
        UPDATE wishlists
        SET title = $1, description = $2, event_date = $3, updated_at = $4
        WHERE id = $5
        RETURNING updated_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		w.Title,
		w.Description,
		w.EventDate,
		w.UpdatedAt,
		w.ID,
	).Scan(&w.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update wishlist: %w", err)
	}

	return nil
}

func (r *WishlistRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM wishlists WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("wishlist with id %s not found", id)
	}

	return nil
}

type WishlistItemRepository struct {
	db *sql.DB
}

func NewWishlistItemRepository(db *sql.DB) *WishlistItemRepository {
	return &WishlistItemRepository{db: db}
}

func (r *WishlistItemRepository) Create(ctx context.Context, item *model.WishlistItem) error {
	query := `
        INSERT INTO wishlist_items (id, wishlist_id, name, description, link, priority, is_reserved, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		item.ID,
		item.WishlistID,
		item.Name,
		item.Description,
		item.Link,
		item.Priority,
		item.IsReserved,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(&item.ID)

	if err != nil {
		return fmt.Errorf("failed to create wishlist item: %w", err)
	}

	return nil
}

func (r *WishlistItemRepository) GetByID(ctx context.Context, id string) (*model.WishlistItem, error) {
	query := `
        SELECT id, wishlist_id, name, description, link, priority, is_reserved, reserved_by, created_at, updated_at
        FROM wishlist_items
        WHERE id = $1
    `

	var item model.WishlistItem
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.WishlistID,
		&item.Name,
		&item.Description,
		&item.Link,
		&item.Priority,
		&item.IsReserved,
		&item.ReservedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}

func (r *WishlistItemRepository) GetByWishlistID(ctx context.Context, wishlistID string) ([]*model.WishlistItem, error) {
	query := `
        SELECT id, wishlist_id, name, description, link, priority, is_reserved, reserved_by, created_at, updated_at
        FROM wishlist_items
        WHERE wishlist_id = $1
        ORDER BY priority DESC, created_at ASC
    `

	rows, err := r.db.QueryContext(ctx, query, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	var items []*model.WishlistItem
	for rows.Next() {
		var item model.WishlistItem
		err := rows.Scan(
			&item.ID,
			&item.WishlistID,
			&item.Name,
			&item.Description,
			&item.Link,
			&item.Priority,
			&item.IsReserved,
			&item.ReservedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating items: %w", err)
	}

	return items, nil
}

func (r *WishlistItemRepository) Update(ctx context.Context, item *model.WishlistItem) error {
	query := `
        UPDATE wishlist_items
        SET name = $1, description = $2, link = $3, priority = $4, updated_at = $5
        WHERE id = $6
        RETURNING updated_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		item.Name,
		item.Description,
		item.Link,
		item.Priority,
		item.UpdatedAt,
		item.ID,
	).Scan(&item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (r *WishlistItemRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM wishlist_items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item with id %s not found", id)
	}

	return nil
}

func (r *WishlistItemRepository) Reserve(ctx context.Context, itemID, reservedBy string) error {
	query := `
        UPDATE wishlist_items
        SET is_reserved = true, reserved_by = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2 AND is_reserved = false
        RETURNING id
    `

	var id string
	err := r.db.QueryRowContext(ctx, query, reservedBy, itemID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("item already reserved or not found")
		}
		return fmt.Errorf("failed to reserve item: %w", err)
	}

	return nil
}
