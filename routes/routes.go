package routes

import (
	"github.com/emyu/ecommer-be/handlers"
	"github.com/emyu/ecommer-be/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Public routes - Auth
	auth := router.Group("/api")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Public routes - Products & Categories
	public := router.Group("/api")
	{
		public.GET("/products", handlers.GetAllProducts)
		public.GET("/products/:id", handlers.GetProductByID)
		public.GET("/categories", handlers.GetAllCategories)
		public.GET("/categories/:id", handlers.GetCategoryByID)
	}

	// Protected routes - Auth only
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", handlers.GetMyProfile)
		protected.POST("/logout", handlers.Logout)

		// Cart
		protected.GET("/carts", handlers.GetUserCart)
		protected.POST("/carts", handlers.CreateCart)
		protected.POST("/cart-items", handlers.AddToCart)
		protected.PUT("/cart-items/:itemId", handlers.UpdateCartItem)
		protected.DELETE("/cart-items/:itemId", handlers.RemoveFromCart)

		// Orders
		protected.GET("/orders", handlers.GetUserOrders)
		protected.GET("/orders/:id", handlers.GetOrderByID)
		protected.POST("/orders", handlers.CreateOrder)

		// Payments
		protected.GET("/payments", handlers.GetPayments)
		protected.POST("/payments", handlers.CreatePayment)
		protected.GET("/payments/:id", handlers.GetPaymentByID)
		protected.PUT("/payments/:id", handlers.UpdatePaymentStatus)

		// Reviews
		protected.GET("/reviews/products/:productId", handlers.GetReviewsByProduct)
		protected.GET("/reviews/user", handlers.GetUserReviews)
		protected.POST("/reviews", handlers.CreateReview)
		protected.PUT("/reviews/:id", handlers.UpdateReview)
		protected.DELETE("/reviews/:id", handlers.DeleteReview)

		// Shipping Addresses
		protected.GET("/shipping-addresses", handlers.GetUserShippingAddresses)
		protected.GET("/shipping-addresses/:id", handlers.GetShippingAddressByID)
		protected.POST("/shipping-addresses", handlers.CreateShippingAddress)
		protected.PUT("/shipping-addresses/:id", handlers.UpdateShippingAddress)
		protected.DELETE("/shipping-addresses/:id", handlers.DeleteShippingAddress)
	}

	// Admin routes
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Products
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)

		// Categories
		admin.POST("/categories", handlers.CreateCategory)
		admin.PUT("/categories/:id", handlers.UpdateCategory)
		admin.DELETE("/categories/:id", handlers.DeleteCategory)

		// Order management
		admin.GET("/orders", handlers.GetAllOrders)
		admin.PUT("/orders/:id", handlers.UpdateOrderStatus)
		admin.DELETE("/orders/:id", handlers.DeleteOrder)

		// User management
		admin.GET("/users", handlers.GetAllUsers)
		admin.GET("/users/:id", handlers.GetUserByID)
		admin.PUT("/users/:id", handlers.UpdateUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)
		admin.GET("/users/:id/stats", handlers.GetUserStats)
	}
}
