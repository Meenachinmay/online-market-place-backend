# API Gateway Implementation Guide

This document defines the data shapes (Requests and Responses) required for the API Gateway to interact with the backend gRPC microservices. The API Gateway is responsible for exposing RESTful endpoints to the frontend/clients and mapping them to the internal gRPC calls.

## Core Concepts & User Identity
*   **User ID (`user_id`)**: The Gateway will authenticate users (e.g., via JWT). The authenticated user's ID is the source of truth.
*   **Buyer vs. Author**: 
    *   When a user makes a purchase, their `user_id` maps to `buyer_id`.
    *   When a user creates a blog, their `user_id` maps to `author_id`.
    *   In the Finance service, `user_id` is the wallet owner (regardless of whether they are a buyer or author).

---

## 1. Product Service

**Base Path**: `/api/v1/products`

### A. List Products
*   **Method**: `GET /`
*   **Description**: Returns a list of all available sneakers.
*   **gRPC Call**: `ProductService.ListProducts`

**Response Body (JSON):**
```json
{
  "products": [
    {
      "id": "uuid-string",
      "name": "Soda Air Max 1",
      "description": "Classic comfort...",
      "price": 100,
      "buyer_reward_points": 10
    }
  ]
}
```

### B. Get Product Details
*   **Method**: `GET /:id`
*   **Description**: Returns details for a specific sneaker.
*   **gRPC Call**: `ProductService.GetProduct`

**Response Body (JSON):**
```json
{
  "id": "uuid-string",
  "name": "Soda Air Max 1",
  "description": "Classic comfort...",
  "price": 100,
  "buyer_reward_points": 10
}
```

---

## 2. Referral Blog Service

**Base Path**: `/api/v1/blogs`

### A. List Blogs (Feed)
*   **Method**: `GET /`
*   **Description**: Returns a list of referral blogs to display on the home feed.
*   **gRPC Call**: `BlogService.ListBlogs`

**Response Body (JSON):**
```json
{
  "blogs": [
    {
      "id": "uuid-string",
      "author_id": "user-uuid-string",
      "content": "Check out these kicks!",
      "linked_product_id": "product-uuid-string"
    }
  ]
}
```

### B. Get Blog Details
*   **Method**: `GET /:id`
*   **Description**: specific blog details.
*   **gRPC Call**: `BlogService.GetBlog`

**Response Body (JSON):**
```json
{
  "id": "uuid-string",
  "author_id": "user-uuid-string",
  "content": "Check out these kicks!",
  "linked_product_id": "product-uuid-string"
}
```

### C. Create Blog (Authoring)
*   **Method**: `POST /`
*   **Description**: Authenticated user creates a new blog post linking to a product.
*   **gRPC Call**: `BlogService.CreateBlog`
*   **Identity Mapping**: Gateway must extract `user_id` from the auth token and pass it as `author_id`.

**Request Body (JSON):**
```json
{
  "content": "These shoes are fire!",
  "product_id": "product-uuid-string"
}
```

**Response Body (JSON):**
```json
{
  "id": "new-blog-uuid",
  "author_id": "user-uuid-from-token",
  "content": "These shoes are fire!",
  "linked_product_id": "product-uuid-string"
}
```

---

## 3. Order Service (Purchasing)

**Base Path**: `/api/v1/orders`

### A. Place Order
*   **Method**: `POST /`
*   **Description**: Authenticated user purchases a product via a blog link.
*   **gRPC Call**: `OrderService.PlaceOrder`
*   **Identity Mapping**: Gateway must extract `user_id` from the auth token and pass it as `buyer_id`.

**Request Body (JSON):**
```json
{
  "product_id": "product-uuid-string",
  "blog_id": "blog-uuid-string"
}
```
*Note: `blog_id` is required to attribute points to the author.*

**Response Body (JSON):**
```json
{
  "order": {
    "id": "order-uuid-string",
    "buyer_id": "user-uuid-from-token",
    "product_id": "product-uuid-string",
    "amount": 100,
    "status": "CONFIRMED",
    "created_at": 1678886400
  }
}
```

---

## 4. Finance Service (Wallet)

**Base Path**: `/api/v1/wallet`

### A. Get My Wallet
*   **Method**: `GET /`
*   **Description**: Returns the current balance and points for the authenticated user.
*   **gRPC Call**: `FinanceService.GetWallet`
*   **Identity Mapping**: Gateway must extract `user_id` from the auth token and pass it as `user_id` in the `UserRequest`.

**Response Body (JSON):**
```json
{
  "user_id": "user-uuid-from-token",
  "soda_points": 1500,
  "soda_balance": 250
}
```

### B. Convert Points
*   **Method**: `POST /convert`
*   **Description**: Converts Soda Points to Soda Balance (Yen).
*   **gRPC Call**: `FinanceService.ConvertPoints`
*   **Identity Mapping**: Gateway must extract `user_id` from the auth token and pass it as `user_id` in the `ConvertRequest`.

**Request Body (JSON):**
```json
{
  "user_id": "user-uuid-from-token",
  "points_to_convert": 1000
}
```
*Note: Backend logic enforces a minimum threshold (e.g., >1000 points required).*

**Response Body (JSON):**
```json
{
  "user_id": "user-uuid-from-token",
  "soda_points": 500,
  "soda_balance": 750
}
```
