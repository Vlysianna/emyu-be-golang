package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/emyu/ecommer-be/middleware"
	"github.com/emyu/ecommer-be/utils"
)

func GetReviewsByProduct(c *gin.Context) {
	productID := c.Param("productId")
	rows, err := database.DB.Query(`
		SELECT id, user_id, product_id, rating, comment, created_at
		FROM reviews WHERE product_id = ? ORDER BY created_at DESC
	`, productID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		rows.Scan(&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt)
		reviews = append(reviews, review)
	}

	if reviews == nil {
		reviews = []models.Review{}
	}

	c.JSON(http.StatusOK, reviews)
}

func GetUserReviews(c *gin.Context) {
	userID := middleware.GetUserID(c)
	rows, err := database.DB.Query(`
		SELECT id, user_id, product_id, rating, comment, created_at
		FROM reviews WHERE user_id = ? ORDER BY created_at DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		rows.Scan(&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt)
		reviews = append(reviews, review)
	}

	if reviews == nil {
		reviews = []models.Review{}
	}

	c.JSON(http.StatusOK, reviews)
}

func CreateReview(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id" binding:"required"`
		Rating    int    `json:"rating" binding:"required,min=1,max=5"`
		Comment   string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	reviewID := utils.GenerateID()

	_, err := database.DB.Exec(`
		INSERT INTO reviews (id, user_id, product_id, rating, comment, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, reviewID, userID, req.ProductID, req.Rating, req.Comment, time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": reviewID, "message": "Review created"})
}

func UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")
	var req struct {
		Rating  int    `json:"rating" binding:"min=1,max=5"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)

	// Check if review belongs to user
	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM reviews WHERE id = ?", reviewID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE reviews SET rating = ?, comment = ? WHERE id = ?
	`, req.Rating, req.Comment, reviewID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated"})
}

func DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	userID := middleware.GetUserID(c)

	// Check if review belongs to user
	var ownerID string
	err := database.DB.QueryRow("SELECT user_id FROM reviews WHERE id = ?", reviewID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM reviews WHERE id = ?", reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted"})
}
