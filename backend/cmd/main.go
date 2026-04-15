// cmd/main.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"test_cdek/backend/internal/handlers"
	"test_cdek/backend/internal/middleware"
	"test_cdek/backend/internal/repository"
	"test_cdek/backend/pkg/jwt"
)

var db *sql.DB

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/mydb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
	}
	jwtManager := jwt.NewManager(jwtSecret)

	userRepo := repository.NewUserRepository(db)
	wishlistRepo := repository.NewWishlistRepository(db)
	itemRepo := repository.NewWishlistItemRepository(db)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	userHandler := handlers.NewUserHandler(userRepo, jwtManager)
	wishlistHandler := handlers.NewWishlistHandler(wishlistRepo)
	itemHandler := handlers.NewWishlistItemHandler(wishlistRepo, itemRepo)
	publicHandler := handlers.NewPublicHandler(wishlistRepo, itemRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", userHandler.Register)
	mux.HandleFunc("POST /api/auth/login", userHandler.Login)

	mux.HandleFunc("GET /api/wishlists", authMiddleware.AuthHandler(wishlistHandler.ListByUser))
	mux.HandleFunc("POST /api/wishlists", authMiddleware.AuthHandler(wishlistHandler.Create))
	mux.HandleFunc("GET /api/wishlists/{id}", authMiddleware.AuthHandler(wishlistHandler.GetByID))
	mux.HandleFunc("PUT /api/wishlists/{id}", authMiddleware.AuthHandler(wishlistHandler.Update))
	mux.HandleFunc("DELETE /api/wishlists/{id}", authMiddleware.AuthHandler(wishlistHandler.Delete))

	mux.HandleFunc("GET /api/wishlists/{wishlistId}/items", authMiddleware.AuthHandler(itemHandler.ListByWishlist))
	mux.HandleFunc("POST /api/wishlists/{wishlistId}/items", authMiddleware.AuthHandler(itemHandler.Create))
	mux.HandleFunc("GET /api/wishlists/{wishlistId}/items/{id}", authMiddleware.AuthHandler(itemHandler.GetByID))
	mux.HandleFunc("PUT /api/wishlists/{wishlistId}/items/{id}", authMiddleware.AuthHandler(itemHandler.Update))
	mux.HandleFunc("DELETE /api/wishlists/{wishlistId}/items/{id}", authMiddleware.AuthHandler(itemHandler.Delete))

	mux.HandleFunc("GET /api/public/wishlists/{token}", publicHandler.GetWishlistByToken)
	mux.HandleFunc("POST /api/public/wishlists/{token}/reserve", publicHandler.ReserveItem)

	mux.HandleFunc("GET /health", healthHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.PingContext(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
