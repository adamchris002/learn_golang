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
	http.HandleFunc("/incompleteTasks", handlers.CallAllIncompleteTasks)
	http.HandleFunc("/todaysTasks", handlers.CallTodaysTasks)
	http.HandleFunc("/oneWeekTasks", handlers.CallOneWeekTasks)
	http.HandleFunc("/updateTaskCompletion", handlers.UpdateTaskCompletion)
	http.HandleFunc("/deleteTask", handlers.DeleteTask)
	http.HandleFunc("/updateTaskValues", handlers.UpdateTaskValues)
	http.HandleFunc("/deleteSubTask", handlers.DeleteSubtask)
	http.HandleFunc("/updateTaskStartDate", handlers.UpdateTaskStartDate)
	// http.HandleFunc("/changeTaskToPending", handlers.ChangeTaskToPending)
	http.HandleFunc("/changeTaskToActive", handlers.ChangeTaskToActive)
}
