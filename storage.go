package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// NewStorage creates and initializes a new Storage instance
func NewStorage() *Storage {
	return &Storage{
		Users:       make(map[string]User),
		Complaints:  make(map[string]Complaint),
		AdminSecret: "admin123", // Simple admin secret for demonstration
	}
}

// generateID generates a unique random ID
func generateID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// generateSecretCode generates a unique secret code for users
func generateSecretCode() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based code if random fails
		return fmt.Sprintf("sc%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// AddUser adds a new user to the storage
func (s *Storage) AddUser(name, email string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if email already exists
	for _, user := range s.Users {
		if user.Email == email {
			return User{}, fmt.Errorf("email already registered")
		}
	}

	// Create new user
	user := User{
		ID:         generateID(),
		SecretCode: generateSecretCode(),
		Name:       name,
		Email:      email,
		Complaints: []Complaint{},
	}

	s.Users[user.ID] = user
	return user, nil
}

// GetUserBySecret retrieves a user by their secret code
func (s *Storage) GetUserBySecret(secretCode string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.Users {
		if user.SecretCode == secretCode {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("user not found")
}

// AddComplaint adds a new complaint to the storage
func (s *Storage) AddComplaint(userID, title, summary string, rating int) (Complaint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user exists
	user, exists := s.Users[userID]
	if !exists {
		return Complaint{}, fmt.Errorf("user not found")
	}

	// Create new complaint
	complaint := Complaint{
		ID:       generateID(),
		Title:    title,
		Summary:  summary,
		Rating:   rating,
		Resolved: false,
		UserID:   userID,
		Date:     time.Now(),
	}

	s.Complaints[complaint.ID] = complaint

	// Add complaint to user's list
	user.Complaints = append(user.Complaints, complaint)
	s.Users[userID] = user

	return complaint, nil
}

// GetComplaint retrieves a complaint by ID
func (s *Storage) GetComplaint(complaintID string) (Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	complaint, exists := s.Complaints[complaintID]
	if !exists {
		return Complaint{}, fmt.Errorf("complaint not found")
	}

	return complaint, nil
}

// GetUserComplaints retrieves all complaints for a specific user
func (s *Storage) GetUserComplaints(userID string) ([]Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.Users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user.Complaints, nil
}

// GetAllComplaints retrieves all complaints in the system
func (s *Storage) GetAllComplaints() ([]Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	complaints := make([]Complaint, 0, len(s.Complaints))
	for _, complaint := range s.Complaints {
		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

// ResolveComplaint marks a complaint as resolved
func (s *Storage) ResolveComplaint(complaintID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	complaint, exists := s.Complaints[complaintID]
	if !exists {
		return fmt.Errorf("complaint not found")
	}

	complaint.Resolved = true
	s.Complaints[complaintID] = complaint

	return nil
}

// IsAdmin checks if the provided secret matches the admin secret
func (s *Storage) IsAdmin(secret string) bool {
	return secret == s.AdminSecret
}