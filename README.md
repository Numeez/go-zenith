# Go Zenith 🏋️‍♂️

Go Zenith is a backend service built in **Golang** that provides **user management and workout tracking** functionality.  
It follows clean backend practices such as authentication, migrations, testing, and containerized infrastructure.

This project is intended to demonstrate real-world backend engineering concepts using Go.

---

## ✨ Features

### 👤 User Management
- User registration
- User login
- Secure password storage using **hashed & encrypted passwords**
- JWT-based authentication and authorization

### 🏋️ Workout Management
- Create workouts
- Update workouts
- Delete workouts
- Fetch user-specific workouts
- Protected routes (only authenticated users can manage workouts)

### 🔐 Authentication & Security
- Password hashing (no plaintext passwords stored)
- JWT tokens for stateless authentication
- Middleware-based auth validation

---

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **ORM / DB Access:** `database/sql`
- **Migrations:** Goose
- **Authentication:** JWT
- **Password Encryption:** bcrypt
- **Containerization:** Docker & Docker Compose
- **Testing:** Go `testing` package
- **Environment Management:** `.env` files


