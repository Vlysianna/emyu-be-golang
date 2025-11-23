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

func GetUserOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	rows, err := database.DB.Query(`
		SELECT id, user_id, order_number, total_amount, shipping_cost, status, payment_method, shipping_address_id, created_at, updated_at
		FROM orders WHERE user_id = ? ORDER BY created_at DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.ShippingCost, &order.Status, &order.PaymentMethod, &order.ShippingAddressID, &order.CreatedAt, &order.UpdatedAt)
		orders = append(orders, order)
	}

	if orders == nil {
		orders = []models.Order{}
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	var order models.Order

	err := database.DB.QueryRow(`
		SELECT id, user_id, order_number, total_amount, shipping_cost, status, payment_method, shipping_address_id, created_at, updated_at
		FROM orders WHERE id = ?
	`, orderID).Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.ShippingCost, &order.Status, &order.PaymentMethod, &order.ShippingAddressID, &order.CreatedAt, &order.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	var req struct {
		TotalAmount       float64 `json:"total_amount" binding:"required"`
		ShippingCost      float64 `json:"shipping_cost" binding:"required"`
		PaymentMethod     string  `json:"payment_method" binding:"required,oneof=qris bank_transfer ewallet"`
		ShippingAddressID string  `json:"shipping_address_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	orderID := utils.GenerateID()
	orderNumber := utils.GenerateOrderNumber()

	_, err := database.DB.Exec(`
		INSERT INTO orders (id, user_id, order_number, total_amount, shipping_cost, status, payment_method, shipping_address_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, orderID, userID, orderNumber, req.TotalAmount, req.ShippingCost, "pending", req.PaymentMethod, req.ShippingAddressID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": orderID, "order_number": orderNumber, "message": "Order created"})
}

func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending paid packed shipped delivered canceled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", req.Status, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM orders WHERE id = ?", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
