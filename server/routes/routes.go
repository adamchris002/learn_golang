package routes

import (
	"net/http"

	"backend/handlers"
)

func SetupRoutes() {
	//user related
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/verifyTokenAndLogin", handlers.VerifyTokenAndLogin)

	//task related
	http.HandleFunc("/addTasks", handlers.AddTask)
	http.HandleFunc("/todaysTasks", handlers.CallTodaysTasks)
	http.HandleFunc("/updateTaskCompletion", handlers.UpdateTaskCompletion)
	http.HandleFunc("/deleteTask", handlers.DeleteTask)
	http.HandleFunc("/updateTaskValues", handlers.UpdateTaskValues)
	http.HandleFunc("/deleteSubTask", handlers.DeleteSubtask)
}
