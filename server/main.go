package main

import (
	
	"fmt"
	"net/http"

	"backend/database"
	"backend/routes"
	"backend/middleware"

	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()

    if err != nil {
        fmt.Println("No .env file found, using environment variables")
    }

    database.Connect()

    routes.SetupRoutes()

    fmt.Println("Server running on port 8080")

    handler := middleware.CORS(http.DefaultServeMux)

    http.ListenAndServe(":8080", handler)
}
