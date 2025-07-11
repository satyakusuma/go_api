package services

import (
	"database/sql"
	"errors"
	"time"
	"github.com/satyakusuma/go-rest-api/internal/config"
	"github.com/satyakusuma/go-rest-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	return err
}

func (s *AuthService) Login(username, password string) (string, error) {
	var user models.User
	var hashedPassword string

	err := s.db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetProfile(userID int64) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, username, created_at FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateProfile(userID int64, currentPassword, newUsername, newPassword string) error {
	// Verify current password
	var hashedPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Prepare update query
	query := "UPDATE users SET"
	args := []interface{}{}
	if newUsername != "" {
		query += " username = ?,"
		args = append(args, newUsername)
	}
	if newPassword != "" {
		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		query += " password = ?,"
		args = append(args, hashedNewPassword)
	}

	// Remove trailing comma and add WHERE clause
	if len(args) == 0 {
		return errors.New("no fields to update")
	}
	query = query[:len(query)-1] + " WHERE id = ?"
	args = append(args, userID)

	// Execute update
	_, err = s.db.Exec(query, args...)
	return err
}

func generateJWT(userID int64) (string, error) {
	secret := config.GetEnv("JWT_SECRET", "your_jwt_secret_key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(secret))
}