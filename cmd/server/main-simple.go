package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("========= Simple version starting =================")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Starting server on port %s\n", port)
	
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"Service is running"}`))
	})
	
	// Main page endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Business Exchange - Simple Version</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100">
    <div class="container mx-auto px-4 py-8">
        <h1 class="text-4xl font-bold text-center text-gray-800 mb-8">
            🚀 Business Exchange Marketplace
        </h1>
        <div class="bg-white rounded-lg shadow-lg p-6 max-w-2xl mx-auto">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Simple Version</h2>
            <p class="text-gray-600 mb-4">This is a simple version to test deployment without database dependencies.</p>
            <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
                ✅ Service is running successfully!
            </div>
            <div class="mt-4 text-sm text-gray-500">
                <p>Port: ` + port + `</p>
                <p>Environment: Cloud Run</p>
                <p>Status: Deployed via GitHub Actions</p>
            </div>
        </div>
    </div>
</body>
</html>`
		w.Write([]byte(html))
	})
	
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 20 * time.Second,
	}
	
	fmt.Printf("Server configured: %+v\n", srv)
	
	go func() {
		fmt.Printf("Server starting on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	
	fmt.Println("Server exited")
}

