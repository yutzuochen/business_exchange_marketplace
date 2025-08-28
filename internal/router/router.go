package router

import (
	logOri "log"
	"net/http"
	"strings"
	"time"

	"trade_company/graph"
	"trade_company/internal/config"
	gqlctx "trade_company/internal/graphql"
	"trade_company/internal/handlers"
	"trade_company/internal/middleware"
	"trade_company/internal/models"

	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, log *zap.Logger, db *gorm.DB, redisClient *redis.Client) http.Handler {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(requestLogger(log))

	// Load templates
	r.LoadHTMLGlob("templates/*.html")

	// Static files
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")

	// Health check endpoints
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "ok",
			"timestamp":  time.Now().UTC(),
			"request_id": c.GetString("request_id"),
		})
	}
	r.GET("/health", healthHandler)
	r.GET("/healthz", healthHandler)

	// Public pages
	r.GET("/", func(c *gin.Context) {
		var txs []models.Transaction
		var listings []models.Listing

		if db != nil {
			_ = db.Order("created_at desc").Limit(10).Find(&txs).Error
			_ = db.Order("id desc").Limit(8).Find(&listings).Error
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"transactions": txs,
			"listings":     listings,
		})
	})

	r.GET("/market", func(c *gin.Context) {
		var txs []models.Transaction
		var listings []models.Listing

		if db != nil {
			_ = db.Order("created_at desc").Limit(10).Find(&txs).Error
			_ = db.Order("id desc").Limit(8).Find(&listings).Error
		}

		c.HTML(http.StatusOK, "market_home.html", gin.H{
			"transactions": txs,
			"listings":     listings,
			"listingPriceRanges": func() []map[string]interface{} {
				ranges := make([]map[string]interface{}, len(listings))
				for i, l := range listings {
					low := int64(float64(l.Price) * 0.85)
					high := int64(float64(l.Price) * 1.15)
					ranges[i] = map[string]interface{}{
						"id":    l.ID,
						"low":   low,
						"high":  high,
						"price": l.Price,
					}
				}
				return ranges
			}(),
		})
	})

	// Search listing by title and redirect to detail page if found
	r.GET("/market/search", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" || db == nil {
			c.Redirect(http.StatusFound, "/market")
			return
		}
		var ls models.Listing
		if err := db.Where("title LIKE ?", "%"+q+"%").Order("id desc").First(&ls).Error; err != nil {
			c.Redirect(http.StatusFound, "/market")
			return
		}
		c.Redirect(http.StatusFound, "/market/listings/"+strconv.FormatUint(uint64(ls.ID), 10))
	})

	// Listing detail page
	r.GET("/market/listings/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if db == nil {
			c.String(http.StatusServiceUnavailable, "database not available")
			return
		}
		var ls models.Listing
		if err := db.First(&ls, idStr).Error; err != nil {
			c.String(http.StatusNotFound, "listing not found")
			return
		}
		var images []models.Image
		_ = db.Where("listing_id = ?", ls.ID).Order("id asc").Find(&images).Error
		// log.Printf("Go syntax: %#v\n", p)
		logOri.Printf("===== LS: %+v\n", ls)
		c.HTML(http.StatusOK, "market_listing.html", gin.H{
			"listing": ls,
			"images":  images,
		})
	})

	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	r.GET("/dashboard", func(c *gin.Context) { c.HTML(http.StatusOK, "dashboard.html", nil) })

	// REST API v1
	authH := &handlers.AuthHandler{DB: db, Cfg: cfg, Log: log}
	listH := &handlers.ListingsHandler{DB: db}
	userH := &handlers.UserHandler{DB: db}
	favH := &handlers.FavoriteHandler{DB: db}
	msgH := &handlers.MessageHandler{DB: db}

	api := r.Group("/api/v1")
	{
		// Public endpoints
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/login", authH.Login)
		api.GET("/listings", listH.List)
		api.GET("/listings/:id", listH.Get)
		api.GET("/categories", listH.GetCategories)

		// Protected endpoints
		authd := api.Group("")
		authd.Use(middleware.JWT(middleware.JWTConfig{
			Secret: cfg.JWTSecret,
			Issuer: cfg.JWTIssuer,
		}, log))
		{
			// User management
			authd.GET("/user/profile", userH.GetProfile)
			authd.PUT("/user/profile", userH.UpdateProfile)
			authd.PUT("/user/password", userH.ChangePassword)

			// Listings
			authd.POST("/listings", listH.Create)
			authd.PUT("/listings/:id", listH.Update)
			authd.DELETE("/listings/:id", listH.Delete)
			authd.POST("/listings/:id/images", listH.UploadImages)

			// Favorites
			authd.GET("/favorites", favH.List)
			authd.POST("/favorites", favH.Add)
			authd.DELETE("/favorites/:id", favH.Remove)

			// Messages
			authd.GET("/messages", msgH.List)
			authd.GET("/messages/:id", msgH.Get)
			authd.POST("/messages", msgH.Create)
			authd.PUT("/messages/:id/read", msgH.MarkAsRead)
		}
	}

	// GraphQL
	es := graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db, Cfg: cfg}})
	gh := handler.NewDefaultServer(es)

	graphqlGroup := r.Group("")
	graphqlGroup.Use(func(c *gin.Context) {
		// Enrich request context with userID if token provided
		ctx := gqlctx.ExtractUserFromAuthHeader(cfg, c.Request.Context(), c.GetHeader("Authorization"))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	graphqlGroup.POST("/graphql", gin.WrapH(gh))
	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL", "/graphql")))

	return r
}

func requestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)

		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = "unknown"
		}

		log.Info("request",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", dur),
		)
	}
}

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := strings.Split(cfg.CORSAllowedOrigins, ",")
	allowedMethods := cfg.CORSAllowedMethods
	allowedHeaders := cfg.CORSAllowedHeaders

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && (cfg.CORSAllowedOrigins == "*" || contains(allowedOrigins, origin)) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		} else if cfg.CORSAllowedOrigins == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", allowedMethods)
		c.Header("Access-Control-Allow-Headers", allowedHeaders)
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) == target {
			return true
		}
	}
	return false
}
