package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Testing database connection...")
	
	// 連接字串
	dsn := "app:app_password@tcp(34.70.172.32:3306)/business_exchange?parseTime=true&loc=Local"
	
	// 設置連接超時
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
	
	// 測試簡單查詢
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	
	fmt.Printf("Query result: %d\n", result)
	fmt.Println("Database test completed successfully!")
}
