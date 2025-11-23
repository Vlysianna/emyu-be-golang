package handlers

import (
	"database/sql"
	"net/http"

	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/middleware"
	"github.com/emyu/ecommer-be/models"
	"github.com/emyu/ecommer-be/utils"
	"github.com/gin-gonic/gin"
)

func GetUserShippingAddresses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	rows, err := database.DB.Query(`
		SELECT id, user_id, address, city, province, postal_code, phone, created_at, updated_at
		FROM shipping_addresses WHERE user_id = ? ORDER BY created_at DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shipping addresses"})
		return
	}
	defer rows.Close()

	var addresses []models.ShippingAddress
	for rows.Next() {
		var addr models.ShippingAddress
		rows.Scan(&addr.ID, &addr.UserID, &addr.Address, &addr.City, &addr.Province, &addr.PostalCode, &addr.Phone, &addr.CreatedAt, &addr.UpdatedAt)
		addresses = append(addresses, addr)
	}

	if addresses == nil {
		addresses = []models.ShippingAddress{}
	}

	c.JSON(http.StatusOK, addresses)
}

func GetShippingAddressByID(c *gin.Context) {
	addressID := c.Param("id")
	userID := middleware.GetUserID(c)
	var addr models.ShippingAddress

	err := database.DB.QueryRow(`
		SELECT id, user_id, address, city, province, postal_code, phone, created_at, updated_at
		FROM shipping_addresses WHERE id = ? AND user_id = ?
	`, addressID, userID).Scan(&addr.ID, &addr.UserID, &addr.Address, &addr.City, &addr.Province, &addr.PostalCode, &addr.Phone, &addr.CreatedAt, &addr.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping address not found"})
		return
	}

	c.JSON(http.StatusOK, addr)
}

func CreateShippingAddress(c *gin.Context) {
	var req struct {
		Address    string `json:"address" binding:"required"`
		City       string `json:"city" binding:"required"`
		Province   string `json:"province" binding:"required"`
		PostalCode string `json:"postal_code" binding:"required"`
		Phone      string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	addressID := utils.GenerateID()

	_, err := database.DB.Exec(`
		INSERT INTO shipping_addresses (id, user_id, address, city, province, postal_code, phone)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, addressID, userID, req.Address, req.City, req.Province, req.PostalCode, req.Phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shipping address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": addressID, "message": "Shipping address created"})
}

func UpdateShippingAddress(c *gin.Context) {
	addressID := c.Param("id")
	var req struct {
		Address    string `json:"address"`
		City       string `json:"city"`
		Province   string `json:"province"`
		PostalCode string `json:"postal_code"`
		Phone      string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)

	// Check if address belongs to user
	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM shipping_addresses WHERE id = ?", addressID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping address not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE shipping_addresses SET address = ?, city = ?, province = ?, postal_code = ?, phone = ? WHERE id = ?
	`, req.Address, req.City, req.Province, req.PostalCode, req.Phone, addressID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipping address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping address updated"})
}

func DeleteShippingAddress(c *gin.Context) {
	addressID := c.Param("id")
	userID := middleware.GetUserID(c)

	// Check if address belongs to user
	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM shipping_addresses WHERE id = ?", addressID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping address not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM shipping_addresses WHERE id = ?", addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shipping address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping address deleted"})
}
