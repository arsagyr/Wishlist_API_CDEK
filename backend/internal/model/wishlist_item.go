package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

type WishlistItem struct {
	ID          string    `db:"id"`
	WishlistID  string    `db:"wishlist_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Link        string    `db:"link"`
	Priority    Priority  `db:"priority"`
	IsReserved  bool      `db:"is_reserved"`
	ReservedBy  *string   `db:"reserved_by"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewWishlistItem(wishlistID, name, description, link string, priority Priority) (*WishlistItem, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if wishlistID == "" {
		return nil, errors.New("wishlist_id cannot be empty")
	}

	now := time.Now()
	return &WishlistItem{
		ID:          uuid.New().String(),
		WishlistID:  wishlistID,
		Name:        name,
		Description: description,
		Link:        link,
		Priority:    priority,
		IsReserved:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (i *WishlistItem) Update(name, description, link string, priority Priority) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	i.Name = name
	i.Description = description
	i.Link = link
	i.Priority = priority
	i.UpdatedAt = time.Now()
	return nil
}

func (i *WishlistItem) Reserve(reservedBy string) error {
	if i.IsReserved {
		return errors.New("item already reserved")
	}
	i.IsReserved = true
	i.ReservedBy = &reservedBy
	i.UpdatedAt = time.Now()
	return nil
}

func (i *WishlistItem) Unreserve() {
	i.IsReserved = false
	i.ReservedBy = nil
	i.UpdatedAt = time.Now()
}
