package handlers

import (
	"encoding/json"
	"net/http"

	"test_cdek/backend/internal/dto"
	"test_cdek/backend/internal/middleware"
	"test_cdek/backend/internal/repository"
)

type WishlistItemHandler struct {
	wishlistRepo *repository.WishlistRepository
	itemRepo     *repository.WishlistItemRepository
}

func NewWishlistItemHandler(wishlistRepo *repository.WishlistRepository, itemRepo *repository.WishlistItemRepository) *WishlistItemHandler {
	return &WishlistItemHandler{
		wishlistRepo: wishlistRepo,
		itemRepo:     itemRepo,
	}
}

func (h *WishlistItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("wishlistId")

	wishlist, err := h.wishlistRepo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req dto.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item, err := req.ToModel(wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.itemRepo.Create(r.Context(), item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToItemResponse(item))
}

func (h *WishlistItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("wishlistId")
	itemID := r.PathValue("id")

	wishlist, err := h.wishlistRepo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	item, err := h.itemRepo.GetByID(r.Context(), itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if item.WishlistID != wishlistID {
		http.Error(w, "Item not found in this wishlist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToItemResponse(item))
}

func (h *WishlistItemHandler) ListByWishlist(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("wishlistId")

	wishlist, err := h.wishlistRepo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	items, err := h.itemRepo.GetByWishlistID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]*dto.ItemResponse, len(items))
	for i, item := range items {
		responses[i] = dto.ToItemResponse(item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *WishlistItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("wishlistId")
	itemID := r.PathValue("id")

	wishlist, err := h.wishlistRepo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	item, err := h.itemRepo.GetByID(r.Context(), itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if item.WishlistID != wishlistID {
		http.Error(w, "Item not found in this wishlist", http.StatusNotFound)
		return
	}

	var req dto.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.ToModel(item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.itemRepo.Update(r.Context(), item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToItemResponse(item))
}

func (h *WishlistItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("wishlistId")
	itemID := r.PathValue("id")

	wishlist, err := h.wishlistRepo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	item, err := h.itemRepo.GetByID(r.Context(), itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if item.WishlistID != wishlistID {
		http.Error(w, "Item not found in this wishlist", http.StatusNotFound)
		return
	}

	if err := h.itemRepo.Delete(r.Context(), itemID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
