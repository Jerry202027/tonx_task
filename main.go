package main

import (
	"context"
	"log"
	"net/http"
	"tonx_task/database"
	"tonx_task/server"
)

func main() {
	// Create root context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Setup cache module.

	// Setup database module.
	database.InitDatabase(ctx)
	defer database.FinalizeDatabase(ctx)

	// Setup HTTP Server
	srv := server.CreateServer(ctx, ":8080")

	log.Println("Server is running on http://localhost:8080")
	// Listen on http://localhost:8080
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
