// @title API Monitor Backend
// @version 1.0
// @description Backend API for monitoring and analyzing API requests
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@apimonitor.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	//_ "api-monitor-backend/docs" // Import generated docs
	"api-monitor-backend/internal/database"
	"api-monitor-backend/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	database.SeedData()
	defer database.CloseDB()

	// Create router
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/requests", handlers.GetRequests).Methods("GET")
	router.HandleFunc("/api/problems", handlers.GetProblems).Methods("GET")
	router.HandleFunc("/api/proxy/{endpoint:.*}", handlers.ProxyRequest).Methods("GET", "POST", "PUT", "DELETE")
	// Swagger docs
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Vite default port
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
