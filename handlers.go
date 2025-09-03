package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Handlers struct holds the storage instance
type Handlers struct {
	storage *Storage
}

// NewHandlers creates a new Handlers instance
func NewHandlers(storage *Storage) *Handlers {
	return &Handlers{storage: storage}
}

// LoginHandler handles user login
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request struct {
		SecretCode string `json:"secret_code"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get user by secret code
	user, err := h.storage.GetUserBySecret(request.SecretCode)
	if err != nil {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	// Return user details (without sensitive information in a real application)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{User: user})
}

// RegisterHandler handles user registration
func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate input
	if strings.TrimSpace(request.Name) == "" || strings.TrimSpace(request.Email) == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	// Add user to storage
	user, err := h.storage.AddUser(request.Name, request.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return user details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{User: user})
}

// SubmitComplaintHandler handles complaint submission
func (h *Handlers) SubmitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user secret from header
	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	// Get user by secret code
	user, err := h.storage.GetUserBySecret(secret)
	if err != nil {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var request struct {
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Rating  int    `json:"rating"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate input
	if strings.TrimSpace(request.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if request.Rating < 1 || request.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	// Add complaint to storage
	complaint, err := h.storage.AddComplaint(user.ID, request.Title, request.Summary, request.Rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return complaint details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ComplaintResponse{Complaint: complaint})
}

// GetAllComplaintsForUserHandler returns all complaints for the authenticated user
func (h *Handlers) GetAllComplaintsForUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user secret from header
	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	// Get user by secret code
	user, err := h.storage.GetUserBySecret(secret)
	if err != nil {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	// Get user's complaints
	complaints, err := h.storage.GetUserComplaints(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return complaints list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ComplaintsListResponse{Complaints: complaints})
}

// GetAllComplaintsForAdminHandler returns all complaints for admin review
func (h *Handlers) GetAllComplaintsForAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get admin secret from header
	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	// Check if secret matches admin secret
	if !h.storage.IsAdmin(secret) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	// Get all complaints
	complaints, err := h.storage.GetAllComplaints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return complaints list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ComplaintsListResponse{Complaints: complaints})
}

// ViewComplaintHandler returns a specific complaint
func (h *Handlers) ViewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get complaint ID from query parameter
	complaintID := r.URL.Query().Get("id")
	if complaintID == "" {
		http.Error(w, "Complaint ID required", http.StatusBadRequest)
		return
	}

	// Get complaint from storage
	complaint, err := h.storage.GetComplaint(complaintID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get user secret from header
	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	// Check if user is admin or the owner of the complaint
	user, userErr := h.storage.GetUserBySecret(secret)
	isAdmin := h.storage.IsAdmin(secret)
	isOwner := userErr == nil && user.ID == complaint.UserID

	if !isAdmin && !isOwner {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Return complaint details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ComplaintResponse{Complaint: complaint})
}

// ResolveComplaintHandler marks a complaint as resolved (admin only)
func (h *Handlers) ResolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get admin secret from header
	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	// Check if secret matches admin secret
	if !h.storage.IsAdmin(secret) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	// Parse request body
	var request struct {
		ComplaintID string `json:"complaint_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Resolve the complaint
	err := h.storage.ResolveComplaint(request.ComplaintID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Complaint resolved successfully"}`))
}