# Emyu E-Commerce API (Golang)

A simple, clean, and fast RESTful API for e-commerce platform built with **Golang**, **Gin Framework**, and **MySQL**.

## 🚀 Features

- **Authentication**: JWT-based Bearer Token (24-hour expiry)
- **Role-based Access**: Admin vs User permissions
- **Product Management**: CRUD operations for products & categories
- **Shopping Cart**: Add/update/remove items from cart
- **Orders**: Create and manage orders with multiple payment methods
- **Payments**: Track payment status (pending, success, failed)
- **Product Variants**: Support for product variants (sizes, colors, etc.)
- **Customizable Products**: Support for custom names and numbers
- **Shipping Addresses**: Multiple addresses per user
- **Clean Architecture**: Well-organized code structure
- **Fast & Lightweight**: Built with Gin framework for high performance

## 📋 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: MySQL 5.7+
- **Authentication**: JWT (golang-jwt)
- **Password Hashing**: bcrypt
- **Environment**: godotenv

## 🛠️ Installation

### Prerequisites
- Go 1.21 or higher
- MySQL 5.7 or higher
- Git

### Step 1: Clone & Setup

```bash
git clone <repo-url>
cd emyu-be-golang
go mod download
```

### Step 2: Database Setup

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE emyu CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Run migrations
mysql -u root -p emyu < database/schema.sql
```

### Step 3: Environment Configuration

```bash
cp .env.example .env
```

Edit `.env` file:
```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=emyu
SERVER_PORT=8080
SERVER_ENV=development
JWT_SECRET=your-super-secret-key-change-this
```

### Step 4: Database Seeding (Optional)

**Populate database with sample data for development:**

```bash
# Option 1: Using Makefile (Recommended)
make db-fresh-seed    # Fresh database + seed data
make db-seed          # Seed existing database

# Option 2: Using binary directly
go run cmd/api/main.go --seed

# Option 3: Using compiled binary
./bin/api --seed
```

**Seeded data includes:**
- 3 Roles: admin (id: 1), user (id: 2), seller (id: 3)
- 1 Admin User: `admin@emyu.com` / `admin123`
- 4 Regular Users: john, jane, budi, siti (all with password: `password123`)
- 4 Product Categories: Sports Jersey, Casual Wear, Uniforms, Merchandise
- 5 Sample Products: each with 6 size variants (XS-XXL)
- Placeholder product images

### Step 5: Run Server

```bash
go run cmd/api/main.go
```

Server will run on `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication
All protected endpoints require Bearer token in header:
```
Authorization: Bearer <your-token>
```

---

## 🔐 Auth Endpoints

### Register
```
POST /api/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "08123456789"
}

Response: 201 Created
{
  "user": {
    "id": "abc123",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "08123456789",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Login
```
POST /api/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}

Response: 200 OK
{
  "user": {...},
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Logout
```
POST /api/logout
Authorization: Bearer <token>

Response: 200 OK
{
  "message": "Logged out successfully"
}
```

---

## 📦 Products & Categories (Public)

### Get All Products
```
GET /api/products

Response: 200 OK
[
  {
    "id": "prod123",
    "name": "Kaos Putih",
    "description": "Kaos putih premium",
    "price": 89000,
    "category_id": "cat123",
    "is_customizable": true,
    "created_at": "2025-11-23T10:00:00Z"
  },
  ...
]
```

### Get Product By ID
```
GET /api/products/:id

Response: 200 OK
{
  "id": "prod123",
  "name": "Kaos Putih",
  "description": "Kaos putih premium",
  "price": 89000,
  "category_id": "cat123",
  "is_customizable": true,
  "created_at": "2025-11-23T10:00:00Z"
}
```

### Get All Categories
```
GET /api/categories

Response: 200 OK
[
  {
    "id": "cat123",
    "name": "Apparel",
    "description": "Clothing items",
    "created_at": "2025-11-23T10:00:00Z"
  },
  ...
]
```

---

## 🛒 Shopping Cart (Protected)

### Get User Cart
```
GET /api/carts
Authorization: Bearer <token>

Response: 200 OK
{
  "id": "cart123",
  "user_id": "user123",
  "items": [
    {
      "id": "item1",
      "cart_id": "cart123",
      "product_variant_id": "var123",
      "quantity": 2,
      "custom_name": "John",
      "custom_number": "1",
      "created_at": "2025-11-23T10:00:00Z"
    }
  ]
}
```

### Create Cart
```
POST /api/carts
Authorization: Bearer <token>

Response: 201 Created
{
  "id": "cart123",
  "message": "Cart created"
}
```

### Add to Cart
```
POST /api/cart-items
Authorization: Bearer <token>
Content-Type: application/json

{
  "cart_id": "cart123",
  "product_variant_id": "var123",
  "quantity": 2,
  "custom_name": "John",
  "custom_number": "1"
}

Response: 201 Created
{
  "id": "item1",
  "message": "Item added to cart"
}
```

### Update Cart Item
```
PUT /api/cart-items/:itemId
Authorization: Bearer <token>
Content-Type: application/json

{
  "quantity": 5
}

Response: 200 OK
{
  "message": "Item updated"
}
```

### Remove from Cart
```
DELETE /api/cart-items/:itemId
Authorization: Bearer <token>

Response: 200 OK
{
  "message": "Item removed from cart"
}
```

---

## 📦 Orders (Protected)

### Get User Orders
```
GET /api/orders
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": "ord123",
    "user_id": "user123",
    "order_number": "ORD-20251123-1234",
    "total_amount": 200000,
    "shipping_cost": 25000,
    "status": "pending",
    "payment_method": "qris",
    "shipping_address_id": "addr123",
    "created_at": "2025-11-23T10:00:00Z"
  }
]
```

### Create Order
```
POST /api/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "total_amount": 200000,
  "shipping_cost": 25000,
  "payment_method": "qris",
  "shipping_address_id": "addr123"
}

Response: 201 Created
{
  "id": "ord123",
  "order_number": "ORD-20251123-1234",
  "message": "Order created"
}
```

### Get Order By ID
```
GET /api/orders/:id
Authorization: Bearer <token>

Response: 200 OK
{
  "id": "ord123",
  ...
}
```

---

## 💳 Payments (Protected)

### Create Payment
```
POST /api/payments
Authorization: Bearer <token>
Content-Type: application/json

{
  "order_id": "ord123"
}

Response: 201 Created
{
  "id": "pay123",
  "payment_code": "PAY123456",
  "message": "Payment created"
}
```

### Update Payment Status
```
PUT /api/payments/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "success"
}

Response: 200 OK
{
  "message": "Payment status updated"
}
```

---

## ⚙️ Admin Endpoints (Admin Only)

### Create Product
```
POST /api/products
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Kaos Biru",
  "description": "Kaos biru premium",
  "price": 89000,
  "category_id": "cat123",
  "is_customizable": true
}

Response: 201 Created
{
  "id": "prod456",
  "message": "Product created"
}
```

### Update Product
```
PUT /api/products/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Kaos Biru Updated",
  "price": 95000
}

Response: 200 OK
{
  "message": "Product updated"
}
```

### Delete Product
```
DELETE /api/products/:id
Authorization: Bearer <admin-token>

Response: 200 OK
{
  "message": "Product deleted"
}
```

### Create Category
```
POST /api/categories
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Apparel",
  "description": "Clothing items"
}

Response: 201 Created
{
  "id": "cat456",
  "message": "Category created"
}
```

### Update Order Status
```
PUT /api/orders/:id
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "status": "shipped"
}

Response: 200 OK
{
  "message": "Order status updated"
}
```

---

## 📁 Project Structure

```
emyu-be-golang/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── config/
│   └── config.go                # Configuration loader
├── database/
│   ├── db.go                    # Database connection
│   └── schema.sql               # Database schema
├── handlers/
│   ├── auth.go                  # Auth endpoints
│   ├── product.go               # Product CRUD
│   ├── category.go              # Category CRUD
│   ├── cart.go                  # Cart management
│   ├── order.go                 # Order management
│   └── payment.go               # Payment handling
├── middleware/
│   └── auth.go                  # JWT & role middleware
├── models/
│   └── models.go                # Data structures
├── routes/
│   └── routes.go                # Route definitions
├── utils/
│   ├── jwt.go                   # JWT utilities
│   └── helpers.go               # Helper functions
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
├── .env.example                 # Environment template
└── README.md                    # Documentation
```

---

## 🔑 Order Statuses

- `pending` - Order created, awaiting payment
- `paid` - Payment received
- `packed` - Items packed and ready to ship
- `shipped` - Order shipped
- `delivered` - Order delivered
- `canceled` - Order canceled

## 💳 Payment Methods

- `qris` - QRIS code
- `bank_transfer` - Bank transfer
- `ewallet` - E-wallet (OVO, Dana, etc.)

## 🧪 Testing with cURL

### Register
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "phone": "08123456789"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Products (with token)
```bash
curl -X GET http://localhost:8080/api/products \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📝 Environment Variables

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=emyu

# Server Configuration
SERVER_PORT=8080
SERVER_ENV=development

# JWT Configuration
JWT_SECRET=your-secret-key-here

# App Info
APP_NAME=Emyu E-Commerce API
```

---

## 🚀 Deployment

### Build for Production

```bash
# Build binary
go build -o bin/api cmd/api/main.go

# Run binary
./bin/api
```

### Docker (Optional)

```bash
# Build image
docker build -t emyu-api .

# Run container
docker run -p 8080:8080 --env-file .env emyu-api
```

---

## 🐛 Troubleshooting

### Database Connection Error
- Check MySQL is running
- Verify credentials in `.env`
- Ensure database is created: `CREATE DATABASE emyu;`

### JWT Token Invalid
- Check JWT_SECRET in `.env`
- Token expires after 24 hours
- Include `Authorization: Bearer` prefix correctly

### Port Already in Use
- Change `SERVER_PORT` in `.env`
- Or kill process: `lsof -ti:8080 | xargs kill -9`

---

## 📄 License

MIT License - feel free to use this for any purpose.

---

## 👨‍💻 Development

Made with ❤️ for simple e-commerce needs.

**Quick Start:**
1. `cp .env.example .env` - Setup config
2. `mysql < database/schema.sql` - Create tables
3. `go run cmd/api/main.go` - Start server
4. Test with Postman or cURL

Happy coding! 🎉
