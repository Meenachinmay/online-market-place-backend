# Role
Act as a Senior Backend Engineer (Go & gRPC Expert).
We are building a monolithic gRPC backend for "The Consult App".

# Context & Architecture
* **Architecture:** Modular Monolith using Go.
* **Communication:** gRPC (Protobuf).
* **Database:** SQL (Postgres or MySQL compatible).
* **Domain:** A referral-based marketplace where users write blogs to sell products.
* **Economy:** Users earn "Soda Points" which can be converted to "Soda Balance" (Yen).

# Task
You need to generate the **Protobuf definitions**, **Database Migrations**, and **Core Business Logic** for the following domains.

## 1. Proto Definitions
Define the messages and RPC services in the following files.
*Note: Do NOT include image URLs in any messages. Images are handled client-side.*

**`product.proto`**
* **Message `Product`**: ID, Name, Description, Price (int64), `buyer_reward_points` (int32).
    * *Constraint:* Do NOT expose `author_reward_points` in the public `Product` message returned to buyers. That data is internal.
* **Service `ProductService`**:
    * `GetProduct(ProductRequest) returns (Product)`
    * `ListProducts(Empty) returns (ProductList)`

**`referral_blog.proto`**
* **Message `Blog`**: ID, AuthorID, Content, `linked_product_id`.
* **Service `BlogService`**:
    * `CreateBlog(CreateBlogRequest) returns (Blog)`: Input includes content and product_id.
    * `ListBlogs(Empty) returns (BlogList)`: For the home screen.
    * `GetBlog(BlogRequest) returns (Blog)`: Detailed view.

**`order.proto`**
* **Message `Order`**: ID, BuyerID, ProductID, Amount, Status, CreatedAt.
* **Service `OrderService`**:
    * `PlaceOrder(PlaceOrderRequest) returns (OrderResponse)`: Input includes `blog_id` (to track who gets the referral points).

**`finance.proto`**
* **Message `Wallet`**: UserID, `soda_points` (int64), `soda_balance` (int64/Yen).
* **Service `FinanceService`**:
    * `GetWallet(UserRequest) returns (Wallet)`
    * `ConvertPoints(ConvertRequest) returns (Wallet)`: Logic defined below.

---

## 2. Database Migrations (SQL)
Write migration scripts to create these tables. ensure strict typing and foreign keys.

1.  **products**: `id`, `name`, `description`, `price`, `buyer_reward_points`, `author_reward_points` (This column exists in DB but not Proto).
2.  **blogs**: `id`, `author_id`, `content`, `product_id` (FK to products).
3.  **orders**: `id`, `buyer_id`, `product_id`, `blog_id` (used for attribution), `amount`, `created_at`.
4.  **wallets**: `user_id` (PK), `soda_points` (default 0), `soda_balance` (default 0).
5.  **transactions**: `id`, `user_id`, `type` (EARNED, SPENT, CONVERTED), `amount`, `related_order_id`.

---

## 3. Business Logic (Go Implementation)
Provide the Go code snippets or pseudo-code for the critical logic flows below.

**Flow A: Placing an Order (`PlaceOrder`)**
When a user buys a product via a blog link:
1.  **Verify:** Check if the user has purchased *any* product before.
2.  **Buyer Reward:** IF it is the user's **first purchase ever**, credit their wallet with `products.buyer_reward_points`.
3.  **Author Reward:** ALWAYS credit the blog author's wallet with `products.author_reward_points`.
4.  **Logging:** Record these point movements in the `transactions` table.

**Flow B: Converting Points (`ConvertPoints`)**
1.  **Rate:** 2 Soda Points = 1 Yen (Soda Balance).
2.  **Threshold:** Conversion is ONLY allowed if `soda_points > 1000`.
3.  **Validation:** Check if user has enough points.
4.  **Execution:** Deduct points, add calculated amount to `soda_balance`. Update atomically.

## Output Requirements
* Produce strict `.proto` syntax (proto3).
* Produce standard SQL for migrations.
* Write clean, idiomatic Go code for the logic, using a transaction manager for the Finance parts (to prevent race conditions on money/points).
* Transport (gRPC) -> Business Logic (Core) -> Data Access (Stores) -> Database (SQLC/PGX).