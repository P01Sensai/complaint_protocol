package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func NewStorage() *Storage {
	return &Storage{
		Users:       make(map[string]User),
		Complaints:  make(map[string]Complaint),
		AdminSecret: "admin123",
	}
}


func generateID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}


func generateSecretCode() string {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
	
		return fmt.Sprintf("sc%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}


func (s *Storage) AddUser(name, email string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	
	for _, user := range s.Users {
		if user.Email == email {
			return User{}, fmt.Errorf("email already registered")
		}
	}

	
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


func (s *Storage) AddComplaint(userID, title, summary string, rating int) (Complaint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()


	user, exists := s.Users[userID]
	if !exists {
		return Complaint{}, fmt.Errorf("user not found")
	}


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

	
	user.Complaints = append(user.Complaints, complaint)
	s.Users[userID] = user

	return complaint, nil
}

func (s *Storage) GetComplaint(complaintID string) (Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	complaint, exists := s.Complaints[complaintID]
	if !exists {
		return Complaint{}, fmt.Errorf("complaint not found")
	}

	return complaint, nil
}


func (s *Storage) GetUserComplaints(userID string) ([]Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.Users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user.Complaints, nil
}

func (s *Storage) GetAllComplaints() ([]Complaint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	complaints := make([]Complaint, 0, len(s.Complaints))
	for _, complaint := range s.Complaints {
		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

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


func (s *Storage) IsAdmin(secret string) bool {
	return secret == s.AdminSecret
}