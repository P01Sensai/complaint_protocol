# Complaint Portal API

## 📋 Project Overview
A RESTful JSON API built in Go for a complaint management system where users can submit complaints and administrators can review and resolve them.

## 🏗️ Project Structure
complaint_portal/
- ├── main.go ---> Server initialization and routing
- ├── go.mod ---> Go module definition
- ├── models.go ---> Data structures (User, Complaint)
- ├── storage.go ---> In-memory database with thread-safe operations
- └── handlers.go ---> HTTP request handlers for all endpoints


## 🔌 API Endpoints

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/register` | POST | Create new user | None |
| `/login` | POST | User login | None |
| `/submitComplaint` | POST | Submit a complaint | User Secret |
| `/getAllComplaintsForUser` | GET | Get user's complaints | User Secret |
| `/getAllComplaintsForAdmin` | GET | Get all complaints (admin view) | Admin Secret |
| `/viewComplaint` | GET | View specific complaint | User/Admin Secret |
| `/resolveComplaint` | POST | Resolve a complaint | Admin Secret |

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or later

### Installation & Running
1. Place all files in the same directory
2. Run the application:
   ```bash
   go run main.go models.go storage.go handlers.go
3. Server starts on http://localhost:8080
## 🔐 Authentication
### User Authentication
- Users receive a secret_code after registration
- Include this secret in the X-Secret-Code header for user endpoints

### Admin Authentication
- Use the admin secret: admin123
- Include in the X-Admin-Secret header for admin endpoints

## Data Models

### User Structure
    ```json
    {
     "id": "string",
     "secret_code": "string",
     "name": "string", 
     "email": "string",
     "complaints": "[]Complaint"
     }

### Complaint Structure

    ```json
    {
     "id": "string",
     "title": "string",
     "summary": "string",
     "rating": "number",
     "resolved": "boolean",
     "user_id": "string",
     "date": "timestamp"
    }
## 🔄 API Flow
### User Registration Flow
- POST to /register with name and email
- Receive user details including secret_code
- Use secret_code for all subsequent authenticated requests

### Complaint Submission Flow
- POST to /submitComplaint with complaint details + X-Secret-Code header
- Receive complaint details including unique complaint_id
- Use complaint_id to view or reference specific complaints

### Admin Management Flow
- Use X-Admin-Secret: admin123 header
- GET /getAllComplaintsForAdmin to view all complaints
- POST /resolveComplaint to mark complaints as resolved

## 📝 Example Usage Scenarios

1. User submits a complaint: Register → Login → Submit Complaint
2. User views their complaints: Login → Get All Complaints
3. Admin reviews complaints: Use admin secret → Get All Complaints (Admin)
4. Admin resolves complaint: Use admin secret → Resolve Complaint
