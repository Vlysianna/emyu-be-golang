package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/emyu/ecommer-be/utils"
	"github.com/gin-gonic/gin"
)

func GetPayments(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, order_id, payment_status, payment_code, paid_at, created_at, updated_at
		FROM payments ORDER BY created_at DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		rows.Scan(&payment.ID, &payment.OrderID, &payment.PaymentStatus, &payment.PaymentCode, &payment.PaidAt, &payment.CreatedAt, &payment.UpdatedAt)
		payments = append(payments, payment)
	}

	if payments == nil {
		payments = []models.Payment{}
	}

	c.JSON(http.StatusOK, payments)
}

func GetPaymentByID(c *gin.Context) {
	paymentID := c.Param("id")
	var payment models.Payment

	err := database.DB.QueryRow(`
		SELECT id, order_id, payment_status, payment_code, paid_at, created_at, updated_at
		FROM payments WHERE id = ?
	`, paymentID).Scan(&payment.ID, &payment.OrderID, &payment.PaymentStatus, &payment.PaymentCode, &payment.PaidAt, &payment.CreatedAt, &payment.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func CreatePayment(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentID := utils.GenerateID()
	paymentCode := utils.GeneratePaymentCode()

	_, err := database.DB.Exec(`
		INSERT INTO payments (id, order_id, payment_status, payment_code)
		VALUES (?, ?, ?, ?)
	`, paymentID, req.OrderID, "pending", paymentCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": paymentID, "payment_code": paymentCode, "message": "Payment created"})
}

func UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending success failed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var paidAt *time.Time
	if req.Status == "success" {
		now := time.Now()
		paidAt = &now
	}

	_, err := database.DB.Exec(
		"UPDATE payments SET payment_status = ?, paid_at = ? WHERE id = ?",
		req.Status, paidAt, paymentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated"})
}
