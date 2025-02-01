package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tonx_task/api"
)

func CreateServer(ctx context.Context, address string) *http.Server {
	// Setup HTTP Server.
	server := &http.Server{
		Addr:    address,
		Handler: api.GetRouter(),
	}

	// Install the shutdown handler.
	installShutdownHandler(ctx, server)

	return server
}

func installShutdownHandler(ctx context.Context, server *http.Server) {
	// Create signal channel & shutdown timeout context.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Catch signals in a separate goroutine.
	go func() {
		// Wait for signals.
		sig := <-sigChan
		signal.Stop(sigChan)
		log.Printf("Received signal: %s. Initiating graceful shutdown...\n", sig)

		timeoutCtx, cancel := context.WithTimeout(ctx,
			10*time.Second)
		defer cancel()

		// Perform graceful shutdown.
		log.Println(ctx, "Initiating graceful shutdown...")

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v\n", err)
		}

		log.Println("Server gracefully shut down.")
	}()
}
