# Go Book Inventory API

A robust, production-ready RESTful API for library and inventory management built with Go (Golang). This project demonstrates advanced backend architecture, relational database management, and secure authorization mechanisms.

## 🚀 Key Features

*   **RESTful JSON API:** Fully headless architecture responding with structured JSON and proper HTTP status codes.
*   **Secure Authentication:** User registration and login utilizing **Bcrypt** for password hashing and **JWT** (JSON Web Tokens) for stateless authentication.
*   **Role-Based Access Control (RBAC):** Distinct permission levels separating `Admin` (full CRUD access) and `Member` (read and borrow access).
*   **Database Transactions:** Safe borrowing and returning mechanisms using GORM's `TX` functions to ensure stock data consistency and prevent race conditions.
*   **Dynamic Queries:** Implemented search filtering (`ILIKE`) and mathematical pagination (limit/offset) for optimal data retrieval on large datasets.
*   **Automated Initialization:** Built-in auto-migration and data seeders to populate initial admin credentials and sample book catalogs upon the first run.

## 🛠️ Tech Stack

*   **Language:** Go (1.16+)
*   **Framework:** Gin-Gonic
*   **ORM:** GORM
*   **Database:** PostgreSQL
*   **Security:** Golang-JWT (v4), X/Crypto (Bcrypt)

## ⚙️ Prerequisites

*   Go installed on your local machine.
*   PostgreSQL server running locally or remotely.

## 📦 Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone [https://github.com/firdhausranggaa/gobook-inventory.git](https://github.com/firdhausranggaa/gobook-inventory.git)
   cd gobook-inventory

```

2. **Environment Variables:**
Create a `.env` file in the root directory and define the following variables:
```env
POSTGRES_URL="postgresql://postgres:password@127.0.0.1/postgres?sslmode=disable"
SUPER_USER="admin"
SUPER_PASS="123"
SUPER_SECRET="your-secure-jwt-secret-key"

```


3. **Install Dependencies:**
```bash
go mod tidy

```


4. **Run the Server:**
```bash
go run main.go

```


*Note: Upon the first successful run, the API will automatically migrate tables and seed the initial data.*

## 📡 API Endpoints

All protected routes require an `Authorization` header with the format: `Bearer <your_token>`.

### Authentication (Public)

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/register` | Register a new member account |
| `POST` | `/api/login` | Authenticate and receive a JWT token |

### Book Management (Protected)

| Method | Endpoint | Access Role | Description |
| --- | --- | --- | --- |
| `GET` | `/api/books` | Admin / Member | List all books. Supports `?search=`, `?page=`, `?limit=` |
| `GET` | `/api/books/:id` | Admin / Member | Get details of a specific book |
| `POST` | `/api/books` | **Admin Only** | Add a new book to the inventory |
| `PUT` | `/api/books/:id` | **Admin Only** | Update an existing book's details |
| `DELETE` | `/api/books/:id` | **Admin Only** | Remove a book from the inventory |

### Borrowing System (Protected)

| Method | Endpoint | Access Role | Description |
| --- | --- | --- | --- |
| `POST` | `/api/borrow` | Admin / Member | Borrow a book (requires `book_id` in JSON body) |
| `POST` | `/api/return/:id` | Admin / Member | Return a borrowed book by borrowing ID |