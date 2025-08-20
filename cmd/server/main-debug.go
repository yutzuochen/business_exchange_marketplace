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

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("========= Debug version starting =================")
	
	// 1. 載入環境變數
	fmt.Println("1. Loading environment variables...")
	_ = godotenv.Load()
	
	// 2. 獲取配置
	fmt.Println("2. Getting configuration...")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("PORT")
	
	if appPort == "" {
		appPort = "8080"
	}
	
	fmt.Printf("DB_HOST: %s\n", dbHost)
	fmt.Printf("DB_PORT: %s\n", dbPort)
	fmt.Printf("DB_USER: %s\n", dbUser)
	fmt.Printf("DB_NAME: %s\n", dbName)
	fmt.Printf("APP_PORT: %s\n", appPort)
	
	// 3. 測試資料庫連接
	fmt.Println("3. Testing database connection...")
	
	// 當使用 Cloud SQL Proxy 時，連接到本地
	var connectionHost string
	if os.Getenv("CLOUDSQL_CONNECTION_NAME") != "" {
		// 使用 Cloud SQL Proxy，連接到本地
		connectionHost = "127.0.0.1"
		fmt.Println("Using Cloud SQL Proxy (local connection)")
	} else {
		// 直接連接
		connectionHost = dbHost
		fmt.Println("Using direct connection")
	}
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", 
		dbUser, dbPassword, connectionHost, dbPort, dbName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// 設置連接池參數
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	
	// 測試連接
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	
	fmt.Println("Database connection successful!")
	
	// 4. 創建 HTTP 服務器
	fmt.Println("4. Creating HTTP server...")
	mux := http.NewServeMux()
	
	// 健康檢查端點
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"Service is running"}`))
	})
	
	// 主頁端點
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Business Exchange - Debug Version</title>
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
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Debug Version</h2>
            <p class="text-gray-600 mb-4">This is a debug version to test database connectivity.</p>
            <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
                ✅ Database connection successful!
            </div>
            <div class="mt-4 text-sm text-gray-500">
                <p>DB Host: ` + dbHost + `</p>
                <p>DB Port: ` + dbPort + `</p>
                <p>DB Name: ` + dbName + `</p>
            </div>
        </div>
    </div>
</body>
</html>`
		w.Write([]byte(html))
	})
	
	// 資料庫測試端點
	mux.HandleFunc("/test-db", func(w http.ResponseWriter, r *http.Request) {
		var result int
		err := db.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error":"Database query failed: %v"}`, err)))
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"ok","result":%d}`, result)))
	})
	
	// 5. 啟動服務器
	fmt.Printf("5. Starting server on port %s...\n", appPort)
	srv := &http.Server{
		Addr:              ":" + appPort,
		Handler:           mux,
		ReadHeaderTimeout: 20 * time.Second,
	}
	
	go func() {
		fmt.Printf("Server starting on port %s\n", appPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// 優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	fmt.Println("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	
	fmt.Println("Server exited")
}
