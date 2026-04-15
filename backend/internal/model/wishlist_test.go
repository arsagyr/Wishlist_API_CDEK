package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWishlist(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		title       string
		description string
		eventDate   *time.Time
		wantErr     bool
	}{
		{
			name:        "valid",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			title:       "Birthday",
			description: "My birthday wishlist",
			eventDate:   nil,
			wantErr:     false,
		},
		{
			name:    "empty userID",
			userID:  "",
			title:   "Birthday",
			wantErr: true,
		},
		{
			name:    "empty title",
			userID:  "550e8400-e29b-41d4-a716-446655440000",
			title:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWishlist(tt.userID, tt.title, tt.description, tt.eventDate)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, w)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, w)
				assert.NotEmpty(t, w.PublicToken)
				assert.NotEmpty(t, w.ID)
			}
		})
	}
}

func TestWishlist_Update(t *testing.T) {
	w, _ := NewWishlist("user-id", "Old Title", "Old Desc", nil)

	err := w.Update("New Title", "New Desc", nil)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", w.Title)
	assert.Equal(t, "New Desc", w.Description)

	err = w.Update("", "New Desc", nil)
	assert.Error(t, err)
}

func TestWishlist_RegenerateToken(t *testing.T) {
	w, _ := NewWishlist("user-id", "Title", "Desc", nil)
	oldToken := w.PublicToken

	w.RegenerateToken()
	assert.NotEqual(t, oldToken, w.PublicToken)
}

func TestNewWishlistItem(t *testing.T) {
	tests := []struct {
		name       string
		wishlistID string
		nameValue  string
		priority   Priority
		wantErr    bool
	}{
		{
			name:       "valid",
			wishlistID: "wishlist-id",
			nameValue:  "iPhone 15",
			priority:   PriorityHigh,
			wantErr:    false,
		},
		{
			name:       "empty wishlistID",
			wishlistID: "",
			nameValue:  "iPhone 15",
			wantErr:    true,
		},
		{
			name:       "empty name",
			wishlistID: "wishlist-id",
			nameValue:  "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewWishlistItem(tt.wishlistID, tt.nameValue, "", "", tt.priority)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, item.ID)
			}
		})
	}
}

func TestWishlistItem_Reserve(t *testing.T) {
	item, _ := NewWishlistItem("wishlist-id", "Gift", "", "", PriorityMedium)

	err := item.Reserve("John")
	assert.NoError(t, err)
	assert.True(t, item.IsReserved)
	assert.Equal(t, "John", *item.ReservedBy)

	err = item.Reserve("Jane")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already reserved")
}

func TestWishlistItem_Unreserve(t *testing.T) {
	item, _ := NewWishlistItem("wishlist-id", "Gift", "", "", PriorityMedium)
	item.Reserve("John")

	item.Unreserve()
	assert.False(t, item.IsReserved)
	assert.Nil(t, item.ReservedBy)
}
