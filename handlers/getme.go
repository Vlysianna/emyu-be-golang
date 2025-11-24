package handlers

import (
	_"fmt"
	"database/sql"
	"net/http"
	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/gin-gonic/gin"
)
func GetMyProfile(c *gin.Context) {
    userID := c.GetString("userID")

    var user models.User

    err := database.DB.QueryRow(`
        SELECT id, name, email, phone, role_id, is_active, created_at, updated_at
        FROM users WHERE id = ?
    `, userID).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Phone,
        &user.RoleID,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if err != nil {
        c.JSON(500, gin.H{
            "error": "SCAN FAILED",
            "details": err.Error(),
            "userID": userID,
        })
        return
    }

    c.JSON(http.StatusOK, user)
}