package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/emyu/ecommer-be/utils"
)

func GetAllProducts(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, description, price, category_id, is_customizable, created_at, updated_at
		FROM products ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.IsCustomizable, &p.CreatedAt, &p.UpdatedAt); err != nil {
			continue
		}
		products = append(products, p)
	}

	if products == nil {
		products = []models.Product{}
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var p models.Product

	err := database.DB.QueryRow(`
		SELECT id, name, description, price, category_id, is_customizable, created_at, updated_at
		FROM products WHERE id = ?
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.IsCustomizable, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func CreateProduct(c *gin.Context) {
	var req struct {
		Name            string  `json:"name" binding:"required"`
		Description     string  `json:"description"`
		Price           float64 `json:"price" binding:"required"`
		CategoryID      string  `json:"category_id"`
		IsCustomizable  bool    `json:"is_customizable"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID := utils.GenerateID()
	_, err := database.DB.Exec(`
		INSERT INTO products (id, name, description, price, category_id, is_customizable)
		VALUES (?, ?, ?, ?, ?, ?)
	`, productID, req.Name, req.Description, req.Price, req.CategoryID, req.IsCustomizable)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": productID, "message": "Product created"})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name            string  `json:"name"`
		Description     string  `json:"description"`
		Price           float64 `json:"price"`
		CategoryID      string  `json:"category_id"`
		IsCustomizable  bool    `json:"is_customizable"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updates []string
	var args []interface{}

	if req.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, req.Name)
	}
	if req.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, req.Description)
	}
	if req.Price > 0 {
		updates = append(updates, "price = ?")
		args = append(args, req.Price)
	}
	if req.CategoryID != "" {
		updates = append(updates, "category_id = ?")
		args = append(args, req.CategoryID)
	}

	args = append(args, id)

	query := "UPDATE products SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	_, err := database.DB.Exec(query, args...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
