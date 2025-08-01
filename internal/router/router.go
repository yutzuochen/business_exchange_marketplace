package router

import (
	"business-marketplace/internal/config"
	"business-marketplace/internal/handlers"
	"business-marketplace/internal/middleware"
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Initialize sets up the router with all routes and middleware
func Initialize(cfg *config.Config, db *sql.DB, redisClient *redis.Client, logger *logrus.Logger) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Configure properly for production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Custom middleware
	r.Use(middleware.Logger(logger))
	r.Use(middleware.RateLimit(cfg))

	// Initialize handlers
	h := handlers.New(cfg, db, redisClient, logger)

	// Serve static files
	r.Static("/static", "./static")
	r.Static("/uploads", cfg.UploadPath)

	// Load HTML templates
	r.LoadHTMLGlob("templates/**/*")

	// Health check
	r.GET("/health", h.HealthCheck)

	// Public routes
	public := r.Group("/")
	{
		public.GET("/", h.Home)
		public.GET("/search", h.Search)
		public.GET("/listing/:slug", h.ListingDetail)
		public.GET("/category/:slug", h.CategoryListings)
	}

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.GET("/login", h.LoginPage)
		auth.POST("/login", h.Login)
		auth.GET("/register", h.RegisterPage)
		auth.POST("/register", h.Register)
		auth.POST("/logout", h.Logout)
		auth.GET("/verify-email/:token", h.VerifyEmail)
		auth.GET("/forgot-password", h.ForgotPasswordPage)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.GET("/reset-password/:token", h.ResetPasswordPage)
		auth.POST("/reset-password/:token", h.ResetPassword)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired(cfg, redisClient))
	{
		protected.GET("/dashboard", h.Dashboard)
		protected.GET("/profile", h.Profile)
		protected.POST("/profile", h.UpdateProfile)

		// Listing management
		protected.GET("/listings/create", h.CreateListingPage)
		protected.POST("/listings/create", h.CreateListing)
		protected.GET("/listings/:id/edit", h.EditListingPage)
		protected.POST("/listings/:id/edit", h.UpdateListing)
		protected.DELETE("/listings/:id", h.DeleteListing)
		protected.GET("/my-listings", h.MyListings)

		// Favorites
		protected.POST("/favorites/:id", h.AddToFavorites)
		protected.DELETE("/favorites/:id", h.RemoveFromFavorites)
		protected.GET("/favorites", h.MyFavorites)

		// Inquiries
		protected.POST("/inquiries", h.SendInquiry)
		protected.GET("/inquiries", h.MyInquiries)
		protected.GET("/inquiries/:id", h.InquiryDetail)
		protected.POST("/inquiries/:id/reply", h.ReplyToInquiry)
	}

	// API routes
	api := r.Group("/api/v1")
	{
		// Public API
		api.GET("/listings", h.APIListings)
		api.GET("/listings/:id", h.APIListingDetail)
		api.GET("/categories", h.APICategories)
		api.GET("/search", h.APISearch)

		// Protected API
		apiProtected := api.Group("/")
		apiProtected.Use(middleware.AuthRequired(cfg, redisClient))
		{
			apiProtected.POST("/listings", h.APICreateListing)
			apiProtected.PUT("/listings/:id", h.APIUpdateListing)
			apiProtected.DELETE("/listings/:id", h.APIDeleteListing)
			apiProtected.POST("/upload", h.APIUploadImage)
		}
	}

	return r
}
