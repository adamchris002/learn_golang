package routes

import (
	"net/http"

	"backend/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
}
