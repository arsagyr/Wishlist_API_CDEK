package handlers

import (
	"encoding/json"
	"net/http"

	"test_cdek/backend/internal/dto"
	"test_cdek/backend/internal/repository"
)

type PublicHandler struct {
	wishlistRepo *repository.WishlistRepository
	itemRepo     *repository.WishlistItemRepository
}

func NewPublicHandler(wishlistRepo *repository.WishlistRepository, itemRepo *repository.WishlistItemRepository) *PublicHandler {
	return &PublicHandler{
		wishlistRepo: wishlistRepo,
		itemRepo:     itemRepo,
	}
}

func (h *PublicHandler) GetWishlistByToken(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	wishlist, err := h.wishlistRepo.GetByPublicToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	items, err := h.itemRepo.GetByWishlistID(r.Context(), wishlist.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToWishlistWithItemsResponse(wishlist, items))
}

func (h *PublicHandler) ReserveItem(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	wishlist, err := h.wishlistRepo.GetByPublicToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	var req dto.ReserveItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.itemRepo.GetByID(r.Context(), req.ItemID)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if item.WishlistID != wishlist.ID {
		http.Error(w, "Item not found in this wishlist", http.StatusNotFound)
		return
	}

	if err := h.itemRepo.Reserve(r.Context(), req.ItemID, req.ReservedBy); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Item reserved successfully",
	})
}
