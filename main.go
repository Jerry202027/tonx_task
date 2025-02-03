package main

import (
	"context"
	"log"
	"net/http"
	"tonx_task/cache"
	"tonx_task/database"
	"tonx_task/server"
	"tonx_task/service"
)

func main() {
	// Create root context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup database module.
	database.InitDatabase(ctx)
	defer database.FinalizeDatabase(ctx)

	// init service
	service.ServiceInit()

	// Setup cache module
	cache.InitRedis("localhost:6379", "", 0)
	log.Println("Redis initialized successfully")

	// only execute once
	// // create the Flight and Booking tables
	// err := service.FlightBookingService.AutoMigrate()
	// if err != nil {
	// 	log.Fatalf("AutoMigrate error: %v\n", err)
	// }

	// // add test case to database
	// err = testcase.InsertSampleFlights()
	// if err != nil {
	// 	log.Fatalf("insert sample flights error: %v\n", err)
	// }

	// Setup HTTP Server
	srv := server.CreateServer(ctx, ":8080")

	log.Println("Server is running on http://localhost:8080")
	// Listen on http://localhost:8080
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
