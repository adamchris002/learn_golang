package routes

import (
	"net/http"

	"backend/handlers"
)

func SetupRoutes() {
	//user related
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)

	//task related
	http.HandleFunc("/addTasks", handlers.AddTask)
	http.HandleFunc("/todaysTasks", handlers.CallTodaysTasks)
	http.HandleFunc("/updateTaskCompletion", handlers.UpdateTaskCompletion)
}
