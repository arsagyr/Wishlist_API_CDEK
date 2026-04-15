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
    "github.com/google/uuid"
    
    "test_cdek/internal/api/dto"
    "test_cdek/internal/domain/model"
    "test_cdek/internal/repository"
)

func main() {
    // Подключение к БД
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://postgres:password@localhost:5432/mydb?sslmode=disable"
    }
    
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Настройка пула соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Проверка соединения
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    // Запуск миграций
    if err := runMigrations(db); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }
    
    // Инициализация репозитория
    userRepo := repository.NewUserRepository(db)
    
    // Настройка HTTP сервера
    mux := http.NewServeMux()
    
    // Эндпоинты API
    mux.HandleFunc("POST /api/auth/register", registerHandler(userRepo))
    mux.HandleFunc("POST /api/auth/login", loginHandler(userRepo))
    mux.HandleFunc("GET /api/users/{id}", getUserHandler(userRepo))
    mux.HandleFunc("PUT /api/users/{id}/email", updateEmailHandler(userRepo))
    mux.HandleFunc("PUT /api/users/{id}/password", updatePasswordHandler(userRepo))
    mux.HandleFunc("DELETE /api/users/{id}", deleteUserHandler(userRepo))
    mux.HandleFunc("GET /api/users", listUsersHandler(userRepo))
    
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // Graceful shutdown
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

// runMigrations - запуск миграций
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

// Обработчики HTTP

func registerHandler(repo *repository.UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req dto.RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        user, err := req.ToEntity()
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        // Генерация UUID если его нет
        if user.ID == "" {
            user.ID = uuid.New().String()
        }
        
        if err := repo.Create(r.Context(), user); err != nil {
            http.Error(w, err.Error(), http.StatusConflict)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(dto.ToUserResponse(user))
    }
}

func loginHandler(repo *repository.UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req dto.LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        user, err := repo.GetByEmail(r.Context(), req.Email)
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
        
        // Обновляем время последнего входа
        user.RecordLogin()
        if err := repo.Update(r.Context(), user); err != nil {
            log.Printf("Failed to update last login: %v", err)
        }
        
        // В реальном приложении здесь нужно генерировать JWT токен
        response := dto.LoginResponse{
            User:      *dto.ToUserResponse(user),
            Token:     "sample-jwt-token", // Заглушка
            ExpiresAt: time.Now().Add(24 * time.Hour),
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

func getUserHandler(repo *repository.UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        
        user, err := repo.GetByID(r.Context(), id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(dto.ToUserResponse(user))
    }
}

func updateEmailHandler(repo *repository.UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        
        var req dto.UpdateEmailRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        user, err := repo.GetByID(r.Context(), id)
        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        
        // Проверяем текущий пароль
        if !user.ValidatePassword(req.Password) {
            http.Error(w, "Invalid password", http.StatusUnauthorized)
            return
        }
        
        // Обновляем email
        if err := user.UpdateEmail(req.NewEmail); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        if err := repo.Update(r.Context(), user); err != nil {
            http.Error(w, "Failed to update email", http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(dto.ToUserResponse(user))
    }
}

func updatePasswordHandler(repo *repository.UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        
        var req dto.UpdatePasswordRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        user, err := repo.GetByID(r.Context(), id)
        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        
        // Проверяем текущий пароль
        if !user.ValidatePassword(req.CurrentPassword) {
            http.Error(w, "Invalid current password", http.StatusUnauthorized)
            return
        }
        
        // Обновляем пароль
        if err := user.UpdatePassword(req.NewPassword); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        if err := repo.Update(r.Context(), user); err != nil {
            http.Error(w, "Failed to update password", http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusNoContent)
    }
}

func deleteUserHandler(repo *repository.UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        
        if err := repo.Delete(r.Context(), id); err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        
        w.WriteHeader(http.StatusNoContent)
    }
}

func listUsersHandler(repo *repository.UserRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        limit := 10
        offset := 0
        
        // Можно добавить парсинг query параметров
        // limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
        // offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
        
        users, err := repo.List(r.Context(), limit, offset)
        if err != nil {
            http.Error(w, "Failed to list users", http.StatusInternalServerError)
            return
        }
        
        responses := make([]*dto.UserResponse, len(users))
        for i, user := range users {
            responses[i] = dto.ToUserResponse(user)
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(responses)
    }

	// health check
	mux.HandleFunc("GET /health", healthHandler)

	func healthHandler(w http.ResponseWriter, r *http.Request) {
		// Проверка подключения к БД
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
}