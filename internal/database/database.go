package database

import (
	"log"
	"strings"
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
	dsn := cfg.MySQLDSN()
	// Log DSN for debugging (without password)
	debugDSN := strings.Replace(dsn, cfg.DBPassword, "***", 1)
	log.Printf("Connecting to database with DSN: %s", debugDSN)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Shorter lifetime for Cloud SQL
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Close idle connections sooner
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
			Title:             "好快樂咖啡 ：）",
			Description:       "這裡是造夢的咖啡鄉，我們的咖啡有特殊秘方，只要一杯，你可感受全身輕飄飄，忘卻世俗一切煩惱，在夢裡，什麼都有",
			Price:             850000,
			Category:          "直營",
			Condition:         "狀況良好，9成新",
			Location:          "台中市西屯區臺灣大道三段99號",
			Status:            "活躍",
			OwnerID:           users[1].ID, // Jane Smith
			ViewCount:         156,
			BrandStory:        "我們曾經是製造業，後來改製造夢想了，我們想造福更多人！！！",
			Rent:              8500,
			Floor:             1,
			Equipment:         "手沖杯，3磅藍山咖啡，一些椅子",
			Decoration:        "夢境風",
			AnnualRevenue:     450000,
			GrossProfitRate:   0.35,
			FastestMovingDate: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0939888888",
			SquareMeters:      120.0,
			Industry:          "餐飲業",
			Deposit:           50000,
		},
		{
			Title:             "城市健身俱樂部",
			Description:       "提供專業教練課程、最新健身器材，會員數超過1500人，穩定現金流。",
			Price:             2300000,
			Category:          "加盟",
			Condition:         "全新裝修",
			Location:          "台北市大安區信義路四段88號",
			Status:            "活躍",
			OwnerID:           users[0].ID, // John Doe
			ViewCount:         320,
			BrandStory:        "我們秉持『動起來，改變生活』的理念，打造友善社群健身空間。",
			Rent:              60000,
			Floor:             3,
			Equipment:         "跑步機、飛輪、重訓器材、瑜伽室",
			Decoration:        "現代工業風",
			AnnualRevenue:     3200000,
			GrossProfitRate:   0.42,
			FastestMovingDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0922001122",
			SquareMeters:      450.5,
			Industry:          "運動健身",
			Deposit:           300000,
		},
		{
			Title:             "手作甜點工坊",
			Description:       "位於人潮熱區，主打無添加甜點，深受年輕族群喜愛。",
			Price:             550000,
			Category:          "直營",
			Condition:         "8成新",
			Location:          "新北市板橋區文化路一段110號",
			Status:            "活躍",
			OwnerID:           users[2].ID, // Bob Wilson
			ViewCount:         210,
			BrandStory:        "以『健康、純粹、美味』為核心，打造甜點的新標準。",
			Rent:              25000,
			Floor:             1,
			Equipment:         "烤箱、冰箱、甜點工作台",
			Decoration:        "溫馨木質風",
			AnnualRevenue:     800000,
			GrossProfitRate:   0.38,
			FastestMovingDate: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0911777999",
			SquareMeters:      65.0,
			Industry:          "餐飲業",
			Deposit:           80000,
		},
		{
			Title:             "小小森林幼兒園",
			Description:       "已營運5年，生源穩定，位於住宅區，交通便利。",
			Price:             4200000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "高雄市鳳山區建國路222號",
			Status:            "活躍",
			OwnerID:           users[3].ID, // Alice Johnso
			ViewCount:         530,
			BrandStory:        "我們相信教育是改變世界的力量，提供孩子最安心的成長環境。",
			Rent:              90000,
			Floor:             2,
			Equipment:         "教學玩具、課桌椅、投影設備",
			Decoration:        "童趣森林風",
			AnnualRevenue:     5200000,
			GrossProfitRate:   0.28,
			FastestMovingDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0966123456",
			SquareMeters:      800.0,
			Industry:          "教育業",
			Deposit:           500000,
		},
		{
			Title:             "時尚美甲沙龍",
			Description:       "鄰近捷運出口，女性消費者為主，回頭率高。",
			Price:             680000,
			Category:          "直營",
			Condition:         "9成新",
			Location:          "台北市松山區南京東路五段66號",
			Status:            "活躍",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         175,
			BrandStory:        "美，是一種生活態度，我們致力於讓每位客人找到專屬風格。",
			Rent:              35000,
			Floor:             1,
			Equipment:         "美甲機、舒適沙發椅、光療工具",
			Decoration:        "簡約時尚風",
			AnnualRevenue:     950000,
			GrossProfitRate:   0.45,
			FastestMovingDate: time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0955123888",
			SquareMeters:      55.5,
			Industry:          "美容業",
			Deposit:           120000,
		},
		{
			Title:             "電玩樂園",
			Description:       "熱門夜市旁，遊戲機台齊全，小朋友與年輕人聚集地。",
			Price:             1500000,
			Category:          "加盟",
			Condition:         "7成新",
			Location:          "台南市中西區民族路88號",
			Status:            "活躍",
			OwnerID:           users[0].ID, // John Doe
			ViewCount:         410,
			BrandStory:        "打造快樂天堂，讓遊戲連結不同世代的回憶。",
			Rent:              50000,
			Floor:             1,
			Equipment:         "夾娃娃機、賽車機、音樂機台",
			Decoration:        "炫彩娛樂風",
			AnnualRevenue:     2800000,
			GrossProfitRate:   0.33,
			FastestMovingDate: time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0977665544",
			SquareMeters:      200.0,
			Industry:          "娛樂業",
			Deposit:           250000,
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
			Filename:  "coffee_shop_1.jpg",
			URL:       "https://hiyori.cc/wp/wp-content/uploads/2021/03/%E5%AD%B8%E6%A0%A1%E5%92%96%E5%95%A1%E9%A4%A8-Ecole-Cafe7.jpg",
			AltText:   "咖啡店：店外觀",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[1].ID,
			Filename:  "fitness_gym_1.jpg",
			URL:       "https://www.worldgymtaiwan.com/files/club/GF/taipei-gf-pc.jpg",
			AltText:   "健身房：跑步機與重量訓練區",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[2].ID,
			Filename:  "dessert_shop_1.jpg",
			URL:       "https://annieko.tw/wp-content/uploads/20241206171518_0_00d18a.jpg",
			AltText:   "甜點店：店內窗景與座位",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[3].ID,
			Filename:  "kindergarten_1.jpg",
			URL:       "https://www-ws.gov.taipei/Download.ashx?icon=.JPG&n=RFNDMDczMjYuSlBH&u=LzAwMS9VcGxvYWQvNTc5L2NrZmlsZS9lYTlmYTk1MC0yNmZhLTQwYzctYWYwZS0wYTc4MmE3NWYwN2MuanBn",
			AltText:   "幼兒園：教室活動空間",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[4].ID,
			Filename:  "nail_salon_1.jpg",
			URL:       "https://cdn.hippolife.tw/wp-content/uploads/2025/03/19154341/DSC05792-edit-2.webp",
			AltText:   "美甲沙龍：店內座位與裝潢",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[5].ID,
			Filename:  "claw_machine_1.jpg",
			URL:       "https://fupo.tw/wp-content/uploads/2023/01/%E5%8F%B0%E5%8D%97%E5%A8%83%E5%A8%83%E6%A9%9F%E6%8E%A8%E8%96%A6%E5%84%AA%E5%93%81%E5%A8%83%E5%A8%83%E5%B1%8B-3-1.jpg",
			AltText:   "夾娃娃機樂園：店外觀",
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
			ListingID: listings[2].ID, // Dessert Shop
		},
		{
			UserID:    users[1].ID,    // John Doe
			ListingID: listings[4].ID, // Nail Salon
		},
		{
			UserID:    users[3].ID,    // Bob Wilson
			ListingID: listings[1].ID, // Fitness Gym
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
