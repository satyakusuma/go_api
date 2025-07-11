package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/satyakusuma/go-rest-api/internal/config"
	"github.com/satyakusuma/go-rest-api/internal/database"
	"github.com/satyakusuma/go-rest-api/internal/routes"
	"github.com/gorilla/mux"
	_ "github.com/swaggo/http-swagger" // Import http-swagger
	_ "github.com/satyakusuma/go-rest-api/docs" // Import generated Swagger docs
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize router
	router := mux.NewRouter()

	// Log all incoming requests
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received %s request for %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Setup routes
	routes.SetupRoutes(router, db)

	// Start server
	port := config.GetEnv("PORT", "8080")
	fmt.Printf("Server running on port %s\n", port)
	fmt.Println("Swagger UI available at http://localhost:" + port + "/swagger/index.html")
	log.Fatal(http.ListenAndServe(":"+port, router))
}