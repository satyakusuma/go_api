package routes

import (
	"database/sql"
	"github.com/satyakusuma/go-rest-api/internal/handlers"
	"github.com/satyakusuma/go-rest-api/internal/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(router *mux.Router, db *sql.DB) {
	authHandler := handlers.NewAuthHandler(db)

	// Public routes
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api/auth").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", authHandler.UpdateProfile).Methods("POST")

	// Swagger UI route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}