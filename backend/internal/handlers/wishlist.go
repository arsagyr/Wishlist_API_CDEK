package handlers

import (
	"encoding/json"
	"net/http"

	"test_cdek/backend/internal/dto"
	"test_cdek/backend/internal/middleware"
	"test_cdek/backend/internal/repository"
)

type WishlistHandler struct {
	repo *repository.WishlistRepository
}

func NewWishlistHandler(repo *repository.WishlistRepository) *WishlistHandler {
	return &WishlistHandler{repo: repo}
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	wishlist, err := req.ToModel(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), wishlist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToWishlistResponse(wishlist))
}

func (h *WishlistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("id")
	wishlist, err := h.repo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToWishlistResponse(wishlist))
}

func (h *WishlistHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlists, err := h.repo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]*dto.WishlistResponse, len(wishlists))
	for i, w := range wishlists {
		responses[i] = dto.ToWishlistResponse(w)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("id")
	wishlist, err := h.repo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req dto.UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := wishlist.Update(req.Title, req.Description, req.EventDate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), wishlist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToWishlistResponse(wishlist))
}

func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlistID := r.PathValue("id")
	wishlist, err := h.repo.GetByID(r.Context(), wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if wishlist.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := h.repo.Delete(r.Context(), wishlistID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
