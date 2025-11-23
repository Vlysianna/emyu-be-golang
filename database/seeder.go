package database

import (
	"fmt"
	"log"

	"github.com/emyu/ecommer-be/utils"
)

func SeedDatabase() error {
	log.Println("🌱 Seeding database...")

	// Seed Roles (already in schema.sql, but verify)
	err := seedRoles()
	if err != nil {
		return err
	}

	// Seed Admin User
	err = seedAdminUser()
	if err != nil {
		return err
	}

	// Seed Regular Users
	err = seedRegularUsers()
	if err != nil {
		return err
	}

	// Seed Categories
	err = seedCategories()
	if err != nil {
		return err
	}

	// Seed Products
	err = seedProducts()
	if err != nil {
		return err
	}

	// Seed Product Variants
	err = seedProductVariants()
	if err != nil {
		return err
	}

	// Seed Product Images
	err = seedProductImages()
	if err != nil {
		return err
	}

	log.Println("✅ Database seeding completed successfully!")
	return nil
}

func seedRoles() error {
	log.Println("  - Seeding roles...")

	roles := []map[string]interface{}{
		{
			"id":          1,
			"name":        "admin",
			"description": "Administrator with full access",
			"is_active":   true,
		},
		{
			"id":          2,
			"name":        "user",
			"description": "Regular customer user",
			"is_active":   true,
		},
		{
			"id":          3,
			"name":        "seller",
			"description": "Seller with product management access",
			"is_active":   true,
		},
	}

	for _, role := range roles {
		_, err := DB.Exec(
			"INSERT IGNORE INTO roles (id, name, description, is_active) VALUES (?, ?, ?, ?)",
			role["id"], role["name"], role["description"], role["is_active"],
		)
		if err != nil {
			return fmt.Errorf("failed to seed role %s: %w", role["name"], err)
		}
	}

	return nil
}

func seedAdminUser() error {
	log.Println("  - Seeding admin user...")

	adminEmail := "admin@emyu.com"

	// Check if admin already exists
	var existingID string
	err := DB.QueryRow("SELECT id FROM users WHERE email = ?", adminEmail).Scan(&existingID)
	if err == nil {
		log.Println("    ℹ️  Admin user already exists")
		return nil
	}

	adminID := utils.GenerateID()
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	_, err = DB.Exec(
		"INSERT INTO users (id, name, email, phone, role_id, password, is_active) VALUES (?, ?, ?, ?, ?, ?, ?)",
		adminID, "Admin User", adminEmail, "08111111111", 1, hashedPassword, true,
	)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	log.Println("    ✓ Admin user created (email: admin@emyu.com, password: admin123)")
	return nil
}

func seedRegularUsers() error {
	log.Println("  - Seeding regular users...")

	users := []map[string]interface{}{
		{
			"name":  "John Doe",
			"email": "john@example.com",
			"phone": "08123456789",
		},
		{
			"name":  "Jane Smith",
			"email": "jane@example.com",
			"phone": "08987654321",
		},
		{
			"name":  "Budi Santoso",
			"email": "budi@example.com",
			"phone": "08555666777",
		},
		{
			"name":  "Siti Nur Azizah",
			"email": "siti@example.com",
			"phone": "08444555666",
		},
	}

	for _, user := range users {
		// Check if user exists
		var existingID string
		err := DB.QueryRow("SELECT id FROM users WHERE email = ?", user["email"]).Scan(&existingID)
		if err == nil {
			continue // User already exists
		}

		userID := utils.GenerateID()
		hashedPassword, err := utils.HashPassword("password123")
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", user["email"], err)
		}

		_, err = DB.Exec(
			"INSERT INTO users (id, name, email, phone, role_id, password, is_active) VALUES (?, ?, ?, ?, ?, ?, ?)",
			userID, user["name"], user["email"], user["phone"], 2, hashedPassword, true,
		)
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user["email"], err)
		}
	}

	log.Println("    ✓ Regular users created")
	return nil
}

func seedCategories() error {
	log.Println("  - Seeding categories...")

	categories := []map[string]interface{}{
		{
			"name":        "Sports Jersey",
			"description": "Custom sports jerseys for teams and athletes",
		},
		{
			"name":        "Casual Wear",
			"description": "Comfortable casual clothing with custom options",
		},
		{
			"name":        "Uniforms",
			"description": "Professional uniforms for organizations",
		},
		{
			"name":        "Merchandise",
			"description": "Custom branded merchandise and accessories",
		},
	}

	for _, cat := range categories {
		// Check if category exists
		var existingID string
		err := DB.QueryRow("SELECT id FROM categories WHERE name = ?", cat["name"]).Scan(&existingID)
		if err == nil {
			continue // Category already exists
		}

		catID := utils.GenerateID()
		_, err = DB.Exec(
			"INSERT INTO categories (id, name, description) VALUES (?, ?, ?)",
			catID, cat["name"], cat["description"],
		)
		if err != nil {
			return fmt.Errorf("failed to seed category %s: %w", cat["name"], err)
		}
	}

	log.Println("    ✓ Categories created")
	return nil
}

func seedProducts() error {
	log.Println("  - Seeding products...")

	// Get first category ID for seeding
	var categoryID string
	err := DB.QueryRow("SELECT id FROM categories LIMIT 1").Scan(&categoryID)
	if err != nil {
		return fmt.Errorf("no categories found for seeding products: %w", err)
	}

	products := []map[string]interface{}{
		{
			"name":             "Premium Sports Jersey",
			"description":      "High-quality polyester jersey with custom name and number printing",
			"price":            299000,
			"is_customizable":  true,
		},
		{
			"name":             "Classic Cotton T-Shirt",
			"description":      "100% cotton t-shirt perfect for casual wear and custom printing",
			"price":            149000,
			"is_customizable":  true,
		},
		{
			"name":             "Team Uniform Polo",
			"description":      "Professional polo shirt for team uniforms",
			"price":            249000,
			"is_customizable":  true,
		},
		{
			"name":             "Basic White T-Shirt",
			"description":      "Simple white t-shirt, great for custom designs",
			"price":            99000,
			"is_customizable":  true,
		},
		{
			"name":             "Hoodie Sweater",
			"description":      "Comfortable hoodie with custom logo option",
			"price":            349000,
			"is_customizable":  true,
		},
	}

	for _, prod := range products {
		// Check if product exists
		var existingID string
		err := DB.QueryRow("SELECT id FROM products WHERE name = ?", prod["name"]).Scan(&existingID)
		if err == nil {
			continue // Product already exists
		}

		prodID := utils.GenerateID()
		_, err = DB.Exec(
			"INSERT INTO products (id, name, description, price, category_id, is_customizable) VALUES (?, ?, ?, ?, ?, ?)",
			prodID, prod["name"], prod["description"], prod["price"], categoryID, prod["is_customizable"],
		)
		if err != nil {
			return fmt.Errorf("failed to seed product %s: %w", prod["name"], err)
		}
	}

	log.Println("    ✓ Products created")
	return nil
}

func seedProductVariants() error {
	log.Println("  - Seeding product variants...")

	// Get all products
	rows, err := DB.Query("SELECT id FROM products")
	if err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	variants := []map[string]interface{}{
		{"name": "XS", "price_adjustment": 0},
		{"name": "S", "price_adjustment": 0},
		{"name": "M", "price_adjustment": 0},
		{"name": "L", "price_adjustment": 0},
		{"name": "XL", "price_adjustment": 25000},
		{"name": "XXL", "price_adjustment": 50000},
	}

	for rows.Next() {
		var productID string
		if err := rows.Scan(&productID); err != nil {
			continue
		}

		for _, variant := range variants {
			variantID := utils.GenerateID()
			_, err = DB.Exec(
				"INSERT INTO product_variants (id, product_id, name, price_adjustment) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id",
				variantID, productID, variant["name"], variant["price_adjustment"],
			)
			if err != nil {
				return fmt.Errorf("failed to seed variant: %w", err)
			}
		}
	}

	log.Println("    ✓ Product variants created")
	return nil
}

func seedProductImages() error {
	log.Println("  - Seeding product images...")

	// Get all products
	rows, err := DB.Query("SELECT id FROM products")
	if err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	imageURLs := []string{
		"https://via.placeholder.com/400x400?text=Product+1",
		"https://via.placeholder.com/400x400?text=Product+2",
		"https://via.placeholder.com/400x400?text=Product+3",
	}

	for rows.Next() {
		var productID string
		if err := rows.Scan(&productID); err != nil {
			continue
		}

		for i, imageURL := range imageURLs {
			imageID := utils.GenerateID()
			_, err = DB.Exec(
				"INSERT INTO product_images (id, product_id, image_url) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id=id",
				imageID, productID, imageURL,
			)
			if err != nil {
				return fmt.Errorf("failed to seed product image: %w", err)
			}

			// Only add 1 image per product for seeding
			if i == 0 {
				break
			}
		}
	}

	log.Println("    ✓ Product images created")
	return nil
}
