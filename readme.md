# Budget Management API

A comprehensive REST API for personal budget and expense management built with Go, Echo framework, and MySQL.

## 🚀 Features

- **User Management**
  - User registration and authentication
  - JWT-based authorization
  - Password reset with email verification
  - Profile management

- **Budget Tracking**
  - Create and manage monthly budgets
  - Associate budgets with multiple categories
  - Track budget amounts by year and month
  - Automatic slug generation

- **Transaction Management**
  - Income and expense tracking
  - Transaction reversal support
  - Wallet-based transactions
  - Category-based filtering
  - Date range filtering

- **Wallet System**
  - Multiple wallet support per user
  - Default wallet generation (Cash, Bank)
  - Balance tracking
  - Wallet-to-wallet transfers

- **Categories**
  - Pre-defined and custom categories
  - Category-based budget allocation
  - Automatic category seeding

## 🛠️ Tech Stack

- **Language**: Go 1.25
- **Framework**: Echo v5
- **Database**: MySQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Validation**: go-playground/validator
- **Email**: gomail v2
- **Password Hashing**: bcrypt

## 📋 Prerequisites

- Go 1.25 or higher
- MySQL 5.7+
- SMTP server for email functionality

## ⚙️ Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd budget-be
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**

Create a `.env` file in the root directory:

```env
# Application
APP_NAME=Budget Manager
APP_PORT=3000

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=your_username
DB_DATABASE=budget_db
DB_PASSWORD=your_password

# JWT
JWT_SECRET=your_jwt_secret_key

# Email
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=your_email@example.com
MAIL_PASSWORD=your_email_password
MAIL_SENDER=noreply@example.com
```

4. **Run database migrations**
```bash
go run internal/database/migrate_up.go
```

5. **Seed categories**
```bash
go run internal/database/seeders/category_seeder.go
```

6. **Start the server**
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:3000`

## 📚 API Endpoints

### Authentication
```
POST   /api/auth/register          - Register new user
POST   /api/auth/login             - Login user
POST   /api/auth/forgot/password   - Request password reset
POST   /api/auth/reset/password    - Reset password
```

### Profile (Protected)
```
GET    /api/profile/authenticated/user  - Get authenticated user
PATCH  /api/profile/update/password     - Change password
```

### Categories (Protected)
```
GET    /api/categories/all         - List all categories
POST   /api/categories/create      - Create custom category
DELETE /api/categories/delete/:id  - Delete category
```

### Budgets (Protected)
```
POST   /api/budgets/create         - Create budget
GET    /api/budgets/all            - List budgets (paginated)
PUT    /api/budgets/:id/update     - Update budget
DELETE /api/budgets/:id/delete     - Delete budget
```

### Wallets (Protected)
```
POST   /api/wallets/create              - Create wallet
GET    /api/wallets/generate-default    - Generate default wallets
GET    /api/wallets/user-list           - List user wallets
```

### Transactions (Protected)
```
POST   /api/transactions/create          - Create transaction
PUT    /api/transactions/:id/reverse     - Reverse transaction
GET    /api/transactions/list            - List transactions (filtered & paginated)
```

### Transfers (Protected)
```
POST   /api/transfer                     - Transfer between wallets
```

## 🔐 Authentication

All protected routes require a Bearer token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## 📊 Query Parameters

### Transactions List
```
?page=1                          - Page number
?limit=10                        - Items per page
?from_date=2024-01-01           - Filter from date (YYYY-MM-DD)
?end_date=2024-01-31            - Filter to date (YYYY-MM-DD)
?category_id=1                   - Filter by category
?wallet_id=1                     - Filter by wallet
?type=income                     - Filter by type (income/expense)
?month=1                         - Filter by month (1-12)
?year=2024                       - Filter by year
```

### Budgets List
```
?page=1                          - Page number
?limit=10                        - Items per page
```

## 🗄️ Database Models

### User
- FirstName, LastName, Email, Password, Gender
- Soft delete support

### Budget
- Title, Description, Amount, Slug
- Year and Month tracking
- Many-to-many relationship with Categories

### Transaction
- Amount, Date, Type (income/expense)
- Reversal support (IsReversal, ParentId)
- Links to Wallet, Category, and User
- Month and Year auto-populated

### Wallet
- Name, Balance
- Unique per user (unique constraint on user_id + name)

### Category
- Name, Slug, IsCustom
- Pre-seeded system categories

## 🔄 Key Features Explained

### Soft Delete
All models use GORM's soft delete feature - deleted records are marked with `deleted_at` timestamp instead of being permanently removed.

### Transaction Reversal
Transactions can be reversed, creating a compensating transaction that:
- Reverses the wallet balance change
- Updates budget amounts (if applicable)
- Links to original transaction via `parent_id`
- Marked with `is_reversal = true`

### Wallet Transfers
Transfers create two transactions:
- Expense transaction from source wallet
- Income transaction to destination wallet
- Atomic operation using database transactions

### Budget Tracking
Budgets automatically update when:
- Expense transactions are created (decrements budget)
- Income reversals occur (increments budget)
- Uses atomic SQL expressions (`amount + ?`) to prevent race conditions

### Concurrent Operations
The API uses goroutines and `sync.WaitGroup` for parallel operations:
- Fetching multiple wallets simultaneously in transfers
- Improved response times for independent database queries

## 🔧 Development

### Project Structure
```
budget-be/
├── cmd/api/
│   ├── filters/          # Query filters
│   ├── handlers/         # HTTP handlers
│   ├── middlewares/      # Custom middlewares
│   ├── requests/         # Request DTOs
│   ├── services/         # Business logic
│   ├── main.go          # Application entry
│   └── routes.go        # Route definitions
├── common/              # Shared utilities
│   ├── pagination.go    # Pagination helper
│   ├── response.go      # Response helpers
│   └── scopes.go        # GORM query scopes
├── internal/
│   ├── database/        # Migrations & seeders
│   ├── mailer/          # Email templates & service
│   ├── models/          # Database models
│   └── custom_errors/   # Custom error types
└── go.mod
```

### Running Tests
```bash
go test ./...
```

## 📝 Example Requests

### Register User
```json
POST /api/auth/register
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "gender": "male"
}
```

### Create Budget
```json
POST /api/budgets/create
Authorization: Bearer <token>

{
  "title": "Monthly Groceries",
  "amount": 500.00,
  "categories": [1, 2],
  "date": "2024-01-01",
  "description": "Food budget for January"
}
```

### Create Transaction
```json
POST /api/transactions/create
Authorization: Bearer <token>

{
  "wallet_id": 1,
  "category_id": 1,
  "amount": 50.00,
  "type": "expense",
  "title": "Grocery shopping",
  "description": "Weekly groceries",
  "date": "2024-01-15"
}
```

### Transfer Between Wallets
```json
POST /api/transfer
Authorization: Bearer <token>

{
  "source_wallet_id": 1,
  "destination_wallet_id": 2,
  "amount": 100.00
}
```

### List Transactions with Filters
```
GET /api/transactions/list?page=1&limit=20&type=expense&month=1&year=2024
Authorization: Bearer <token>
```

## 🎯 Validation Rules

### Transaction
- Amount: Required, minimum 1
- Type: Required, must be "income" or "expense"
- WalletId: Required
- Date: Optional, format YYYY-MM-DD
- Title: Optional, max 100 characters
- Description: Optional, max 490 characters

### Budget
- Title: Required, max 100 characters
- Amount: Required, minimum 0
- Date: Required, format YYYY-MM-DD
- Categories: Required array of category IDs

### Wallet
- Name: Required, max 100 characters
- Amount: Required, minimum 0

## 🔒 Security Features

- Password hashing using bcrypt
- JWT token-based authentication
- Protected routes with authentication middleware
- User-scoped data access (users can only access their own data)
- Email verification for password reset
- Unique constraints on sensitive fields

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👤 Author

Shivanshu Tripathi - [LinkedIn](https://www.linkedin.com/in/shivanshu-tripathi-662405217/)

## 🙏 Acknowledgments

- Echo Framework
- GORM
- Go Playground Validator
- golang-jwt/jwt
