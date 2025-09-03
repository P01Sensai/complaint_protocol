package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handlers struct {
	storage *Storage
}

func NewHandlers(storage *Storage) *Handlers {
	return &Handlers{storage: storage}
}


func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	
	var request struct {
		SecretCode string `json:"secret_code"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.storage.GetUserBySecret(request.SecretCode)
	if err != nil {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{User: user})
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.Name) == "" || strings.TrimSpace(request.Email) == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	user, err := h.storage.AddUser(request.Name, request.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{User: user})
}

func (h *Handlers) SubmitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	user, err := h.storage.GetUserBySecret(secret)
	if err != nil {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	var request struct {
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Rating  int    `json:"rating"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if request.Rating < 1 || request.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	complaint, err := h.storage.AddComplaint(user.ID, request.Title, request.Summary, request.Rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ComplaintResponse{Complaint: complaint})
}

func (h *Handlers) GetAllComplaintsForUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	user, err := h.storage.GetUserBySecret(secret)
	if err != nil {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	complaints, err := h.storage.GetUserComplaints(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ComplaintsListResponse{Complaints: complaints})
}

func (h *Handlers) GetAllComplaintsForAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	if !h.storage.IsAdmin(secret) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	complaints, err := h.storage.GetAllComplaints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ComplaintsListResponse{Complaints: complaints})
}

func (h *Handlers) ViewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	complaintID := r.URL.Query().Get("id")
	if complaintID == "" {
		http.Error(w, "Complaint ID required", http.StatusBadRequest)
		return
	}

	complaint, err := h.storage.GetComplaint(complaintID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	user, userErr := h.storage.GetUserBySecret(secret)
	isAdmin := h.storage.IsAdmin(secret)
	isOwner := userErr == nil && user.ID == complaint.UserID

	if !isAdmin && !isOwner {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ComplaintResponse{Complaint: complaint})
}

func (h *Handlers) ResolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := r.Header.Get("X-Secret-Code")
	if secret == "" {
		http.Error(w, "Secret code required", http.StatusUnauthorized)
		return
	}

	if !h.storage.IsAdmin(secret) {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}

	var request struct {
		ComplaintID string `json:"complaint_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.storage.ResolveComplaint(request.ComplaintID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Complaint resolved successfully"}`))
}