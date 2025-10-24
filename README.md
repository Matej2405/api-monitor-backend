# API Monitor - Backend

High-performance REST API backend for real-time API monitoring and analysis, built with Go.

## 🚀 Live Demo

**Production API:** [https://api-monitor-backend-production.up.railway.app](https://api-monitor-backend-production.up.railway.app)

**Health Check:** [https://api-monitor-backend-production.up.railway.app/api/requests](https://api-monitor-backend-production.up.railway.app/api/requests)

The application is deployed on **Railway** with automatic deployments from the main branch.

---

## 🎨 Design Reference

This backend implements the data layer for **Treblle's official design system** from the "Ship Happens" Hackathon 2025:

**Figma Prototype:**
- [Interactive Prototype](https://www.figma.com/proto/9RokOq6XAby6le7ePwNTj0/Treblle-"Ship-Happens"-Hackathon-2025---additional-task)

---

## 🛠 Tech Stack

### **Core Technologies:**
- **Go 1.23** - Backend language ✅ BONUS
- **Gorilla Mux** - HTTP router and middleware
- **SQLite** - Embedded database ✅ BONUS
- **CORS middleware** - Cross-origin resource sharing

### **Documentation:**
- **Swagger/OpenAPI** - API documentation ✅ BONUS
- **swaggo** - Swagger annotation generator

### **Database:**
- **SQLite3** - Local file-based database
- **database/sql** - Go standard library SQL driver
- **mattn/go-sqlite3** - CGO-enabled SQLite driver

---

## ✨ Features

### **Core Functionality:**
- ✅ **RESTful API design** with proper HTTP methods
- ✅ **Request logging** - Captures all API requests with detailed metrics
- ✅ **Problem detection** - Automatically identifies and categorizes API issues
- ✅ **Advanced filtering** - By method, response code, time range, response time
- ✅ **Flexible sorting** - By creation time or response time (asc/desc)
- ✅ **Search capability** - Full-text search across request paths
- ✅ **JSONPlaceholder proxy** - Test API with real-world data
- ✅ **Automatic seeding** - Populates database with sample data on startup

### **Problem Categories:**
The system automatically detects and categorizes these problems:
- **5xx errors** - Server-side errors (500-599)
- **4xx errors** - Client-side errors (400-499)
- **Slow responses** - Requests exceeding 1000ms
- **Timeouts** - Requests that timed out
- **Rate limiting** - HTTP 429 responses

---

## 🚀 Local Development

### **Prerequisites:**
- Go 1.23 or higher installed
- GCC compiler (for SQLite CGO support)
  - **Windows:** Install MinGW or TDM-GCC
  - **macOS:** Xcode Command Line Tools (`xcode-select --install`)
  - **Linux:** `sudo apt-get install build-essential`

### **Installation:**
```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/api-monitor-backend.git
cd api-monitor-backend

# Download dependencies
go mod download

# Verify CGO is enabled (should show "CGO_ENABLED=1")
go env CGO_ENABLED
```

### **Run Development Server:**
```bash
go run cmd/api/main.go
```

The server will start on **http://localhost:8080**

### **Build for Production:**
```bash
# Build binary
go build -o api-monitor ./cmd/api

# Run the binary
./api-monitor
```

---

## 📚 API Documentation

### **OpenAPI/Swagger Documentation:**

When running locally, full interactive API documentation is available at:

**Swagger UI:** [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

The API is documented using OpenAPI 3.0 specification with complete request/response schemas.

---

## 🔌 API Endpoints

### **Requests**

#### **GET /api/requests**
Get all API requests with optional filtering and sorting.

**Query Parameters:**
- `method` (string) - Filter by HTTP method (GET, POST, PUT, DELETE)
- `response_code` (int) - Filter by exact response code
- `min_response_code` (int) - Minimum response code
- `max_response_code` (int) - Maximum response code
- `min_response_time` (int) - Minimum response time in ms
- `max_response_time` (int) - Maximum response time in ms
- `start_date` (string) - Filter by start date (RFC3339 format)
- `end_date` (string) - Filter by end date (RFC3339 format)
- `search` (string) - Search in request path
- `sort_by` (string) - Sort field: `created_at` or `response_time` (default: created_at)
- `order` (string) - Sort order: `asc` or `desc` (default: desc)

**Example:**
```bash
curl "http://localhost:8080/api/requests?method=GET&min_response_code=200&max_response_code=299&sort_by=response_time&order=asc"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "method": "GET",
      "path": "/users",
      "response_code": 200,
      "response_time": 45,
      "created_at": "2025-10-24T02:44:37Z"
    }
  ]
}
```

---

### **Problems**

#### **GET /api/problems**
Get all detected API problems with optional filtering and sorting.

**Query Parameters:**
- `problem_type` (string) - Filter by type: `5xx_error`, `4xx_error`, `slow_response`, `timeout`, `rate_limit`
- `severity` (string) - Filter by severity: `high`, `medium`, `low`
- `sort_by` (string) - Sort field
- `order` (string) - Sort order: `asc` or `desc`

**Example:**
```bash
curl "http://localhost:8080/api/problems?problem_type=5xx_error&severity=high"
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "request_id": 42,
      "problem_type": "5xx_error",
      "severity": "high",
      "description": "Server error occurred",
      "detected_at": "2025-10-24T02:44:37Z"
    }
  ]
}
```

---

### **Proxy**

#### **ALL /api/proxy/{endpoint}**
Proxy requests to JSONPlaceholder API (https://jsonplaceholder.typicode.com).

Supports all HTTP methods: GET, POST, PUT, DELETE

**Example:**
```bash
# Get all users
curl http://localhost:8080/api/proxy/users

# Get specific user
curl http://localhost:8080/api/proxy/users/1

# Create post
curl -X POST http://localhost:8080/api/proxy/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","body":"Content","userId":1}'
```

All proxied requests are automatically logged and analyzed for problems.

---

## 🗄 Database Schema

### **api_requests table**
```sql
CREATE TABLE api_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    response_code INTEGER NOT NULL,
    response_time INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### **problems table**
```sql
CREATE TABLE problems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id INTEGER NOT NULL,
    problem_type TEXT NOT NULL,
    severity TEXT NOT NULL,
    description TEXT,
    detected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES api_requests(id)
);
```

---

## 🐳 Docker

The application is fully dockerized with multi-stage builds for optimal image size.

### **Run with Docker Compose:**

From the **root directory** (parent of frontend and backend):
```bash
docker-compose up --build
```

The backend will be available at **http://localhost:8080**

### **Docker Setup Details:**
- **Build stage:** golang:alpine with gcc for SQLite compilation
- **Runtime stage:** alpine:latest with minimal dependencies
- **Database:** SQLite file persisted via Docker volume
- **Auto-deployment:** Railway rebuilds on every git push

### **Dockerfile Features:**
- ✅ Multi-stage build for smaller image size
- ✅ CGO enabled for SQLite support
- ✅ Security: runs as non-root user
- ✅ Health checks included
- ✅ Automatic schema migration on startup

---

## 📁 Project Structure
```
api-monitor-backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── database/
│   │   ├── database.go          # Database connection & initialization
│   │   └── schema.sql           # Database schema
│   ├── handlers/
│   │   └── handlers.go          # HTTP request handlers
│   ├── models/
│   │   └── models.go            # Data models & types
│   └── utils/
│       └── seed.go              # Database seeding utility
├── docs/                        # Generated Swagger documentation
├── Dockerfile                   # Container configuration
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
└── README.md                    # This file
```

---

## 🔧 Configuration

### **Environment Variables:**

The application uses these environment variables (all optional):

- `PORT` - Server port (default: 8080)
- `DB_PATH` - SQLite database file path (default: ./api_monitor.db)
- `CORS_ORIGINS` - Allowed CORS origins (default: *)

### **Example .env file:**
```env
PORT=8080
DB_PATH=./data/api_monitor.db
CORS_ORIGINS=https://your-frontend.vercel.app
```

---

## 🚀 Deployment

### **Railway Deployment:**

The application is configured for automatic deployment on Railway:

1. **Connect GitHub repository** to Railway
2. **Railway auto-detects** Go application
3. **Automatic builds** on every push to main branch
4. **Environment variables** configured in Railway dashboard
5. **Persistent storage** for SQLite database via Railway volumes

**Production URL:** https://api-monitor-backend-production.up.railway.app

### **Deployment Configuration:**
- **Build command:** Automatic (Railway detects Go)
- **Start command:** `./api-monitor`
- **Health check:** `/api/requests` endpoint
- **CORS:** Configured to allow all origins

---

## 🧪 Testing

### **Manual Testing:**
```bash
# Test requests endpoint
curl http://localhost:8080/api/requests

# Test with filters
curl "http://localhost:8080/api/requests?method=GET&sort_by=response_time"

# Test problems endpoint
curl http://localhost:8080/api/problems

# Test proxy
curl http://localhost:8080/api/proxy/users/1
```

### **Health Check:**
```bash
curl http://localhost:8080/api/requests
# Should return JSON with data array
```

---

## 📊 Database Seeding

The application automatically seeds the database with sample data on startup if it's empty:

- **50+ sample requests** across various endpoints
- **Mix of HTTP methods** (GET, POST, PUT, DELETE)
- **Various response codes** (200, 400, 404, 500, 503)
- **Realistic response times** (10ms - 3000ms)
- **Automatic problem detection** based on response codes and times

This provides immediate data for testing and demo purposes.

---

## 🔗 Related

**Frontend Repository:** [api-monitor-frontend](https://github.com/YOUR_USERNAME/api-monitor-frontend)

**Live Frontend:** [Vercel Deployment URL]

---

## 📝 License

This project was created for the **Treblle "Ship Happens" Hackathon 2025**.

---

## 👨‍💻 Author

Built with ❤️ for the Treblle Hackathon

---

## 🏆 Hackathon Requirements Met

### **Base Requirements:** ✅ 100%
- ✅ RESTful API design
- ✅ Full git flow with separate repository
- ✅ List/table view support with all required fields
- ✅ Sorting by created_at and response_time
- ✅ Filtering by method, response code, and time
- ✅ Search functionality
- ✅ Problem object with full CRUD operations

### **Bonus Requirements:** ✅ 6/9
- ✅ Backend written in **Go** (not PHP)
- ✅ **OpenAPI/Swagger** documentation
- ✅ **SQLite** local database
- ✅ **Dockerized** with multi-stage builds
- ✅ **Deployed** to production (Railway)
- ❌ Tests not implemented (time constraints)

---

## 🐛 Troubleshooting

### **CGO errors:**
If you get CGO-related errors, ensure you have a C compiler installed:
```bash
# Check CGO is enabled
go env CGO_ENABLED

# Should return: 1
```

### **Database locked:**
If you get "database is locked" errors:
```bash
# Stop all running instances
pkill api-monitor

# Remove lock file
rm api_monitor.db-shm api_monitor.db-wal
```

### **Port already in use:**
```bash
# Change port
PORT=8081 go run cmd/api/main.go
```

---

**🚀 Ready to monitor APIs at scale!**
