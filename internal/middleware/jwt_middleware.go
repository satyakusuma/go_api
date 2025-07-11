package middleware

import (
	"context"
	"net/http"
	"strings"
	"github.com/satyakusuma/go-rest-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := config.GetEnv("JWT_SECRET", "your_jwt_secret_key")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}
		ctx := context.WithValue(r.Context(), "user_id", int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}