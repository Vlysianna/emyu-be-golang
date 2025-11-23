package handlers

import (
	"database/sql"
	"net/http"

	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, email, phone, role, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}

	if users == nil {
		users = []models.User{}
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	err := database.DB.QueryRow(`
		SELECT id, name, email, phone, role, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Role  string `json:"role" binding:"oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(`
		UPDATE users SET name = ?, email = ?, phone = ?, role = ? WHERE id = ?
	`, req.Name, req.Email, req.Phone, req.Role, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func GetUserStats(c *gin.Context) {
	userID := c.Param("id")

	stats := struct {
		TotalOrders  int64   `json:"total_orders"`
		TotalSpent   float64 `json:"total_spent"`
		TotalReviews int64   `json:"total_reviews"`
	}{}

	// Get total orders
	database.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = ?", userID).Scan(&stats.TotalOrders)

	// Get total spent
	database.DB.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE user_id = ?", userID).Scan(&stats.TotalSpent)

	// Get total reviews
	database.DB.QueryRow("SELECT COUNT(*) FROM reviews WHERE user_id = ?", userID).Scan(&stats.TotalReviews)

	c.JSON(http.StatusOK, stats)
}
