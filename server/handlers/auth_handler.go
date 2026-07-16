package handlers

import (
	"encoding/json"
	"net/http"

	"backend/database"
	"backend/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	database.DB.Create(&user)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login endpoint",
	})
}
