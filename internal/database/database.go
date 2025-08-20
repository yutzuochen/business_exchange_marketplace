package database

import (
	"log"
	"time"

	"trade_company/internal/config"
	"trade_company/internal/models"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config, _ any) (*gorm.DB, error) {
	logMode := logger.Info
	if cfg.AppEnv == "production" {
		logMode = logger.Warn
	}
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN()), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Listing{},
		&models.Image{},
		&models.Favorite{},
		&models.Message{},
		&models.Transaction{},
	)
}

// SeedData adds sample data to the database for testing
func SeedData(db *gorm.DB) error {
	log.Println("Seeding database with sample data...")

	// Check if users already exist
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)

	log.Printf("Current user count: %d", userCount)

	// For testing purposes, always seed data
	// if userCount > 0 {
	// 	log.Println("Database already has data, skipping seed")
	// 	return nil
	// }

	// Clear existing data first
	log.Println("Clearing existing data...")
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM messages")
	db.Exec("DELETE FROM favorites")
	db.Exec("DELETE FROM images")
	db.Exec("DELETE FROM listings")
	db.Exec("DELETE FROM users")

	log.Println("Existing data cleared, starting to seed...")

	// Create sample users
	users := []models.User{
		{
			Email:        "admin@example.com",
			Username:     "admin",
			PasswordHash: hashPassword("admin123"),
			FirstName:    "Admin",
			LastName:     "User",
			Role:         "admin",
			IsActive:     true,
		},
		{
			Email:        "john.doe@example.com",
			Username:     "johndoe",
			PasswordHash: hashPassword("password123"),
			FirstName:    "John",
			LastName:     "Doe",
			Role:         "user",
			IsActive:     true,
		},
		{
			Email:        "jane.smith@example.com",
			Username:     "janesmith",
			PasswordHash: hashPassword("password123"),
			FirstName:    "Jane",
			LastName:     "Smith",
			Role:         "user",
			IsActive:     true,
		},
		{
			Email:        "bob.wilson@example.com",
			Username:     "bobwilson",
			PasswordHash: hashPassword("password123"),
			FirstName:    "Bob",
			LastName:     "Wilson",
			Role:         "user",
			IsActive:     true,
		},
		{
			Email:        "alice.johnson@example.com",
			Username:     "alicejohnson",
			PasswordHash: hashPassword("password123"),
			FirstName:    "Alice",
			LastName:     "Johnson",
			Role:         "user",
			IsActive:     true,
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Printf("Failed to create user %s: %v", users[i].Username, err)
			return err
		}
		log.Printf("Created user: %s (%s)", users[i].Username, users[i].Email)
	}

	log.Printf("Created %d users successfully", len(users))

	// Create sample listings
	listings := []models.Listing{
		{
			Title:       "MacBook Pro 2023 - 14 inch",
			Description: "Excellent condition MacBook Pro with M2 Pro chip. 16GB RAM, 512GB SSD. Perfect for development and design work.",
			Price:       180000, // $1,800.00
			Category:    "Electronics",
			Condition:   "excellent",
			Location:    "San Francisco, CA",
			Status:      "active",
			OwnerID:     users[1].ID, // John Doe
			ViewCount:   45,
		},
		{
			Title:       "Vintage Leather Office Chair",
			Description: "Beautiful vintage leather office chair from the 1960s. High-quality leather, very comfortable. Great for home office.",
			Price:       45000, // $450.00
			Category:    "Furniture",
			Condition:   "good",
			Location:    "New York, NY",
			Status:      "active",
			OwnerID:     users[2].ID, // Jane Smith
			ViewCount:   23,
		},
		{
			Title:       "Professional Camera Lens Set",
			Description: "Complete set of professional camera lenses: 24-70mm f/2.8, 70-200mm f/2.8, and 50mm f/1.4. Perfect for photography.",
			Price:       320000, // $3,200.00
			Category:    "Electronics",
			Condition:   "excellent",
			Location:    "Los Angeles, CA",
			Status:      "active",
			OwnerID:     users[3].ID, // Bob Wilson
			ViewCount:   67,
		},
		{
			Title:       "Antique Wooden Dining Table",
			Description: "Stunning antique wooden dining table with 6 chairs. Solid oak construction, beautiful craftsmanship. Perfect for family gatherings.",
			Price:       120000, // $1,200.00
			Category:    "Furniture",
			Condition:   "good",
			Location:    "Chicago, IL",
			Status:      "active",
			OwnerID:     users[4].ID, // Alice Johnson
			ViewCount:   34,
		},
		{
			Title:       "Mountain Bike - Trek Fuel EX 8",
			Description: "High-end mountain bike in great condition. Carbon frame, full suspension, perfect for trail riding. Includes helmet and accessories.",
			Price:       280000, // $2,800.00
			Category:    "Sports",
			Condition:   "excellent",
			Location:    "Denver, CO",
			Status:      "active",
			OwnerID:     users[1].ID, // John Doe
			ViewCount:   89,
		},
	}

	for i := range listings {
		if err := db.Create(&listings[i]).Error; err != nil {
			log.Printf("Failed to create listing %s: %v", listings[i].Title, err)
			return err
		}
		log.Printf("Created listing: %s ($%.2f)", listings[i].Title, float64(listings[i].Price)/100)
	}

	log.Printf("Created %d listings successfully", len(listings))

	// Create sample images for listings
	images := []models.Image{
		// MacBook Pro images
		{
			ListingID: listings[0].ID,
			Filename:  "macbook_pro_1.jpg",
			URL:       "/static/images/macbook_pro_1.jpg",
			AltText:   "MacBook Pro front view",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[0].ID,
			Filename:  "macbook_pro_2.jpg",
			URL:       "/static/images/macbook_pro_2.jpg",
			AltText:   "MacBook Pro side view",
			Order:     1,
			IsPrimary: false,
		},
		// Office Chair images
		{
			ListingID: listings[1].ID,
			Filename:  "office_chair_1.jpg",
			URL:       "/static/images/office_chair_1.jpg",
			AltText:   "Vintage leather office chair",
			Order:     0,
			IsPrimary: true,
		},
		// Camera Lens images
		{
			ListingID: listings[2].ID,
			Filename:  "camera_lens_1.jpg",
			URL:       "/static/images/camera_lens_1.jpg",
			AltText:   "Professional camera lens set",
			Order:     0,
			IsPrimary: true,
		},
		// Dining Table images
		{
			ListingID: listings[3].ID,
			Filename:  "dining_table_1.jpg",
			URL:       "/static/images/dining_table_1.jpg",
			AltText:   "Antique wooden dining table",
			Order:     0,
			IsPrimary: true,
		},
		// Mountain Bike images
		{
			ListingID: listings[4].ID,
			Filename:  "mountain_bike_1.jpg",
			URL:       "/static/images/mountain_bike_1.jpg",
			AltText:   "Trek Fuel EX 8 mountain bike",
			Order:     0,
			IsPrimary: true,
		},
	}

	for i := range images {
		if err := db.Create(&images[i]).Error; err != nil {
			log.Printf("Failed to create image %s: %v", images[i].Filename, err)
			return err
		}
		log.Printf("Created image: %s", images[i].Filename)
	}

	log.Printf("Created %d images successfully", len(images))

	// Create sample favorites
	favorites := []models.Favorite{
		{
			UserID:    users[2].ID,    // Jane Smith
			ListingID: listings[0].ID, // MacBook Pro
		},
		{
			UserID:    users[3].ID,    // Bob Wilson
			ListingID: listings[1].ID, // Office Chair
		},
		{
			UserID:    users[4].ID,    // Alice Johnson
			ListingID: listings[2].ID, // Camera Lens
		},
		{
			UserID:    users[1].ID,    // John Doe
			ListingID: listings[3].ID, // Dining Table
		},
	}

	for i := range favorites {
		if err := db.Create(&favorites[i]).Error; err != nil {
			log.Printf("Failed to create favorite: %v", err)
			return err
		}
		log.Printf("Created favorite for user %d, listing %d", favorites[i].UserID, favorites[i].ListingID)
	}

	log.Printf("Created %d favorites successfully", len(favorites))

	// Create sample messages
	messages := []models.Message{
		{
			SenderID:   users[2].ID,     // Jane Smith
			ReceiverID: users[1].ID,     // John Doe
			ListingID:  &listings[0].ID, // MacBook Pro
			Subject:    "Question about MacBook Pro",
			Content:    "Hi John, I'm interested in your MacBook Pro. Is it still available? Can you tell me more about its condition?",
			IsRead:     false,
		},
		{
			SenderID:   users[1].ID,     // John Doe
			ReceiverID: users[2].ID,     // Jane Smith
			ListingID:  &listings[0].ID, // MacBook Pro
			Subject:    "Re: Question about MacBook Pro",
			Content:    "Hi Jane, yes it's still available! It's in excellent condition, barely used. I can send you more photos if you'd like.",
			IsRead:     true,
		},
		{
			SenderID:   users[3].ID,     // Bob Wilson
			ReceiverID: users[2].ID,     // Jane Smith
			ListingID:  &listings[1].ID, // Office Chair
			Subject:    "Office Chair Inquiry",
			Content:    "Hi Jane, I love your vintage office chair! Would you be willing to ship it to LA? I can cover shipping costs.",
			IsRead:     false,
		},
	}

	for i := range messages {
		if err := db.Create(&messages[i]).Error; err != nil {
			log.Printf("Failed to create message: %v", err)
			return err
		}
		log.Printf("Created message from user %d to user %d", messages[i].SenderID, messages[i].ReceiverID)
	}

	log.Printf("Created %d messages successfully", len(messages))

	// Create sample transactions
	transactions := []models.Transaction{
		{
			ListingID:     listings[4].ID, // Mountain Bike
			BuyerID:       users[3].ID,    // Bob Wilson
			SellerID:      users[1].ID,    // John Doe
			Amount:        280000,         // $2,800.00
			Status:        "completed",
			PaymentMethod: "PayPal",
			CompletedAt:   &[]time.Time{time.Now().Add(-24 * time.Hour)}[0], // 1 day ago
		},
		{
			ListingID:     listings[2].ID, // Camera Lens
			BuyerID:       users[4].ID,    // Alice Johnson
			SellerID:      users[3].ID,    // Bob Wilson
			Amount:        320000,         // $3,200.00
			Status:        "pending",
			PaymentMethod: "Credit Card",
		},
	}

	for i := range transactions {
		if err := db.Create(&transactions[i]).Error; err != nil {
			log.Printf("Failed to create transaction: %v", err)
			return err
		}
		log.Printf("Created transaction: $%.2f for listing %d", float64(transactions[i].Amount)/100, transactions[i].ListingID)
	}

	log.Printf("Created %d transactions successfully", len(transactions))

	log.Println("Database seeding completed successfully!")
	return nil
}

// hashPassword creates a bcrypt hash of the password
func hashPassword(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return ""
	}
	return string(hashedBytes)
}
