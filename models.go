package main

import (
	"sync"
	"time"
)

// User represents a user in the system
type User struct {
	ID         string      `json:"id"`
	SecretCode string      `json:"secret_code"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints,omitempty"` // Optional field
}

// Complaint represents a complaint submitted by a user
type Complaint struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Rating   int       `json:"rating"` // Severity rating (1-5 or similar)
	Resolved bool      `json:"resolved"`
	UserID   string    `json:"user_id"` // ID of the user who submitted the complaint
	Date     time.Time `json:"date"`
}

// Storage represents the in-memory database with mutex for concurrency safety
type Storage struct {
	mu         sync.RWMutex
	Users      map[string]User
	Complaints map[string]Complaint
	// Admin credentials (hardcoded for simplicity, in real app use proper auth)
	AdminSecret string
}

// Response structures for JSON responses
type LoginResponse struct {
	User User `json:"user"`
}

type RegisterResponse struct {
	User User `json:"user"`
}

type ComplaintResponse struct {
	Complaint Complaint `json:"complaint"`
}

type ComplaintsListResponse struct {
	Complaints []Complaint `json:"complaints"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}