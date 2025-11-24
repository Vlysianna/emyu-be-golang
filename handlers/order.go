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
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.ShippingCost, &order.Status, &order.PaymentMethod, &order.ShippingAddressID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}

		// Fetch order items
		order.Items, _ = getOrderItems(order.ID)

		// Fetch shipping address
		order.ShippingAddress, _ = getShippingAddressDetails(order.ShippingAddressID)

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading orders"})
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	c.JSON(http.StatusOK, orders)
}

// GetAllOrders - Admin endpoint to get all orders
func GetAllOrders(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, order_number, total_amount, shipping_cost, status, payment_method, shipping_address_id, created_at, updated_at
		FROM orders ORDER BY created_at DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.ShippingCost, &order.Status, &order.PaymentMethod, &order.ShippingAddressID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}

		// Fetch order items
		order.Items, _ = getOrderItems(order.ID)

		// Fetch shipping address
		order.ShippingAddress, _ = getShippingAddressDetails(order.ShippingAddressID)

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading orders"})
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	c.JSON(http.StatusOK, orders)
}

// Helper function to get order items with product details
func getOrderItems(orderID string) ([]models.OrderItem, error) {
	rows, err := database.DB.Query(`
		SELECT oi.id, oi.order_id, oi.product_variant_id, oi.quantity, oi.price,
		       pv.id, pv.product_id, pv.name,
		       p.id, p.name, p.price
		FROM order_items oi
		LEFT JOIN product_variants pv ON oi.product_variant_id = pv.id
		LEFT JOIN products p ON pv.product_id = p.id
		WHERE oi.order_id = ?
	`, orderID)
	if err != nil {
		return []models.OrderItem{}, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var variantID, variantName, productID, productName sql.NullString
		var productPrice sql.NullFloat64
		var variantProductID sql.NullString

		rows.Scan(&item.ID, &item.OrderID, &item.ProductVariantID, &item.Quantity, &item.Price,
			&variantID, &variantProductID, &variantName,
			&productID, &productName, &productPrice)

		// Build ProductVariant with product info - but we need the product name in the item
		// For now, just return the items with the variant
		if variantID.Valid {
			item.ProductVariant = &models.ProductVariant{
				ID:        variantID.String,
				ProductID: variantProductID.String,
				Name:      variantName.String,
			}
		}

		items = append(items, item)
	}

	return items, nil
}

// Helper function to get shipping address details
func getShippingAddressDetails(addressID string) (*models.ShippingAddress, error) {
	if addressID == "" {
		return nil, nil
	}

	var address models.ShippingAddress
	err := database.DB.QueryRow(`
		SELECT id, user_id, address, city, province, postal_code, phone, created_at, updated_at
		FROM shipping_addresses WHERE id = ?
	`, addressID).Scan(&address.ID, &address.UserID, &address.Address, &address.City, &address.Province, &address.PostalCode, &address.Phone, &address.CreatedAt, &address.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &address, nil
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
		PaymentMethod     string  `json:"payment_method" binding:"required,oneof=qris bank_transfer ewallet credit_card e_wallet"`
		ShippingAddressID string  `json:"shipping_address_id" binding:"required"`
		Items             []struct {
			ProductVariantID string  `json:"product_variant_id" binding:"required"`
			Quantity         int     `json:"quantity" binding:"required,min=1"`
			Price            float64 `json:"price" binding:"required"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	orderID := utils.GenerateID()
	orderNumber := utils.GenerateOrderNumber()

	// Begin transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Create order
	_, err = tx.Exec(`
		INSERT INTO orders (id, user_id, order_number, total_amount, shipping_cost, status, payment_method, shipping_address_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, orderID, userID, orderNumber, req.TotalAmount, req.ShippingCost, "pending", req.PaymentMethod, req.ShippingAddressID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Add order items
	for _, item := range req.Items {
		itemID := utils.GenerateID()
		_, err := tx.Exec(`
			INSERT INTO order_items (id, order_id, product_variant_id, quantity, price)
			VALUES (?, ?, ?, ?, ?)
		`, itemID, orderID, item.ProductVariantID, item.Quantity, item.Price)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order items"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": orderID, "order_number": orderNumber, "message": "Order created successfully"})
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
