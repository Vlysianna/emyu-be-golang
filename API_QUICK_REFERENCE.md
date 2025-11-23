# 🚀 Emyu E-Commerce API - Quick Reference

## Base URL
```
http://localhost:8080/api
```

## 🔐 Authentication
All protected endpoints require:
```
Authorization: Bearer <your_jwt_token>
```

---

## 📋 Quick API Reference

### Auth Endpoints
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| POST | `/register` | ❌ | Create new user account |
| POST | `/login` | ❌ | Login and get JWT token |
| POST | `/logout` | ✅ | Logout (delete token client-side) |

### Products & Categories (Public)
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/products` | ❌ | List all products |
| GET | `/products/:id` | ❌ | Get product details |
| POST | `/products` | ✅ Admin | Create product |
| PUT | `/products/:id` | ✅ Admin | Update product |
| DELETE | `/products/:id` | ✅ Admin | Delete product |
| GET | `/categories` | ❌ | List all categories |
| GET | `/categories/:id` | ❌ | Get category details |
| POST | `/categories` | ✅ Admin | Create category |
| PUT | `/categories/:id` | ✅ Admin | Update category |
| DELETE | `/categories/:id` | ✅ Admin | Delete category |

### Shopping Cart
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/carts` | ✅ | Get user cart |
| POST | `/carts` | ✅ | Create cart |
| POST | `/cart-items` | ✅ | Add item to cart |
| PUT | `/cart-items/:itemId` | ✅ | Update cart item |
| DELETE | `/cart-items/:itemId` | ✅ | Remove item from cart |

### Orders
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/orders` | ✅ | Get user orders |
| GET | `/orders/:id` | ✅ | Get order details |
| POST | `/orders` | ✅ | Create order (checkout) |
| PUT | `/orders/:id` | ✅ Admin | Update order status |
| DELETE | `/orders/:id` | ✅ Admin | Delete order |

### Payments
| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/payments` | ✅ | List payments |
| GET | `/payments/:id` | ✅ | Get payment details |
| POST | `/payments` | ✅ | Create payment |
| PUT | `/payments/:id` | ✅ | Update payment status |

---

## 💡 Common Request Examples

### 1. Register User
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

### 2. Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Get All Products
```bash
curl -X GET http://localhost:8080/api/products
```

### 4. Add to Cart
```bash
curl -X POST http://localhost:8080/api/cart-items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "cart_id": "cart123",
    "product_variant_id": "var123",
    "quantity": 2,
    "custom_name": "John",
    "custom_number": "1"
  }'
```

### 5. Create Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "total_amount": 200000,
    "shipping_cost": 25000,
    "payment_method": "qris",
    "shipping_address_id": "addr123"
  }'
```

### 6. Create Product (Admin)
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "name": "Kaos Putih",
    "description": "Kaos putih premium",
    "price": 89000,
    "category_id": "cat123",
    "is_customizable": true
  }'
```

---

## 📊 Data Models

### User
```json
{
  "id": "abc123",
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "08123456789",
  "role": "user|admin",
  "created_at": "2025-11-23T10:00:00Z",
  "updated_at": "2025-11-23T10:00:00Z"
}
```

### Product
```json
{
  "id": "prod123",
  "name": "Kaos Putih",
  "description": "Kaos putih premium",
  "price": 89000,
  "category_id": "cat123",
  "is_customizable": true,
  "created_at": "2025-11-23T10:00:00Z",
  "updated_at": "2025-11-23T10:00:00Z"
}
```

### Cart Item
```json
{
  "id": "item1",
  "cart_id": "cart123",
  "product_variant_id": "var123",
  "quantity": 2,
  "custom_name": "John",
  "custom_number": "1",
  "created_at": "2025-11-23T10:00:00Z",
  "updated_at": "2025-11-23T10:00:00Z"
}
```

### Order
```json
{
  "id": "ord123",
  "user_id": "user123",
  "order_number": "ORD-20251123-1234",
  "total_amount": 200000,
  "shipping_cost": 25000,
  "status": "pending|paid|packed|shipped|delivered|canceled",
  "payment_method": "qris|bank_transfer|ewallet",
  "shipping_address_id": "addr123",
  "created_at": "2025-11-23T10:00:00Z",
  "updated_at": "2025-11-23T10:00:00Z"
}
```

### Payment
```json
{
  "id": "pay123",
  "order_id": "ord123",
  "payment_status": "pending|success|failed",
  "payment_code": "PAY123456",
  "paid_at": "2025-11-23T10:00:00Z",
  "created_at": "2025-11-23T10:00:00Z",
  "updated_at": "2025-11-23T10:00:00Z"
}
```

---

## 🔑 Status Values

### Order Status
- `pending` - Order created, awaiting payment
- `paid` - Payment confirmed
- `packed` - Ready to ship
- `shipped` - In transit
- `delivered` - Successfully delivered
- `canceled` - Order canceled

### Payment Status
- `pending` - Awaiting payment
- `success` - Payment successful
- `failed` - Payment failed

### Payment Methods
- `qris` - QRIS QR code
- `bank_transfer` - Bank transfer
- `ewallet` - E-wallet (OVO, Dana, etc.)

### User Roles
- `user` - Regular customer
- `admin` - Administrator

---

## 🔓 Protected Routes

### User Protected (Auth Required)
- GET `/carts`
- POST `/carts`
- POST `/cart-items`
- PUT `/cart-items/:itemId`
- DELETE `/cart-items/:itemId`
- GET `/orders`
- GET `/orders/:id`
- POST `/orders`
- GET `/payments`
- GET `/payments/:id`
- POST `/payments`
- PUT `/payments/:id`
- POST `/logout`

### Admin Protected (Auth + Admin Role Required)
- POST `/products`
- PUT `/products/:id`
- DELETE `/products/:id`
- POST `/categories`
- PUT `/categories/:id`
- DELETE `/categories/:id`
- PUT `/orders/:id`
- DELETE `/orders/:id`

---

## ⚠️ Error Responses

### Bad Request (400)
```json
{
  "error": "Invalid request format"
}
```

### Unauthorized (401)
```json
{
  "error": "Authorization header required"
}
```

### Forbidden (403)
```json
{
  "error": "Admin access required"
}
```

### Not Found (404)
```json
{
  "error": "Product not found"
}
```

### Conflict (409)
```json
{
  "error": "Email already registered"
}
```

### Server Error (500)
```json
{
  "error": "Failed to create user"
}
```

---

## 🚀 Development Setup

```bash
# 1. Clone and install
git clone <repo-url>
cd emyu-be-golang
make install

# 2. Setup database
make db-create
make db-migrate

# 3. Configure environment
cp .env.example .env
# Edit .env with your settings

# 4. Start development server
make dev

# 5. Test API (use Postman or import the collection)
```

---

## 📱 Frontend Integration Tips

### 1. Store Token
```javascript
// After login/register
localStorage.setItem('token', response.data.token);
localStorage.setItem('user', JSON.stringify(response.data.user));
```

### 2. Use Token in Requests
```javascript
const config = {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`
  }
};
axios.get('/api/orders', config);
```

### 3. Handle Logout
```javascript
// Clear token and user
localStorage.removeItem('token');
localStorage.removeItem('user');
// Redirect to login
```

### 4. Error Handling
```javascript
try {
  const response = await axios.post('/api/orders', data);
} catch (error) {
  if (error.response?.status === 401) {
    // Redirect to login
  } else {
    // Show error message
    console.error(error.response?.data?.error);
  }
}
```

---

## 🎯 Workflow Example

```
1. Register/Login
   POST /api/register or POST /api/login
   ↓ Get JWT token
   
2. Browse Products
   GET /api/products
   GET /api/products/:id
   ↓
   
3. Add to Cart
   POST /api/carts (create if needed)
   POST /api/cart-items (add items)
   ↓
   
4. Checkout
   POST /api/shipping-addresses (add address)
   POST /api/orders (create order)
   ↓
   
5. Payment
   POST /api/payments (create payment)
   PUT /api/payments/:id (update status)
   ↓
   
6. Order Tracking
   GET /api/orders (view all orders)
   GET /api/orders/:id (view specific order)
```

---

## 📞 Support

For issues or questions:
- Check README.md for detailed documentation
- Review database schema in `database/schema.sql`
- Check handler implementations in `handlers/` folder
- Review middleware in `middleware/auth.go`
