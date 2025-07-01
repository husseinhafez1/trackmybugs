package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	// Make the first registered user an admin
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user count"})
		return
	}
	if userCount == 0 {
		registerData.Role = "admin"
	} else if registerData.Role == "" {
		registerData.Role = "user"
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
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
	search := c.Query("search")
	limit := 10
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	total, err := countProjectsByUser(userID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count projects"})
		return
	}

	projects, err := getProjectsByUserPaginated(userID, search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	if projects == nil {
		projects = []Project{}
	}
	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"message":  "Projects fetched successfully",
	})
}

func countProjectsByUser(userID, search string) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM projects WHERE created_by = $1`
	args := []interface{}{userID}
	idx := 2
	if search != "" {
		query += ` AND (LOWER(name) LIKE $` + strconv.Itoa(idx) + ` OR LOWER(description) LIKE $` + strconv.Itoa(idx) + `)`
		searchTerm := "%" + search + "%"
		args = append(args, strings.ToLower(searchTerm))
		idx++
	}
	err := db.QueryRow(query, args...).Scan(&total)
	return total, err
}

func getProjectsByUserPaginated(userID, search string, limit, offset int) ([]Project, error) {
	query := `
	SELECT id, name, description, created_by, created_at, updated_at
	FROM projects
	WHERE created_by = $1`
	args := []interface{}{userID}
	idx := 2
	if search != "" {
		query += ` AND (LOWER(name) LIKE $` + strconv.Itoa(idx) + ` OR LOWER(description) LIKE $` + strconv.Itoa(idx) + `)`
		searchTerm := "%" + search + "%"
		args = append(args, strings.ToLower(searchTerm))
		idx++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(idx) + ` OFFSET $` + strconv.Itoa(idx+1)
	args = append(args, limit, offset)
	rows, err := db.Query(query, args...)
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

// RBAC middleware: only allow admins
func adminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		user, err := getUserByID(userID)
		if err != nil || user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Issues handlers
func getIssuesHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Query("project_id")
	status := c.Query("status")
	priority := c.Query("priority")
	assignedTo := c.Query("assigned_to")
	search := c.Query("search")
	limit := 10
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	var issues []Issue
	var total int
	var err error

	if projectID != "" {
		total, err = countIssuesByProjectFiltered(projectID, status, priority, assignedTo, search)
		if err == nil {
			issues, err = getIssuesByProjectPaginatedFiltered(projectID, status, priority, assignedTo, search, limit, offset)
		}
	} else {
		total, err = countIssuesByUserFiltered(userID, status, priority, assignedTo, search)
		if err == nil {
			issues, err = getIssuesByUserPaginatedFiltered(userID, status, priority, assignedTo, search, limit, offset)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch issues"})
		return
	}
	if issues == nil {
		issues = []Issue{}
	}
	c.JSON(http.StatusOK, gin.H{
		"issues":  issues,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"message": "Issues fetched successfully",
	})
}

func countIssuesByProjectFiltered(projectID, status, priority, assignedTo, search string) (int, error) {
	query := `SELECT COUNT(*) FROM issues WHERE project_id = $1`
	args := []interface{}{projectID}
	idx := 2
	if status != "" {
		query += ` AND status = $` + strconv.Itoa(idx)
		args = append(args, status)
		idx++
	}
	if priority != "" {
		query += ` AND priority = $` + strconv.Itoa(idx)
		args = append(args, priority)
		idx++
	}
	if assignedTo != "" {
		query += ` AND assigned_to = $` + strconv.Itoa(idx)
		args = append(args, assignedTo)
		idx++
	}
	if search != "" {
		query += ` AND (LOWER(title) LIKE $` + strconv.Itoa(idx) + ` OR LOWER(description) LIKE $` + strconv.Itoa(idx) + `)`
		searchTerm := "%" + search + "%"
		args = append(args, strings.ToLower(searchTerm))
		idx++
	}
	var total int
	if err := db.QueryRow(query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func getIssuesByProjectPaginatedFiltered(projectID, status, priority, assignedTo, search string, limit, offset int) ([]Issue, error) {
	query := `
	SELECT id, title, description, status, priority, project_id, created_by, assigned_to, created_at, updated_at
	FROM issues
	WHERE project_id = $1`
	args := []interface{}{projectID}
	idx := 2
	if status != "" {
		query += ` AND status = $` + strconv.Itoa(idx)
		args = append(args, status)
		idx++
	}
	if priority != "" {
		query += ` AND priority = $` + strconv.Itoa(idx)
		args = append(args, priority)
		idx++
	}
	if assignedTo != "" {
		query += ` AND assigned_to = $` + strconv.Itoa(idx)
		args = append(args, assignedTo)
		idx++
	}
	if search != "" {
		query += ` AND (LOWER(title) LIKE $` + strconv.Itoa(idx) + ` OR LOWER(description) LIKE $` + strconv.Itoa(idx) + `)`
		searchTerm := "%" + search + "%"
		args = append(args, strings.ToLower(searchTerm))
		idx++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(idx) + ` OFFSET $` + strconv.Itoa(idx+1)
	args = append(args, limit, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var issues []Issue
	for rows.Next() {
		var issue Issue
		err := rows.Scan(&issue.ID, &issue.Title, &issue.Description, &issue.Status, &issue.Priority, &issue.ProjectID, &issue.CreatedBy, &issue.AssignedTo, &issue.CreatedAt, &issue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func countIssuesByUserFiltered(userID, status, priority, assignedTo, search string) (int, error) {
	query := `SELECT COUNT(*) FROM issues WHERE created_by = $1`
	args := []interface{}{userID}
	idx := 2
	if status != "" {
		query += ` AND status = $` + strconv.Itoa(idx)
		args = append(args, status)
		idx++
	}
	if priority != "" {
		query += ` AND priority = $` + strconv.Itoa(idx)
		args = append(args, priority)
		idx++
	}
	if assignedTo != "" {
		query += ` AND assigned_to = $` + strconv.Itoa(idx)
		args = append(args, assignedTo)
		idx++
	}
	if search != "" {
		query += ` AND (LOWER(title) LIKE $` + strconv.Itoa(idx) + ` OR LOWER(description) LIKE $` + strconv.Itoa(idx) + `)`
		searchTerm := "%" + search + "%"
		args = append(args, strings.ToLower(searchTerm))
		idx++
	}
	var total int
	if err := db.QueryRow(query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func getIssuesByUserPaginatedFiltered(userID, status, priority, assignedTo, search string, limit, offset int) ([]Issue, error) {
	query := `
	SELECT id, title, description, project_id, created_by, created_at, updated_at
	FROM issues
	WHERE created_by = $1`
	args := []interface{}{userID}
	idx := 2
	if status != "" {
		query += ` AND status = $` + strconv.Itoa(idx)
		args = append(args, status)
		idx++
	}
	if priority != "" {
		query += ` AND priority = $` + strconv.Itoa(idx)
		args = append(args, priority)
		idx++
	}
	if assignedTo != "" {
		query += ` AND assigned_to = $` + strconv.Itoa(idx)
		args = append(args, assignedTo)
		idx++
	}
	if search != "" {
		query += ` AND (LOWER(title) LIKE $` + strconv.Itoa(idx) + ` OR LOWER(description) LIKE $` + strconv.Itoa(idx) + `)`
		searchTerm := "%" + search + "%"
		args = append(args, strings.ToLower(searchTerm))
		idx++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(idx) + ` OFFSET $` + strconv.Itoa(idx+1)
	args = append(args, limit, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var issues []Issue
	for rows.Next() {
		var issue Issue
		err := rows.Scan(&issue.ID, &issue.Title, &issue.Description, &issue.ProjectID, &issue.CreatedBy, &issue.CreatedAt, &issue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func createIssueHandler(c *gin.Context) {
	log.Println("createIssueHandler called")
	var issue Issue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue.ID = uuid.New().String()
	issue.CreatedBy = c.GetString("user_id")
	issue.CreatedAt = time.Now().Format(time.RFC3339)
	issue.UpdatedAt = issue.CreatedAt

	if err := createIssue(issue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create issue"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Issue created successfully"})
}

func createIssue(issue Issue) error {
	// Set default values if not provided
	if issue.Status == "" {
		issue.Status = "open"
	}
	if issue.Priority == "" {
		issue.Priority = "medium"
	}

	query := `
	INSERT INTO issues (id, title, description, status, priority, project_id, created_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.Exec(query, issue.ID, issue.Title, issue.Description, issue.Status, issue.Priority, issue.ProjectID, issue.CreatedBy, issue.CreatedAt, issue.UpdatedAt)
	if err != nil {
		log.Printf("Database error creating issue: %v", err)
	}
	return err
}

func getIssueHandler(c *gin.Context) {
	issueID := c.Param("id")
	issue, err := getIssueByID(issueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}
	c.JSON(http.StatusOK, issue)
}

func getIssueByID(issueID string) (Issue, error) {
	var issue Issue
	query := `
	SELECT id, title, description, status, priority, project_id, created_by, assigned_to, created_at, updated_at
	FROM issues
	WHERE id = $1
	`
	err := db.QueryRow(query, issueID).Scan(&issue.ID, &issue.Title, &issue.Description, &issue.Status, &issue.Priority, &issue.ProjectID, &issue.CreatedBy, &issue.AssignedTo, &issue.CreatedAt, &issue.UpdatedAt)
	return issue, err
}

func updateIssueHandler(c *gin.Context) {
	issueID := c.Param("id")
	issue, err := getIssueByID(issueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := updateIssue(issue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update issue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue updated successfully"})
}

func updateIssue(issue Issue) error {
	query := `
	UPDATE issues
	SET title = $1, description = $2, status = $3, priority = $4, updated_at = $5
	WHERE id = $6
	`
	_, err := db.Exec(query, issue.Title, issue.Description, issue.Status, issue.Priority, issue.UpdatedAt, issue.ID)
	return err
}

func deleteIssueHandler(c *gin.Context) {
	issueID := c.Param("id")
	if err := deleteIssue(issueID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete issue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Issue deleted successfully"})
}

func deleteIssue(issueID string) error {
	query := `
	DELETE FROM issues
	WHERE id = $1
	`
	_, err := db.Exec(query, issueID)
	return err
}

func createCommentHandler(c *gin.Context) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.ID = uuid.New().String()
	comment.CreatedBy = c.GetString("user_id")
	comment.CreatedAt = time.Now().Format(time.RFC3339)
	comment.UpdatedAt = comment.CreatedAt

	if err := createComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully"})
}

func createComment(comment Comment) error {
	query := `
	INSERT INTO comments (id, issue_id, content, created_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(query, comment.ID, comment.IssueID, comment.Content, comment.CreatedBy, comment.CreatedAt, comment.UpdatedAt)
	if err != nil {
		log.Printf("Database error creating comment: %v", err)
	}
	return err
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
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

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

func getCommentsHandler(c *gin.Context) {
	issueID := c.Param("issueId")
	limit := 10
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	total, err := countCommentsByIssue(issueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count comments"})
		return
	}
	comments, err := getCommentsByIssuePaginated(issueID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	if comments == nil {
		comments = []Comment{}
	}
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"message":  "Comments fetched successfully",
	})
}

func countCommentsByIssue(issueID string) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM comments WHERE issue_id = $1`
	err := db.QueryRow(query, issueID).Scan(&total)
	return total, err
}

func getCommentsByIssuePaginated(issueID string, limit, offset int) ([]Comment, error) {
	query := `
	SELECT id, issue_id, created_by, content, created_at, updated_at
	FROM comments
	WHERE issue_id = $1
	ORDER BY created_at ASC
	LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(query, issueID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.IssueID, &comment.CreatedBy, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func updateCommentHandler(c *gin.Context) {
	commentID := c.Param("id")
	userID := c.GetString("user_id")

	// Get the existing comment
	comment, err := getCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if the current user is the creator of the comment
	if comment.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own comments"})
		return
	}

	// Parse the update data
	var updateData struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the comment
	comment.Content = updateData.Content
	comment.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := updateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

func getCommentByID(commentID string) (Comment, error) {
	var comment Comment
	query := `
	SELECT id, issue_id, created_by, content, created_at, updated_at
	FROM comments
	WHERE id = $1
	`
	err := db.QueryRow(query, commentID).Scan(&comment.ID, &comment.IssueID, &comment.CreatedBy, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	return comment, err
}

func updateComment(comment Comment) error {
	query := `
	UPDATE comments
	SET content = $1, updated_at = $2
	WHERE id = $3
	`
	_, err := db.Exec(query, comment.Content, comment.UpdatedAt, comment.ID)
	if err != nil {
		log.Printf("Database error updating comment: %v", err)
	}
	return err
}

func deleteCommentHandler(c *gin.Context) {
	commentID := c.Param("id")
	userID := c.GetString("user_id")

	// Get the existing comment
	comment, err := getCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if the current user is the creator of the comment
	if comment.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
		return
	}

	// Delete the comment
	if err := deleteComment(commentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func deleteComment(commentID string) error {
	query := `
	DELETE FROM comments
	WHERE id = $1
	`
	_, err := db.Exec(query, commentID)
	if err != nil {
		log.Printf("Database error deleting comment: %v", err)
	}
	return err
}

func getUsersHandler(c *gin.Context) {
	users, err := getAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func getAllUsers() ([]User, error) {
	query := `
	SELECT id, email, first_name, last_name, role, created_at, updated_at
	FROM users
	ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getProfileHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func getUserByID(userID string) (User, error) {
	var user User
	query := `
	SELECT id, email, first_name, last_name, role, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

// Update profile handler
func updateProfileHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	var updateData struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := updateUserProfile(userID, updateData.FirstName, updateData.LastName, updateData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	user, err := getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUserProfile(userID, firstName, lastName, email string) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := db.Exec(query, firstName, lastName, email, userID)
	return err
}

// Admin: update user role
func updateUserRoleHandler(c *gin.Context) {
	userID := c.Param("id")
	var body struct {
		Role string `json:"role" binding:"required,oneof=admin user"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := setUserRole(userID, body.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}
	user, err := getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func setUserRole(userID, role string) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	_, err := db.Exec(query, role, userID)
	return err
}
