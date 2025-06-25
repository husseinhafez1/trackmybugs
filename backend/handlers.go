package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func registerHandler(c *gin.Context) {
	var registerData struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Role      string `json:"role"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Set default role
	if registerData.Role == "" {
		registerData.Role = "user"
	}

	// Create user
	user := User{
		ID:           uuid.New().String(),
		Email:        registerData.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    registerData.FirstName,
		LastName:     registerData.LastName,
		Role:         registerData.Role,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	if err := createUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func createUser(user User) error {
	query := `
	INSERT INTO users (id, email, password_hash, first_name, last_name, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := db.Exec(query, user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Printf("Database error creating user: %v", err)
	}
	return err
}

func loginHandler(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	user, err := getUserByEmail(loginData.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func getUserByEmail(email string) (User, error) {
	var user User
	query := `
	SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
	FROM users
	WHERE email = $1
	`
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func getProjectsHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	projects, err := getProjectsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"message":  "Projects fetched successfully",
	})
}

func getProjectsByUser(userID string) ([]Project, error) {
	query := `
	SELECT id, name, description, created_by, created_at, updated_at
	FROM projects
	WHERE created_by = $1
	ORDER BY created_at DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func createProjectHandler(c *gin.Context) {
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = uuid.New().String()
	project.CreatedBy = c.GetString("user_id")
	project.CreatedAt = time.Now().Format(time.RFC3339)
	project.UpdatedAt = project.CreatedAt

	if err := createProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

func createProject(project Project) error {
	query := `
	INSERT INTO projects (id, name, description, created_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(query, project.ID, project.Name, project.Description, project.CreatedBy, project.CreatedAt, project.UpdatedAt)
	return err
}

func getProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	project, err := getProjectByID(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func getProjectByID(projectID string) (Project, error) {
	var project Project
	query := `
	SELECT id, name, description, created_by, created_at, updated_at
	FROM projects
	WHERE id = $1
	`
	err := db.QueryRow(query, projectID).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt)
	return project, err
}

func updateProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	project, err := getProjectByID(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := updateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

func updateProject(project Project) error {
	query := `
	UPDATE projects
	SET name = $1, description = $2, updated_at = $3
	WHERE id = $4
	`
	_, err := db.Exec(query, project.Name, project.Description, project.UpdatedAt, project.ID)
	return err
}

func deleteProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	if err := deleteProject(projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func deleteProject(projectID string) error {
	query := `
	DELETE FROM projects
	WHERE id = $1
	`
	_, err := db.Exec(query, projectID)
	return err
}

// Issues handlers
func getIssuesHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get issues - not implemented yet"})
}

func createIssueHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Create issue - not implemented yet"})
}

func getIssueHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get issue - not implemented yet"})
}

func updateIssueHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update issue - not implemented yet"})
}

func deleteIssueHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete issue - not implemented yet"})
}

// Comments handlers
func getCommentsHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get comments - not implemented yet"})
}

func createCommentHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Create comment - not implemented yet"})
}

func updateCommentHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update comment - not implemented yet"})
}

func deleteCommentHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete comment - not implemented yet"})
}

// Users handlers
func getUsersHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get users - not implemented yet"})
}

func getProfileHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get profile - not implemented yet"})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["sub"].(string)
			c.Set("user_id", userID)
		}

		c.Next()
	}
}
