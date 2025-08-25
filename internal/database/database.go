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
		// Index: 0
		{
			Title:             "好快樂咖啡 ：(",
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
		// Index: 1
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
		// Index: 2
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
		// Index: 3
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
		// Index: 4
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
		// Index: 5
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
		// Index: 6
		{
			Title:             "珍珠研究所（手搖飲）",
			Description:       "每日現煮黑糖珍珠，主打減糖健康風，午晚高峰穩定排隊。",
			Price:             780000,
			Category:          "直營",
			Condition:         "9成新",
			Location:          "台北市信義區永春路100號",
			Status:            "活躍",
			OwnerID:           users[0].ID, // John Doe
			ViewCount:         248,
			BrandStory:        "用最簡單的配方，做最真誠的好味道。",
			Rent:              38000,
			Floor:             1,
			Equipment:         "不鏽鋼工作台、封口機、煮茶鍋、煮珍珠鍋",
			Decoration:        "清新簡約",
			AnnualRevenue:     1100000,
			GrossProfitRate:   0.58,
			FastestMovingDate: time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-111-111",
			SquareMeters:      22.0,
			Industry:          "餐飲業",
			Deposit:           100000,
		},
		// Index: 7
		{
			Title:             "科技便當（外帶快餐）",
			Description:       "鄰近園區，主打高蛋白低油餐盒，合作企業訂單穩定。",
			Price:             1650000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "新竹市東區光復路二段200號",
			Status:            "活躍",
			OwnerID:           users[1].ID, // Jane Smith
			ViewCount:         301,
			BrandStory:        "讓忙碌工程師也能吃得健康又省時。",
			Rent:              52000,
			Floor:             1,
			Equipment:         "四口瓦斯爐、電鍋多台、冷藏展示櫃",
			Decoration:        "工業風",
			AnnualRevenue:     2600000,
			GrossProfitRate:   0.42,
			FastestMovingDate: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-222-222",
			SquareMeters:      48.0,
			Industry:          "餐飲業",
			Deposit:           150000,
		},
		// Index: 8
		{
			Title:             "逗點書店",
			Description:       "社區型閱讀空間，導入選書策展與講座活動，會員制經營。",
			Price:             520000,
			Category:          "直營",
			Condition:         "8成新",
			Location:          "台中市北區文心路一段220號",
			Status:            "活躍",
			OwnerID:           users[2].ID, // Bob Wilson
			ViewCount:         187,
			BrandStory:        "在繁忙城市裡，留下讓人喘口氣的閱讀逗點。",
			Rent:              23000,
			Floor:             1,
			Equipment:         "書架、收銀機、條碼系統、活動投影機",
			Decoration:        "木質暖色",
			AnnualRevenue:     720000,
			GrossProfitRate:   0.32,
			FastestMovingDate: time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-333-333",
			SquareMeters:      36.0,
			Industry:          "零售業",
			Deposit:           60000,
		},
		// Index: 9
		{
			Title:             "微笑洗衣店（自助+代洗）",
			Description:       "24小時自助洗烘加代洗服務，社區大樓密集，回頭率高。",
			Price:             980000,
			Category:          "加盟",
			Condition:         "9成新",
			Location:          "高雄市苓雅區三多一路88號",
			Status:            "活躍",
			OwnerID:           users[3].ID, // Alice Johnso
			ViewCount:         269,
			BrandStory:        "把生活的小麻煩交給我們，換你更多的微笑時光。",
			Rent:              40000,
			Floor:             1,
			Equipment:         "投幣洗衣機×8、烘衣機×6、摺衣桌",
			Decoration:        "亮色清爽",
			AnnualRevenue:     1450000,
			GrossProfitRate:   0.47,
			FastestMovingDate: time.Date(2025, 11, 5, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-444-444",
			SquareMeters:      50.0,
			Industry:          "生活服務",
			Deposit:           120000,
		},
		// Index: 10
		{
			Title:             "小橘子花店",
			Description:       "婚禮佈置＋節慶禮盒，企業合作穩定，線上下單系統完整。",
			Price:             680000,
			Category:          "直營",
			Condition:         "9成新",
			Location:          "台南市安平區安北路300號",
			Status:            "活躍",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         214,
			BrandStory:        "用花朵，把日常的平凡變成值得紀念的驚喜。",
			Rent:              26000,
			Floor:             1,
			Equipment:         "冷藏花庫、修剪工具、包裝台",
			Decoration:        "法式小清新",
			AnnualRevenue:     950000,
			GrossProfitRate:   0.55,
			FastestMovingDate: time.Date(2025, 9, 25, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-555-555",
			SquareMeters:      28.0,
			Industry:          "零售業",
			Deposit:           80000,
		},
		// Index: 11
		{
			Title:             "沐日瑜珈",
			Description:       "小班制與孕婦課專班，周邊商品與線上課程營收成長。",
			Price:             1250000,
			Category:          "直營",
			Condition:         "全新裝修",
			Location:          "桃園市中壢區中山東路二段160號",
			Status:            "活躍",
			OwnerID:           users[0].ID, // John Doe
			ViewCount:         162,
			BrandStory:        "在呼吸之間，與自己重新對話。",
			Rent:              38000,
			Floor:             2,
			Equipment:         "瑜珈墊、輔具、空間音響、濕度控制",
			Decoration:        "日系無印風",
			AnnualRevenue:     1750000,
			GrossProfitRate:   0.48,
			FastestMovingDate: time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-666-666",
			SquareMeters:      90.0,
			Industry:          "運動健身",
			Deposit:           150000,
		},
		// Index: 12
		{
			Title:             "小日子攝影工作室",
			Description:       "親子＆形象照為主，附妝髮區與自然光棚，社群口碑佳。",
			Price:             880000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "新北市新店區北新路二段150號",
			Status:            "活躍",
			OwnerID:           users[1].ID, // Jane Smith
			ViewCount:         141,
			BrandStory:        "把平凡的一天，拍成值得珍藏的一天。",
			Rent:              42000,
			Floor:             3,
			Equipment:         "棚燈三組、反光板、背景紙、4K修圖螢幕",
			Decoration:        "極簡自然光",
			AnnualRevenue:     1350000,
			GrossProfitRate:   0.52,
			FastestMovingDate: time.Date(2025, 11, 12, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-777-777",
			SquareMeters:      65.0,
			Industry:          "攝影服務",
			Deposit:           120000,
		},
		// Index: 13
		{
			Title:             "海風旅店（簡約旅宿）",
			Description:       "步行可到港區與夜市，滿房率穩定，OTA 評價 4.6。",
			Price:             5200000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "基隆市仁愛區愛三路60號",
			Status:            "活躍",
			OwnerID:           users[2].ID, // Bob Wilson
			ViewCount:         403,
			BrandStory:        "在海風裡醒來，旅行也有家的溫度。",
			Rent:              0,
			Floor:             5,
			Equipment:         "客房10間、前台系統、清潔備品",
			Decoration:        "海洋風",
			AnnualRevenue:     6800000,
			GrossProfitRate:   0.39,
			FastestMovingDate: time.Date(2025, 10, 5, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-888-888",
			SquareMeters:      480.0,
			Industry:          "旅宿業",
			Deposit:           600000,
		},
		// Index: 14
		{
			Title:             "漁夫海味小舖",
			Description:       "嚴選產地直送海鮮，冷凍宅配與門市並行，節慶檔期爆量。",
			Price:             1350000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "屏東縣東港鎮中正路110號",
			Status:            "活躍",
			OwnerID:           users[3].ID, // Alice Johnso
			ViewCount:         199,
			BrandStory:        "從海上到餐桌，縮短美味的距離。",
			Rent:              22000,
			Floor:             1,
			Equipment:         "冷凍櫃、真空包裝機、溫控物流合作",
			Decoration:        "藍白海港風",
			AnnualRevenue:     2200000,
			GrossProfitRate:   0.31,
			FastestMovingDate: time.Date(2025, 9, 28, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0912-999-999",
			SquareMeters:      42.0,
			Industry:          "生鮮零售",
			Deposit:           100000,
		},
		// Index: 15
		{
			Title:             "山谷民宿咖啡",
			Description:       "山景第一排，下午茶＋住宿一泊二食方案，假日爆滿。",
			Price:             3900000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "花蓮縣花蓮市中正路50號",
			Status:            "活躍",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         356,
			BrandStory:        "在山與雲的中間，留一席給咖啡與你。",
			Rent:              0,
			Floor:             2,
			Equipment:         "義式咖啡機、烤箱、房務清潔設備",
			Decoration:        "自然木質",
			AnnualRevenue:     5200000,
			GrossProfitRate:   0.37,
			FastestMovingDate: time.Date(2025, 12, 2, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-111-000",
			SquareMeters:      380.0,
			Industry:          "旅宿餐飲",
			Deposit:           450000,
		},
		// Index: 16
		{
			Title:             "青田文具行",
			Description:       "鄰近校園，開學季營收高峰，客製化印章刻印服務。",
			Price:             430000,
			Category:          "直營",
			Condition:         "8成新",
			Location:          "宜蘭縣羅東鎮中正路210號",
			Status:            "活躍",
			OwnerID:           users[0].ID, // John Doe
			ViewCount:         133,
			BrandStory:        "用文具陪伴每一段學習與創作。",
			Rent:              18000,
			Floor:             1,
			Equipment:         "POS、影印機、刻印機、展示架",
			Decoration:        "實用陳列",
			AnnualRevenue:     620000,
			GrossProfitRate:   0.28,
			FastestMovingDate: time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-222-000",
			SquareMeters:      30.0,
			Industry:          "零售業",
			Deposit:           50000,
		},
		// Index: 17
		{
			Title:             "春田機車行",
			Description:       "保養維修、事故協力、外送車隊合作，地點醒目。",
			Price:             850000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "苗栗縣竹南鎮博愛街90號",
			Status:            "活躍",
			OwnerID:           users[1].ID, // Jane Smith
			ViewCount:         177,
			BrandStory:        "讓每天的通勤更安全、更放心。",
			Rent:              20000,
			Floor:             1,
			Equipment:         "舉升機、氣動工具、電瓶測試儀",
			Decoration:        "機能取向",
			AnnualRevenue:     1350000,
			GrossProfitRate:   0.36,
			FastestMovingDate: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-333-000",
			SquareMeters:      55.0,
			Industry:          "維修服務",
			Deposit:           100000,
		},
		// Index: 18
		{
			Title:             "家家乾洗",
			Description:       "社區代收點多據點合作，禮服與西裝精緻洗護口碑好。",
			Price:             720000,
			Category:          "加盟",
			Condition:         "9成新",
			Location:          "新竹縣竹北市文興路100號",
			Status:            "活躍",
			OwnerID:           users[2].ID, // Bob Wilson
			ViewCount:         159,
			BrandStory:        "為每一件衣服恢復初見時的心動。",
			Rent:              28000,
			Floor:             1,
			Equipment:         "水洗機、乾洗機、蒸氣熨燙台",
			Decoration:        "明亮整潔",
			AnnualRevenue:     1200000,
			GrossProfitRate:   0.41,
			FastestMovingDate: time.Date(2025, 11, 8, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-444-000",
			SquareMeters:      40.0,
			Industry:          "生活服務",
			Deposit:           90000,
		},
		// Index: 19
		{
			Title:             "玩具倉庫（親子選物）",
			Description:       "益智教具與桌遊為主，假日親子活動帶動銷售。",
			Price:             690000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "台北市士林區文林路150號",
			Status:            "活躍",
			OwnerID:           users[3].ID, // Alice Johnso
			ViewCount:         201,
			BrandStory:        "把快樂變成能分享的禮物。",
			Rent:              37000,
			Floor:             1,
			Equipment:         "展示層架、收銀系統、活動區桌椅",
			Decoration:        "繽紛童趣",
			AnnualRevenue:     1120000,
			GrossProfitRate:   0.34,
			FastestMovingDate: time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-555-000",
			SquareMeters:      45.0,
			Industry:          "零售業",
			Deposit:           110000,
		},
		// Index: 20
		{
			Title:             "豆香手工豆花",
			Description:       "古早味路線，使用非基改黃豆，每日限量售完為止。",
			Price:             430000,
			Category:          "直營",
			Condition:         "8成新",
			Location:          "嘉義市西區文化路120號",
			Status:            "活躍",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         188,
			BrandStory:        "一碗豆花，留住童年的味道。",
			Rent:              16000,
			Floor:             1,
			Equipment:         "蒸煮鍋、冷藏櫃、保溫桶",
			Decoration:        "復古小店",
			AnnualRevenue:     680000,
			GrossProfitRate:   0.49,
			FastestMovingDate: time.Date(2025, 9, 26, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-666-000",
			SquareMeters:      20.0,
			Industry:          "餐飲業",
			Deposit:           50000,
		},
		// Index: 21
		{
			Title:             "稻香便當站",
			Description:       "強調產地溯源的白米與在地蔬菜，外送佔比 40%。",
			Price:             780000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "台東縣池上鄉中正路88號",
			Status:            "活躍",
			OwnerID:           users[0].ID, // John Doe
			ViewCount:         144,
			BrandStory:        "用好米，做出記憶中的家常味。",
			Rent:              12000,
			Floor:             1,
			Equipment:         "電鍋、保溫餐車、冷藏展示櫃",
			Decoration:        "樸實清爽",
			AnnualRevenue:     980000,
			GrossProfitRate:   0.43,
			FastestMovingDate: time.Date(2025, 10, 18, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-777-000",
			SquareMeters:      26.0,
			Industry:          "餐飲業",
			Deposit:           80000,
		},
		// Index: 22
		{
			Title:             "晨光托育園",
			Description:       "鄰近公園，戶外活動空間大，社區口碑高。",
			Price:             4200000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "新竹縣新豐鄉建興路60號",
			Status:            "活躍",
			OwnerID:           users[1].ID, // Jane Smith
			ViewCount:         329,
			BrandStory:        "把安全與愛，變成每天可見的日常。",
			Rent:              68000,
			Floor:             2,
			Equipment:         "教具、監視系統、室外遊具",
			Decoration:        "童趣自然",
			AnnualRevenue:     5600000,
			GrossProfitRate:   0.27,
			FastestMovingDate: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-888-000",
			SquareMeters:      600.0,
			Industry:          "教育業",
			Deposit:           500000,
		},
		// Index: 23
		{
			Title:             "髮藝沙龍（三重）",
			Description:       "捷運商圈，燙染護比例高，會員儲值穩定。",
			Price:             980000,
			Category:          "直營",
			Condition:         "9成新",
			Location:          "新北市三重區重新路三段120號",
			Status:            "活躍",
			OwnerID:           users[2].ID, // Bob Wilson
			ViewCount:         246,
			BrandStory:        "髮絲之間，讓自信自然流露。",
			Rent:              45000,
			Floor:             2,
			Equipment:         "洗髮椅、造型椅、染燙設備",
			Decoration:        "都會簡約",
			AnnualRevenue:     1750000,
			GrossProfitRate:   0.46,
			FastestMovingDate: time.Date(2025, 10, 7, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0920-999-000",
			SquareMeters:      70.0,
			Industry:          "美容美髮",
			Deposit:           150000,
		},
		// Index: 24
		{
			Title:             "創客共用空間",
			Description:       "3D列印、雷射切割、社群講座，每月固定會員 120+。",
			Price:             2100000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "台中市西區公益路200號",
			Status:            "活躍",
			OwnerID:           users[3].ID, // Alice Johnso
			ViewCount:         318,
			BrandStory:        "把點子做成作品，把作品變成事業。",
			Rent:              98000,
			Floor:             3,
			Equipment:         "3D印表機×6、雷射切割機、工作台",
			Decoration:        "開放工坊風",
			AnnualRevenue:     3600000,
			GrossProfitRate:   0.33,
			FastestMovingDate: time.Date(2025, 11, 20, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0930-111-222",
			SquareMeters:      320.0,
			Industry:          "共享空間",
			Deposit:           300000,
		},
		// Index: 25
		{
			Title:             "晨曦烘焙坊",
			Description:       "每日現烤歐式麵包與天然酵母，下午出爐秒殺。",
			Price:             980000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "雲林縣斗六市中山路66號",
			Status:            "活躍",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         207,
			BrandStory:        "用時間換來的麥香，值得等候。",
			Rent:              21000,
			Floor:             1,
			Equipment:         "雙層烤箱、發酵箱、行星攪拌機",
			Decoration:        "歐式鄉村",
			AnnualRevenue:     1600000,
			GrossProfitRate:   0.44,
			FastestMovingDate: time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0930-222-333",
			SquareMeters:      38.0,
			Industry:          "餐飲業",
			Deposit:           120000,
		},
		// Index: 26
		{
			Title:             "樂活寵物美容",
			Description:       "犬貓洗護＋基礎訓練，周邊商品搭配銷售。",
			Price:             830000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "新北市板橋區文化路二段88號",
			Status:            "活躍",
			OwnerID:           users[0].ID, // John Doe
			ViewCount:         173,
			BrandStory:        "讓毛孩更舒服，讓飼主更放心。",
			Rent:              33000,
			Floor:             1,
			Equipment:         "美容桌、烘箱、吹水機、剪具",
			Decoration:        "溫馨寵物友善",
			AnnualRevenue:     1250000,
			GrossProfitRate:   0.45,
			FastestMovingDate: time.Date(2025, 10, 2, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0930-333-444",
			SquareMeters:      40.0,
			Industry:          "寵物服務",
			Deposit:           100000,
		},
		// Index: 27
		{
			Title:             "清泉自助洗車",
			Description:       "雙車位＋吸塵區，鄰近社區停車場，夜間人流穩定。",
			Price:             1680000,
			Category:          "直營",
			Condition:         "9成新",
			Location:          "桃園市桃園區中華路500號",
			Status:            "活躍",
			OwnerID:           users[1].ID, // Jane Smith
			ViewCount:         220,
			BrandStory:        "讓車子在十分鐘內煥然一新。",
			Rent:              45000,
			Floor:             1,
			Equipment:         "高壓水柱、泡沫槍、投幣吸塵器",
			Decoration:        "戶外站點",
			AnnualRevenue:     2300000,
			GrossProfitRate:   0.51,
			FastestMovingDate: time.Date(2025, 11, 3, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0930-444-555",
			SquareMeters:      180.0,
			Industry:          "汽車服務",
			Deposit:           250000,
		},
		// Index: 28
		{
			Title:             "亮亮眼鏡館",
			Description:       "醫師配鏡合作、快速取件，學生與上班族客群穩定。",
			Price:             1150000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "台北市中山區南京東路二段120號",
			Status:            "活躍",
			OwnerID:           users[2].ID, // Bob Wilson
			ViewCount:         195,
			BrandStory:        "讓視界清晰，讓生活更輕鬆。",
			Rent:              52000,
			Floor:             1,
			Equipment:         "驗光儀、研磨機、鏡框展示牆",
			Decoration:        "現代簡約",
			AnnualRevenue:     2100000,
			GrossProfitRate:   0.43,
			FastestMovingDate: time.Date(2025, 10, 14, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0930-555-666",
			SquareMeters:      52.0,
			Industry:          "零售服務",
			Deposit:           180000,
		},
		// Index: 29
		{
			Title:             "十里鍋物（小火鍋）",
			Description:       "個人鍋快翻桌、高 CP 值，外送平台口碑 4.7。",
			Price:             1750000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "新北市永和區中山路一段180號",
			Status:            "活躍",
			OwnerID:           users[3].ID, // Alice Johnso
			ViewCount:         287,
			BrandStory:        "用好湯底，走十里都要回頭吃。",
			Rent:              68000,
			Floor:             1,
			Equipment:         "商用電磁爐、冷藏冷凍庫、前場點餐系統",
			Decoration:        "溫暖木質",
			AnnualRevenue:     3600000,
			GrossProfitRate:   0.38,
			FastestMovingDate: time.Date(2025, 11, 11, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0930-666-777",
			SquareMeters:      120.0,
			Industry:          "餐飲業",
			Deposit:           300000,
		},
		// Index: 30
		{
			Title:             "學園家教中心",
			Description:       "國高中數理專班，小班制與一對一並行，升學績效佳。",
			Price:             2350000,
			Category:          "直營",
			Condition:         "良好",
			Location:          "台南市東區東寧路260號",
			Status:            "活躍",
			OwnerID:           users[4].ID, // Alice Johnson
			ViewCount:         334,
			BrandStory:        "讓學習變得有方法、有成就感。",
			Rent:              58000,
			Floor:             3,
			Equipment:         "白板、投影機、分組教室、講義系統",
			Decoration:        "明亮教室",
			AnnualRevenue:     5200000,
			GrossProfitRate:   0.29,
			FastestMovingDate: time.Date(2025, 10, 30, 0, 0, 0, 0, time.UTC),
			PhoneNumber:       "0930-777-888",
			SquareMeters:      240.0,
			Industry:          "教育業",
			Deposit:           350000,
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
		// Index: 0 - 快樂咖啡館
		{
			ListingID: listings[0].ID,
			Filename:  "happy_coffee.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/happy_coffee.jpg",
			AltText:   "快樂咖啡館",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 1 - 樂活寵物美容
		{
			ListingID: listings[1].ID,
			Filename:  "pet_grooming.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/pet_grooming.jpg",
			AltText:   "寵物美容：美容台與設備",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 2 - 髮藝沙龍（三重）
		{
			ListingID: listings[2].ID,
			Filename:  "nail_salon.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/nail_salon.jpg",
			AltText:   "髮藝沙龍：造型座位區",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 3 - 晨曦烘焙坊
		{
			ListingID: listings[3].ID,
			Filename:  "bakery.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/bakery.jpg",
			AltText:   "烘焙坊：麵包陳列櫃",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 4 - 創客共用空間
		{
			ListingID: listings[4].ID,
			Filename:  "photo_studio.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/photo_studio.jpg",
			AltText:   "創客空間：工作檯與設備",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 5 - 稻香便當站
		{
			ListingID: listings[5].ID,
			Filename:  "bento_shop.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/bento_shop.jpg",
			AltText:   "便當店：餐盒展示",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 6 - 豆香手工豆花
		{
			ListingID: listings[6].ID,
			Filename:  "dessert_shop.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/dessert_shop.jpg",
			AltText:   "豆花店：甜品陳列",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 7 - 玩具倉庫（親子選物）
		{
			ListingID: listings[7].ID,
			Filename:  "toy_store.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/toy_store.jpg",
			AltText:   "玩具店：商品陳列",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 8 - 家家乾洗
		{
			ListingID: listings[8].ID,
			Filename:  "dry_cleaning.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/dry_cleaning.jpg",
			AltText:   "乾洗店：洗衣設備",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 9 - 春田機車行
		{
			ListingID: listings[9].ID,
			Filename:  "scooter_shop.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/scooter_shop.jpg",
			AltText:   "機車行：維修區",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 10 - 青田文具行
		{
			ListingID: listings[10].ID,
			Filename:  "stationery_store.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/stationery_store.jpg",
			AltText:   "文具行：商品陳列",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 11 - 沐日瑜珈
		{
			ListingID: listings[11].ID,
			Filename:  "yoga_studio.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/yoga_studio.jpg",
			AltText:   "瑜珈教室：練習空間",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 12 - 小日子攝影工作室
		{
			ListingID: listings[12].ID,
			Filename:  "photo_studio.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/photo_studio.jpg",
			AltText:   "攝影工作室：拍攝空間",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 13 - 海風旅店（簡約旅宿）
		{
			ListingID: listings[13].ID,
			Filename:  "hotel_room.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/hotel_room.jpg",
			AltText:   "旅店：客房環境",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 14 - 漁夫海味小舖
		{
			ListingID: listings[14].ID,
			Filename:  "seafood_market.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/seafood_market.jpg",
			AltText:   "海鮮店：新鮮海產",
			Order:     0,
			IsPrimary: true,
		},
		// Index: 15 - 山谷民宿咖啡
		{
			ListingID: listings[15].ID,
			Filename:  "mountain_cafe.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/mountain_cafe.jpg",
			AltText:   "山谷咖啡：景觀座位",
			Order:     0,
			IsPrimary: true,
		},
		{ListingID: listings[16].ID, Filename: "stationery_store.jpg", URL: "http://127.0.0.1:8080/static/images/listings/stationery_store.jpg", AltText: "文具店通道與貨架示意", Order: 0, IsPrimary: true},
		{ListingID: listings[17].ID, Filename: "scooter_repair_shop.jpg", URL: "http://127.0.0.1:8080/static/images/listings/scooter_shop.jpg", AltText: "機車維修工位與工具示意", Order: 0, IsPrimary: true},
		{ListingID: listings[18].ID, Filename: "dry_clean_shop.jpg", URL: "http://127.0.0.1:8080/static/images/listings/dry_cleaning.jpg", AltText: "洗烘併設的洗衣空間示意", Order: 0, IsPrimary: true},
		{ListingID: listings[19].ID, Filename: "toy_store_aisle.jpg", URL: "http://127.0.0.1:8080/static/images/listings/toy_store.jpg", AltText: "玩具賣場走道與貨架示意", Order: 0, IsPrimary: true},
		{ListingID: listings[20].ID, Filename: "douhua_shop.jpg", URL: "http://127.0.0.1:8080/static/images/listings/dessert_shop.jpg", AltText: "豆花甜品與內用座位示意", Order: 0, IsPrimary: true},
		{ListingID: listings[21].ID, Filename: "bento_counter.jpg", URL: "http://127.0.0.1:8080/static/images/listings/bento_shop.jpg", AltText: "便當餐盒展示與出餐示意", Order: 0, IsPrimary: true},
		{ListingID: listings[22].ID, Filename: "daycare_classroom.jpg", URL: "http://127.0.0.1:8080/static/images/listings/kindergarten.jpg", AltText: "幼兒教室與遊戲區示意", Order: 0, IsPrimary: true},
		{ListingID: listings[23].ID, Filename: "hair_salon_interior.jpg", URL: "http://127.0.0.1:8080/static/images/listings/nail_salon.jpg", AltText: "髮廊工業風內裝與座位示意", Order: 0, IsPrimary: true},
		{ListingID: listings[24].ID, Filename: "makerspace_workshop.jpg", URL: "http://127.0.0.1:8080/static/images/listings/photo_studio.jpg", AltText: "3D列印與手作空間示意", Order: 0, IsPrimary: true},
		{ListingID: listings[25].ID, Filename: "bakery_storefront.jpg", URL: "http://127.0.0.1:8080/static/images/listings/bakery.jpg", AltText: "烘焙坊麵包陳列與店面示意", Order: 0, IsPrimary: true},
		{
			ListingID: listings[26].ID,
			Filename:  "pet_grooming_waiting_area.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/pet_grooming.jpg",
			AltText:   "寵物美容：明亮等待區與遊戲室",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[27].ID,
			Filename:  "self_service_car_wash.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/car_wash.jpg",
			AltText:   "自助洗車場：戶外洗車隔間",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[28].ID,
			Filename:  "eyeglass_store_interior.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/eyeglass_store.jpg",
			AltText:   "眼鏡門市：展示區與鏡架",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[29].ID,
			Filename:  "hotpot_interior.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/hotpot_restaurant.jpg",
			AltText:   "小火鍋店：內用空間與自助吧",
			Order:     0,
			IsPrimary: true,
		},
		{
			ListingID: listings[30].ID,
			Filename:  "tutoring_classroom.jpg",
			URL:       "http://127.0.0.1:8080/static/images/listings/tutoring_center.jpg",
			AltText:   "家教/補習班：電腦教室座位",
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

	// Create sample favorites (ensure no duplicates)
	favorites := []models.Favorite{
		{
			UserID:    users[0].ID,    // John Doe
			ListingID: listings[1].ID, // Fitness Gym
		},
		{
			UserID:    users[1].ID,    // Jane Smith
			ListingID: listings[0].ID, // Coffee Shop
		},
		{
			UserID:    users[2].ID,    // Bob Wilson
			ListingID: listings[3].ID, // Kindergarten
		},
		{
			UserID:    users[2].ID,    // Bob Wilson
			ListingID: listings[5].ID, // Game Arcade
		},
		{
			UserID:    users[3].ID,    // Alice Johnson (user index 3)
			ListingID: listings[2].ID, // Dessert Shop
		},
		{
			UserID:    users[4].ID,    // Alice Johnson (user index 4)
			ListingID: listings[4].ID, // Nail Salon
		},
		{
			UserID:    users[1].ID,    // Jane Smith
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
