package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Database connection
var db *sql.DB

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	initDB()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "TrackMyBugs API is running",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", registerHandler)
			auth.POST("/login", loginHandler)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware())
		{
			// Projects
			projects := protected.Group("/projects")
			{
				projects.GET("", getProjectsHandler)
				projects.POST("", createProjectHandler)
				projects.GET("/:id", getProjectHandler)
				projects.PUT("/:id", updateProjectHandler)
				projects.DELETE("/:id", adminOnly(), deleteProjectHandler)
			}

			// Issues
			issues := protected.Group("/issues")
			{
				issues.GET("", getIssuesHandler)
				issues.POST("", createIssueHandler)
				issues.GET("/:id", getIssueHandler)
				issues.PUT("/:id", updateIssueHandler)
				issues.DELETE("/:id", deleteIssueHandler)
			}

			// Comments
			comments := protected.Group("/comments")
			{
				comments.GET("/issue/:issueId", getCommentsHandler)
				comments.POST("", createCommentHandler)
				comments.PUT("/:id", updateCommentHandler)
				comments.DELETE("/:id", deleteCommentHandler)
			}

			// Users
			users := protected.Group("/users")
			{
				users.GET("", getUsersHandler)
				users.GET("/profile", getProfileHandler)
				users.PUT("/profile", updateProfileHandler)
				users.PUT("/:id/role", adminOnly(), updateUserRoleHandler)
			}
		}
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Initialize database connection
func initDB() {
	// Get database connection details from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "trackmybugs")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPassword)

	// Open database connection
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
