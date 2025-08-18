# 🎵 Music Server API

A RESTful Music Server API built in **Golang**, supporting user signup/login, authentication, logging, and database persistence. Designed for modularity, concurrency, and production readiness.

---

- **User Authentication**: Signup and login endpoints with password hashing (`bcrypt`) and JWT tokens, cookies too.  
- **Middleware**: Logging and authentication middleware for secure API routes.  
- **Routing**: Handled using `Gorilla Mux`.  
- **Database**: PostgreSQL (`avien`) integration via `GORM` ORM.  
- **Configuration**: Environment variables managed with `godotenv`.  
- **Daemon Compilation**: Server can run as a background daemon.  
- **Models & Structs**: Clear separation of models and request/response structs.

---
## Tech Stack

- **Language**: Go  
- **Routing**: Gorilla Mux  
- **ORM**: GORM  
- **Database**: PostgreSQL  
- **Authentication**: JWT + bcrypt  
- **Environment Config**: godotenv  
- **Logging**: Custom middleware  
- **Build/Daemon**: go build / systemd or manual daemon mode
---

## Installation

### Prerequisites

- Go 1.20+ installed  
- PostgreSQL database running (avien)  
- `git` installed  

### Steps

1. Clone the repository:

```bash
git clone https://github.com/yourusername/music-server-api.git
cd music-server-api

2. Install Dependencies:
go mod tidy

3. Create .env file in the project root:

PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=youruser
DB_PASS=yourpassword
DB_NAME=musicdb
JWT_SECRET=yourjwtsecret

4. go run main.go

API Endpoints
AuthMethod	Endpoint	Description
POST		/signup		Register new user
POST		/login		Login user

All protected routes require Authorization: Bearer <JWT> header.

Health
Method	Endpoint	Description
GET	/health	Health check
Middleware

Logging: Logs request method, path, and status code.

Authentication: Verifies JWT for protected routes.

Dependencies

github.com/gorilla/mux – routing

gorm.io/gorm – ORM

gorm.io/driver/postgres – PostgreSQL driver

golang.org/x/crypto/bcrypt – password hashing

github.com/golang-jwt/jwt – JWT handling

github.com/joho/godotenv – env management


License

MIT License © [franzego]
