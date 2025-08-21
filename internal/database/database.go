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
		PrepareStmt:                              true,
		Logger:                                   logger.Default.LogMode(logMode),
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableNestedTransaction:                 true,
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

// AutoMigrate is deprecated - use golang-migrate instead
// func AutoMigrate(db *gorm.DB) error {
// 	return db.AutoMigrate(
// 		&models.User{},
// 		&models.Listing{},
// 		&models.Image{},
// 		&models.Favorite{},
// 		&models.Message{},
// 		&models.Transaction{},
// 	)
// }

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

	log.Printf("Created %d users successfully :))))", len(users))

	// Create sample listings
	listings := []models.Listing{
		{
			Title:             "MacBook Pro 2023 - 14 inch",
			Description:       "Excellent condition MacBook Pro with M2 Pro chip. 16GB RAM, 512GB SSD. Perfect for development and design work.",
			Price:             180000, // $1,800.00
			Category:          "Electronics",
			Condition:         "excellent",
			Location:          "San Francisco, CA",
			Status:            "active",
			OwnerID:           users[1].ID, // John Doe
			ViewCount:         45,
			BrandStory:        "Apple's flagship laptop, designed for professionals who demand the best performance and build quality.",
			Rent:              0, // Not for rent
			Floor:             0, // Not applicable
			Equipment:         "Includes original charger, protective case, and documentation",
			Decoration:        "modern",
			AnnualRevenue:     0,                                           // Not applicable for personal items
			GrossProfitRate:   0.0,                                         // Not applicable
			FastestMovingDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Default date
			PhoneNumber:       "+1-555-0123",
			SquareMeters:      0.0, // Not applicable
			Industry:          "Technology",
			Deposit:           0, // Not applicable
		},
		{
			Title:             "Vintage Leather Office Chair",
			Description:       "Beautiful vintage leather office chair from the 1960s. High-quality leather, very comfortable. Great for home office.",
			Price:             45000, // $450.00
			Category:          "Furniture",
			Condition:         "good",
			Location:          "New York, NY",
			Status:            "active",
			OwnerID:           users[2].ID, // Jane Smith
			ViewCount:         23,
			BrandStory:        "Authentic mid-century modern design, crafted by skilled artisans in the golden age of American furniture making.",
			Rent:              0, // Not for rent
			Floor:             0, // Not applicable
			Equipment:         "Original leather upholstery, chrome base, swivel mechanism",
			Decoration:        "vintage",
			AnnualRevenue:     0,                                           // Not applicable for personal items
			GrossProfitRate:   0.0,                                         // Not applicable
			FastestMovingDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Default date
			PhoneNumber:       "+1-555-0124",
			SquareMeters:      0.0, // Not applicable
			Industry:          "Furniture",
			Deposit:           0, // Not applicable
		},
		{
			Title:             "Professional Camera Lens Set",
			Description:       "Complete set of professional camera lenses: 24-70mm f/2.8, 70-200mm f/2.8, and 50mm f/1.4. Perfect for photography.",
			Price:             320000, // $3,200.00
			Category:          "Electronics",
			Condition:         "excellent",
			Location:          "Los Angeles, CA",
			Status:            "active",
			OwnerID:           users[3].ID, // Bob Wilson
			ViewCount:         67,
			BrandStory:        "Canon L-series professional lenses, used by award-winning photographers worldwide for their exceptional optical quality.",
			Rent:              0, // Not for rent
			Floor:             0, // Not applicable
			Equipment:         "Includes lens caps, hoods, carrying case, and cleaning kit",
			Decoration:        "professional",
			AnnualRevenue:     0,                                           // Not applicable for personal items
			GrossProfitRate:   0.0,                                         // Not applicable
			FastestMovingDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Default date
			PhoneNumber:       "+1-555-0125",
			SquareMeters:      0.0, // Not applicable
			Industry:          "Photography",
			Deposit:           0, // Not applicable
		},
		{
			Title:             "Antique Wooden Dining Table",
			Description:       "Stunning antique wooden dining table with 6 chairs. Solid oak construction, beautiful craftsmanship. Perfect for family gatherings.",
			Price:             120000, // $1,200.00
			Category:          "Furniture",
			Condition:         "good",
			Location:          "Chicago, IL",
			Status:            "active",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         34,
			BrandStory:        "Handcrafted in the early 1900s by master woodworkers, this table has been the centerpiece of family celebrations for generations.",
			Rent:              0, // Not for rent
			Floor:             0, // Not applicable
			Equipment:         "Table with 6 matching chairs, table runner, and care instructions",
			Decoration:        "antique",
			AnnualRevenue:     0,                                           // Not applicable for personal items
			GrossProfitRate:   0.0,                                         // Not applicable
			FastestMovingDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Default date
			PhoneNumber:       "+1-555-0126",
			SquareMeters:      0.0, // Not applicable
			Industry:          "Furniture",
			Deposit:           0, // Not applicable
		},
		{
			Title:             "Mountain Bike - Trek Fuel EX 8",
			Description:       "High-end mountain bike in great condition. Carbon frame, full suspension, perfect for trail riding. Includes helmet and accessories.",
			Price:             280000, // $2,800.00
			Category:          "Sports",
			Condition:         "excellent",
			Location:          "Denver, CO",
			Status:            "active",
			OwnerID:           users[1].ID, // John Doe
			ViewCount:         89,
			BrandStory:        "Trek's premium mountain bike series, designed for serious riders who demand performance and reliability on challenging trails.",
			Rent:              0, // Not for rent
			Floor:             0, // Not applicable
			Equipment:         "Bike, helmet, pump, repair kit, and maintenance guide",
			Decoration:        "sporty",
			AnnualRevenue:     0,                                           // Not applicable for personal items
			GrossProfitRate:   0.0,                                         // Not applicable
			FastestMovingDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Default date
			PhoneNumber:       "+1-555-0127",
			SquareMeters:      0.0, // Not applicable
			Industry:          "Sports",
			Deposit:           0, // Not applicable
		},
		{
			Title:             "Downtown Coffee Shop for Sale",
			Description:       "Profitable coffee shop in prime downtown location. Established customer base, modern equipment, great potential for growth.",
			Price:             850000, // $8,500.00
			Category:          "Business",
			Condition:         "excellent",
			Location:          "Seattle, WA",
			Status:            "active",
			OwnerID:           users[2].ID, // Jane Smith
			ViewCount:         156,
			BrandStory:        "A beloved local coffee shop that has been serving the community for over 5 years, known for quality coffee and warm atmosphere.",
			Rent:              8500, // Monthly rent
			Floor:             1,    // Ground floor
			Equipment:         "Commercial espresso machines, grinders, refrigerators, furniture, POS system",
			Decoration:        "modern",
			AnnualRevenue:     450000, // $450,000 annual revenue
			GrossProfitRate:   0.35,   // 35% gross profit margin
			FastestMovingDate: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "+1-555-0128",
			SquareMeters:      120.0, // 120 square meters
			Industry:          "Food & Beverage",
			Deposit:           50000, // $50,000 deposit
		},
		{
			Title:             "Tech Startup Office Space",
			Description:       "Modern office space perfect for tech startups. Open floor plan, meeting rooms, high-speed internet, parking included.",
			Price:             0, // For rent only
			Category:          "Real Estate",
			Condition:         "excellent",
			Location:          "Austin, TX",
			Status:            "active",
			OwnerID:           users[3].ID, // Bob Wilson
			ViewCount:         78,
			BrandStory:        "A state-of-the-art office building designed specifically for modern tech companies, with amenities that foster innovation and collaboration.",
			Rent:              12000, // Monthly rent
			Floor:             3,     // Third floor
			Equipment:         "Furnished workstations, conference rooms, kitchen, lounge area, gym",
			Decoration:        "contemporary",
			AnnualRevenue:     0,   // Not applicable for rental properties
			GrossProfitRate:   0.0, // Not applicable
			FastestMovingDate: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "+1-555-0129",
			SquareMeters:      300.0, // 300 square meters
			Industry:          "Real Estate",
			Deposit:           24000, // 2 months deposit
		},
		{
			Title:             "E-commerce Business - Fashion Brand",
			Description:       "Established online fashion brand with strong social media presence. Includes website, inventory, and customer database.",
			Price:             650000, // $650,000
			Category:          "Business",
			Condition:         "good",
			Location:          "Miami, FL",
			Status:            "active",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         203,
			BrandStory:        "A trendy fashion brand that has successfully built a loyal following through social media marketing and quality products.",
			Rent:              0, // Online business
			Floor:             0, // Not applicable
			Equipment:         "Website, inventory management system, social media accounts, customer database",
			Decoration:        "trendy",
			AnnualRevenue:     320000, // $320,000 annual revenue
			GrossProfitRate:   0.42,   // 42% gross profit margin
			FastestMovingDate: time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "+1-555-0130",
			SquareMeters:      0.0, // Online business
			Industry:          "Fashion & Retail",
			Deposit:           0, // Not applicable
		},
		{
			Title:             "Manufacturing Equipment - CNC Machines",
			Description:       "Complete set of CNC machining equipment for manufacturing business. Includes 3 CNC mills, 2 lathes, and support equipment.",
			Price:             420000, // $420,000
			Category:          "Industrial",
			Condition:         "excellent",
			Location:          "Detroit, MI",
			Status:            "active",
			OwnerID:           users[1].ID, // John Doe
			ViewCount:         45,
			BrandStory:        "Professional-grade CNC equipment from leading manufacturers, maintained to the highest standards for precision manufacturing.",
			Rent:              0, // Not for rent
			Floor:             0, // Not applicable
			Equipment:         "3 CNC mills, 2 CNC lathes, tooling, measuring equipment, safety gear",
			Decoration:        "industrial",
			AnnualRevenue:     0,   // Not applicable for equipment sales
			GrossProfitRate:   0.0, // Not applicable
			FastestMovingDate: time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "+1-555-0131",
			SquareMeters:      0.0, // Not applicable
			Industry:          "Manufacturing",
			Deposit:           0, // Not applicable
		},
		{
			Title:             "Restaurant Kitchen Equipment",
			Description:       "Complete commercial kitchen setup including ovens, grills, refrigerators, and prep stations. Perfect for new restaurant.",
			Price:             180000, // $180,000
			Category:          "Restaurant",
			Condition:         "good",
			Location:          "Portland, OR",
			Status:            "active",
			OwnerID:           users[2].ID, // Jane Smith
			ViewCount:         67,
			BrandStory:        "Professional kitchen equipment from a successful restaurant that has been serving quality food for over 10 years.",
			Rent:              50000, // Not for rent
			Floor:             66,    // Not applicable
			Equipment:         "Commercial ovens, grills, fryers, refrigerators, prep tables, dishwashers",
			Decoration:        "commercial",
			AnnualRevenue:     600, // Not applicable for equipment sales
			GrossProfitRate:   3.0, // Not applicable
			FastestMovingDate: time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "+1-555-0132",
			SquareMeters:      34.5, // Not applicable
			Industry:          "Food Service",
			Deposit:           40000, // Not applicable
		},
		{
			Title:             "Franchise Opportunity - Fast Food",
			Description:       "Established fast food franchise with proven business model. Includes training, marketing support, and operational guidelines.",
			Price:             250000, // $250,000 franchise fee
			Category:          "Franchise",
			Condition:         "excellent",
			Location:          "Phoenix, AZ",
			Status:            "active",
			OwnerID:           users[3].ID, // Bob Wilson
			ViewCount:         134,
			BrandStory:        "A nationally recognized fast food brand with over 500 locations, offering entrepreneurs a proven path to business success.",
			Rent:              0, // Franchise opportunity
			Floor:             0, // Not applicable
			Equipment:         "Franchise license, training materials, marketing support, operational manual",
			Decoration:        "branded",
			AnnualRevenue:     180000, // $180,000 average annual revenue
			GrossProfitRate:   0.28,   // 28% gross profit margin
			FastestMovingDate: time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "+1-555-0133",
			SquareMeters:      0.0, // Not applicable
			Industry:          "Food Service",
			Deposit:           50000, // $50,000 initial deposit
		},
	}
	log.Printf("============= start to create listings =============")
	for i := range listings {
		log.Printf("listings[i]: %+v\n", listings[i])
		if err := db.Create(&listings[i]).Error; err != nil {
			log.Printf("Failed QQ to create listing %s: %v", listings[i].Title, err)
			// fmt.Printf("listings[i]: %+v\n", listings[i])

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
		// Coffee Shop images
		{
			ListingID: listings[5].ID,
			Filename:  "coffee_shop_1.jpg",
			URL:       "/static/images/coffee_shop_1.jpg",
			AltText:   "Downtown coffee shop interior",
			Order:     0,
			IsPrimary: true,
		},
		// Office Space images
		{
			ListingID: listings[6].ID,
			Filename:  "office_space_1.jpg",
			URL:       "/static/images/office_space_1.jpg",
			AltText:   "Modern office space interior",
			Order:     0,
			IsPrimary: true,
		},
		// Fashion Brand images
		{
			ListingID: listings[7].ID,
			Filename:  "fashion_brand_1.jpg",
			URL:       "/static/images/fashion_brand_1.jpg",
			AltText:   "Fashion brand website screenshot",
			Order:     0,
			IsPrimary: true,
		},
		// CNC Equipment images
		{
			ListingID: listings[8].ID,
			Filename:  "cnc_equipment_1.jpg",
			URL:       "/static/images/cnc_equipment_1.jpg",
			AltText:   "CNC machining equipment",
			Order:     0,
			IsPrimary: true,
		},
		// Kitchen Equipment images
		{
			ListingID: listings[9].ID,
			Filename:  "kitchen_equipment_1.jpg",
			URL:       "/static/images/kitchen_equipment_1.jpg",
			AltText:   "Commercial kitchen equipment",
			Order:     0,
			IsPrimary: true,
		},
		// Franchise images
		{
			ListingID: listings[10].ID,
			Filename:  "franchise_1.jpg",
			URL:       "/static/images/franchise_1.jpg",
			AltText:   "Fast food franchise storefront",
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
		{
			UserID:    users[2].ID,    // Jane Smith
			ListingID: listings[5].ID, // Coffee Shop
		},
		{
			UserID:    users[4].ID,    // Alice Johnson
			ListingID: listings[6].ID, // Office Space
		},
		{
			UserID:    users[1].ID,    // John Doe
			ListingID: listings[7].ID, // Fashion Brand
		},
		{
			UserID:    users[3].ID,    // Bob Wilson
			ListingID: listings[8].ID, // CNC Equipment
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
		{
			SenderID:   users[4].ID,     // Alice Johnson
			ReceiverID: users[2].ID,     // Jane Smith
			ListingID:  &listings[5].ID, // Coffee Shop
			Subject:    "Coffee Shop Investment",
			Content:    "Hi Jane, I'm very interested in your coffee shop! Can we schedule a meeting to discuss the business details?",
			IsRead:     false,
		},
		{
			SenderID:   users[1].ID,     // John Doe
			ReceiverID: users[3].ID,     // Bob Wilson
			ListingID:  &listings[6].ID, // Office Space
			Subject:    "Office Space Rental",
			Content:    "Hi Bob, your office space looks perfect for our startup! Is it still available for rent?",
			IsRead:     false,
		},
		{
			SenderID:   users[2].ID,     // Jane Smith
			ReceiverID: users[4].ID,     // Alice Johnson
			ListingID:  &listings[7].ID, // Fashion Brand
			Subject:    "Fashion Brand Partnership",
			Content:    "Hi Alice, I love your fashion brand concept! Would you be interested in a potential partnership?",
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
