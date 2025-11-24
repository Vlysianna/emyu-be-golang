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
		SELECT id, name, email, phone, role_id, created_at, updated_at
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
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading users"})
		return
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
		SELECT id, name, email, phone, role_id, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		RoleID int    `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(`
		UPDATE users SET name = ?, email = ?, phone = ?, role_id = ? WHERE id = ?
	`, req.Name, req.Email, req.Phone, req.RoleID, userID)

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
	err := database.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = ?", userID).Scan(&stats.TotalOrders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order stats"})
		return
	}

	// Get total spent
	err = database.DB.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE user_id = ?", userID).Scan(&stats.TotalSpent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spending stats"})
		return
	}

	// Get total reviews
	err = database.DB.QueryRow("SELECT COUNT(*) FROM reviews WHERE user_id = ?", userID).Scan(&stats.TotalReviews)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
