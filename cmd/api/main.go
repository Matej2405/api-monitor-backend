package main

import (
	"api-monitor-backend/internal/database"
	"api-monitor-backend/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Vite default port
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
