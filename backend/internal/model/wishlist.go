package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	EventDate   *time.Time `db:"event_date"`
	PublicToken string     `db:"public_token"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func NewWishlist(userID, title, description string, eventDate *time.Time) (*Wishlist, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if userID == "" {
		return nil, errors.New("user_id cannot be empty")
	}

	now := time.Now()
	return &Wishlist{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: description,
		EventDate:   eventDate,
		PublicToken: uuid.New().String(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (w *Wishlist) Update(title, description string, eventDate *time.Time) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	w.Title = title
	w.Description = description
	w.EventDate = eventDate
	w.UpdatedAt = time.Now()
	return nil
}

func (w *Wishlist) RegenerateToken() {
	w.PublicToken = uuid.New().String()
	w.UpdatedAt = time.Now()
}
