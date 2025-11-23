package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/emyu/ecommer-be/database"
	"github.com/emyu/ecommer-be/models"
	"github.com/emyu/ecommer-be/utils"
)

func GetAllCategories(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, description, created_at, updated_at FROM categories ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category

	err := database.DB.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM categories WHERE id = ?",
		id,
	).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	catID := utils.GenerateID()
	_, err := database.DB.Exec(
		"INSERT INTO categories (id, name, description) VALUES (?, ?, ?)",
		catID, req.Name, req.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": catID, "message": "Category created"})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(
		"UPDATE categories SET name = ?, description = ? WHERE id = ?",
		req.Name, req.Description, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
