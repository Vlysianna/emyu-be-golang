package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/emyu/ecommer-be/middleware"
	"github.com/emyu/ecommer-be/utils"
)

func GetUserCart(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var cart models.Cart

	err := database.DB.QueryRow(
		"SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = ?",
		userID,
	).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Get cart items
	rows, _ := database.DB.Query(
		"SELECT id, cart_id, product_variant_id, quantity, custom_name, custom_number, created_at, updated_at FROM cart_items WHERE cart_id = ?",
		cart.ID,
	)
	defer rows.Close()

	for rows.Next() {
		var item models.CartItem
		rows.Scan(&item.ID, &item.CartID, &item.ProductVariantID, &item.Quantity, &item.CustomName, &item.CustomNumber, &item.CreatedAt, &item.UpdatedAt)
		cart.Items = append(cart.Items, item)
	}

	c.JSON(http.StatusOK, cart)
}

func CreateCart(c *gin.Context) {
	userID := middleware.GetUserID(c)
	cartID := utils.GenerateID()

	_, err := database.DB.Exec(
		"INSERT INTO carts (id, user_id) VALUES (?, ?)",
		cartID, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": cartID, "message": "Cart created"})
}

func AddToCart(c *gin.Context) {
	var req struct {
		CartID           string `json:"cart_id" binding:"required"`
		ProductVariantID string `json:"product_variant_id" binding:"required"`
		Quantity         int    `json:"quantity" binding:"required,min=1"`
		CustomName       string `json:"custom_name"`
		CustomNumber     string `json:"custom_number"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemID := utils.GenerateID()
	_, err := database.DB.Exec(`
		INSERT INTO cart_items (id, cart_id, product_variant_id, quantity, custom_name, custom_number)
		VALUES (?, ?, ?, ?, ?, ?)
	`, itemID, req.CartID, req.ProductVariantID, req.Quantity, req.CustomName, req.CustomNumber)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": itemID, "message": "Item added to cart"})
}

func RemoveFromCart(c *gin.Context) {
	itemID := c.Param("itemId")
	_, err := database.DB.Exec("DELETE FROM cart_items WHERE id = ?", itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

func UpdateCartItem(c *gin.Context) {
	itemID := c.Param("itemId")
	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec("UPDATE cart_items SET quantity = ? WHERE id = ?", req.Quantity, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated"})
}
