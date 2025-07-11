package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/satyakusuma/go-rest-api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{authService: services.NewAuthService(db)}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(input.Username, input.Password); err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(input.Username, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}


func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}


func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	var input struct {
		CurrentPassword string `json:"current_password"`
		NewUsername     string `json:"new_username"`
		NewPassword     string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.CurrentPassword == "" {
		http.Error(w, "Current password is required", http.StatusBadRequest)
		return
	}

	if input.NewUsername == "" && input.NewPassword == "" {
		http.Error(w, "At least one of new_username or new_password must be provided", http.StatusBadRequest)
		return
	}

	if err := h.authService.UpdateProfile(userID, input.CurrentPassword, input.NewUsername, input.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}