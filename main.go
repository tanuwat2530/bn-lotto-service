package main

import (
	"fmt"
	"log"
	postgres "lotto-backend-api/db"
	"lotto-backend-api/routes"
	"net/http"
)

func main() {
	// --- 1. Initialize Database ---
	// Your InitDB function from the `db` package should handle pool configuration.
	// It's good practice for it to return both the db object and an error.
	db, err := postgres.InitDB()
	if err != nil {
		log.Fatalf("FATAL: Could not initialize database: %v", err)
	}
	fmt.Println("✅ Database connection pool initialized.")

	// --- 2. Setup Routes and Inject Dependencies ---
	// Pass the database connection pool to your router setup function.
	// This is the dependency injection step.
	router := routes.SetupRoutes(db)
	fmt.Println("✅ Routes have been configured.")

	// --- 3. Start the Server ---
	// The router returned by SetupRoutes is used as the server's handler.
	port := "8080"
	fmt.Printf("🚀 Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("FATAL: Error starting server: %v", err)
	}

}
