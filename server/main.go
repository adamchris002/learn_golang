package main

import (
	
	"fmt"
	"net/http"

	"backend/database"
	"backend/models"
	"backend/routes"
	"backend/middleware"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Cannot load .env")
	}

	database.Connect()

	database.DB.AutoMigrate(
		&models.User{},
	)
	

	routes.SetupRoutes()
	
	fmt.Println("Server running on port 8080")

	handler := middleware.CORS(http.DefaultServeMux)

	http.ListenAndServe(":8080", handler)
}
