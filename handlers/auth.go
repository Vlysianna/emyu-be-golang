package handlers

import (
	"database/sql"
	"net/http"

	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/emyu/ecommer-be/utils"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email exists
	var existingUser models.User
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	userID := utils.GenerateID()
	// Role ID 2 is default "user" role (set in schema)
	user := models.User{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   2, // Default user role
		Password: hashedPassword,
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (id, name, email, phone, role_id, password) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Name, user.Email, user.Phone, user.RoleID, user.Password,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Get role info for token
	var roleID int
	var roleName string
	var permissions []string
	err = database.DB.QueryRow(
		"SELECT id, name FROM roles WHERE id = ?",
		user.RoleID,
	).Scan(&roleID, &roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user role"})
		return
	}

	// Fetch role permissions
	permRows, err := database.DB.Query(
		"SELECT permission FROM role_permissions WHERE role_id = ?",
		roleID,
	)
	if err == nil {
		defer permRows.Close()
		for permRows.Next() {
			var perm string
			if err := permRows.Scan(&perm); err == nil {
				permissions = append(permissions, perm)
			}
		}
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, roleID, roleName, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, name, email, phone, role_id, password FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.RoleID, &user.Password)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Get role info for token
	var roleID int
	var roleName string
	var permissions []string
	err = database.DB.QueryRow(
		"SELECT id, name FROM roles WHERE id = ?",
		user.RoleID,
	).Scan(&roleID, &roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user role"})
		return
	}

	// Fetch role permissions
	permRows, err := database.DB.Query(
		"SELECT permission FROM role_permissions WHERE role_id = ?",
		roleID,
	)
	if err == nil {
		defer permRows.Close()
		for permRows.Next() {
			var perm string
			if err := permRows.Scan(&perm); err == nil {
				permissions = append(permissions, perm)
			}
		}
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, roleID, roleName, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

func Logout(c *gin.Context) {
	// JWT is stateless, logout is handled by frontend (delete token)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
