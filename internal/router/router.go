package router

import (
	"net/http"
	"strings"
	"time"

	"trade_company/graph"
	"trade_company/internal/config"
	gqlctx "trade_company/internal/graphql"
	"trade_company/internal/handlers"
	"trade_company/internal/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, log *zap.Logger, db *gorm.DB, _ *redis.Client) http.Handler {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger(log))
	r.Use(middleware.RequestID())
	r.Use(corsMiddleware(cfg))

	// load templates
	r.LoadHTMLGlob("templates/*.html")

	// pages
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	r.GET("/dashboard", func(c *gin.Context) { c.HTML(http.StatusOK, "dashboard.html", nil) })

	// health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// REST API v1
	authH := &handlers.AuthHandler{DB: db, Cfg: cfg}
	listH := &handlers.ListingsHandler{DB: db}
	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/login", authH.Login)

		api.GET("/listings", listH.List)
		api.GET("/listings/:id", listH.Get)
		authd := api.Group("")
		authd.Use(middleware.JWTAuth(cfg))
		authd.POST("/listings", listH.Create)
	}

	// GraphQL
	es := graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db, Cfg: cfg}})
	gh := handler.NewDefaultServer(es)

	graphqlGroup := r.Group("")
	graphqlGroup.Use(func(c *gin.Context) {
		// enrich request context with userID if token provided
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
		log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
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
