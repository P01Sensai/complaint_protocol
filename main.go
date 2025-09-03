package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	
	storage := NewStorage()
	handlers := NewHandlers(storage)
	

	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/submitComplaint", handlers.SubmitComplaintHandler)
	http.HandleFunc("/getAllComplaintsForUser", handlers.GetAllComplaintsForUserHandler)
	http.HandleFunc("/getAllComplaintsForAdmin", handlers.GetAllComplaintsForAdminHandler)
	http.HandleFunc("/viewComplaint", handlers.ViewComplaintHandler)
	http.HandleFunc("/resolveComplaint", handlers.ResolveComplaintHandler)
	
	
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}