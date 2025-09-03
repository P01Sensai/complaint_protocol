package main

import (
	"sync"
	"time"
)


type User struct {
	ID         string      `json:"id"`
	SecretCode string      `json:"secret_code"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints,omitempty"` 
}

type Complaint struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Rating   int       `json:"rating"` 
	Resolved bool      `json:"resolved"`
	UserID   string    `json:"user_id"` 
	Date     time.Time `json:"date"`
}


type Storage struct {
	mu         sync.RWMutex
	Users      map[string]User
	Complaints map[string]Complaint
	AdminSecret string
}


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