package dto

import (
	"time"

	"test_cdek/backend/internal/model"
)

type CreateWishlistRequest struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Description string     `json:"description"`
	EventDate   *time.Time `json:"event_date"`
}

type UpdateWishlistRequest struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Description string     `json:"description"`
	EventDate   *time.Time `json:"event_date"`
}

type WishlistResponse struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	EventDate   *time.Time     `json:"event_date,omitempty"`
	PublicToken string         `json:"public_token"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Items       []ItemResponse `json:"items,omitempty"`
}

func ToWishlistResponse(w *model.Wishlist) *WishlistResponse {
	resp := &WishlistResponse{
		ID:          w.ID,
		Title:       w.Title,
		Description: w.Description,
		EventDate:   w.EventDate,
		PublicToken: w.PublicToken,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
	return resp
}

func ToWishlistWithItemsResponse(w *model.Wishlist, items []*model.WishlistItem) *WishlistResponse {
	resp := ToWishlistResponse(w)
	resp.Items = make([]ItemResponse, len(items))
	for i, item := range items {
		resp.Items[i] = *ToItemResponse(item)
	}
	return resp
}

func (r *CreateWishlistRequest) ToModel(userID string) (*model.Wishlist, error) {
	return model.NewWishlist(userID, r.Title, r.Description, r.EventDate)
}

type ItemResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Priority    int       `json:"priority"`
	IsReserved  bool      `json:"is_reserved"`
	ReservedBy  *string   `json:"reserved_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToItemResponse(i *model.WishlistItem) *ItemResponse {
	return &ItemResponse{
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		Link:        i.Link,
		Priority:    int(i.Priority),
		IsReserved:  i.IsReserved,
		ReservedBy:  i.ReservedBy,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

type CreateItemRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    int    `json:"priority" validate:"min=0,max=2"`
}

func (r *CreateItemRequest) ToModel(wishlistID string) (*model.WishlistItem, error) {
	return model.NewWishlistItem(wishlistID, r.Name, r.Description, r.Link, model.Priority(r.Priority))
}

type UpdateItemRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    int    `json:"priority" validate:"min=0,max=2"`
}

func (r *UpdateItemRequest) ToModel(item *model.WishlistItem) error {
	return item.Update(r.Name, r.Description, r.Link, model.Priority(r.Priority))
}

type ReserveItemRequest struct {
	ItemID     string `json:"item_id" validate:"required"`
	ReservedBy string `json:"reserved_by" validate:"required"`
}
