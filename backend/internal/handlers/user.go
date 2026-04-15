package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"test_cdek/backend/internal/dto"
	"test_cdek/backend/internal/repository"
	"test_cdek/backend/pkg/jwt"
)

type UserHandler struct {
	repo       *repository.UserRepository
	jwtManager *jwt.Manager
}

func NewUserHandler(repo *repository.UserRepository, jwtManager *jwt.Manager) *UserHandler {
	return &UserHandler{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := req.ToModel()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()

	if err := h.repo.Create(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToUserResponse(user))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.ValidatePassword(req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.IsActiveUser() {
		http.Error(w, "Account is deactivated", http.StatusForbidden)
		return
	}

	user.RecordLogin()
	if err := h.repo.Update(r.Context(), user); err != nil {
		log.Printf("Failed to update last login: %v", err)
	}

	token, err := h.jwtManager.Generate(user.ID, user.Email, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := dto.LoginResponse{
		User:      *dto.ToUserResponse(user),
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
