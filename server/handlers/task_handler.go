package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/models"
)

type SubtaskRequest struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type UpdateTaskRequest struct {
	DueDate     string           `json:"dueDate"`
	Description string           `json:"description"`
	SubTask     []SubtaskRequest `json:"subTask"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	database.DB.Create(&task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Add Task Success",
		"message":      "Task added",
	})
}

func CallTasksBasedOnWeek(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	weekDateStr := r.URL.Query().Get("weekDate")
	if weekDateStr == "" {
		http.Error(w, "Missing week date", http.StatusBadRequest)
		return
	}

	layout := "02/01/2006"

	translateWeekToTime, err := time.Parse(layout, weekDateStr)
	if err != nil {
		http.Error(w, "Invalid week date format", http.StatusBadRequest)
		return
	}

	weekday := int(translateWeekToTime.Weekday())
	daysSinceMonday := (weekday + 6) % 7

	startOfWeek := time.Date(translateWeekToTime.Year(), translateWeekToTime.Month(), translateWeekToTime.Day(), 0, 0, 0, 0, translateWeekToTime.Location()).
		AddDate(0, 0, -daysSinceMonday)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var tasks []models.Task
	result := database.DB.Preload("Subtasks").
		Where("user_id = ? AND TO_DATE(due_date, 'DD/MM/YYYY') >= ? AND TO_DATE(due_date, 'DD/MM/YYYY') < ?", userID, startOfWeek, endOfWeek).
		Find(&tasks)

	if result.Error != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func CallTaskBasedOnMonth(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	layout := "02/01/2006"

	monthDateStr := r.URL.Query().Get("monthDate")
	if monthDateStr == "" {
		http.Error(w, "Missing month date", http.StatusBadRequest)
		return
	}

	translateMonthToTime, err := time.Parse(layout, monthDateStr)
	if err != nil {
		http.Error(w, "Invalid month date format", http.StatusBadRequest)
		return
	}

	startOfMonth := time.Date(translateMonthToTime.Year(), translateMonthToTime.Month(), 1, 0, 0, 0, 0, translateMonthToTime.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var tasks []models.Task

	result := database.DB.Preload("Subtasks").Where("tasks.user_id = ? AND TO_DATE(tasks.due_date, 'DD/MM/YYYY') >= ? AND TO_DATE(tasks.due_date, 'DD/MM/YYYY') < ?", userID, startOfMonth, endOfMonth).Find(&tasks)

	if result.Error != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func CallAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	userID := r.URL.Query().Get("userId")

	result := database.DB.
		Preload("Subtasks").
		Where("tasks.user_id = ?", userID).
		Find(&tasks)

	if result.Error != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func CallAllIncompleteTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	userID := r.URL.Query().Get("userId")

	result := database.DB.
		Preload("Subtasks").
		Where("tasks.user_id = ? AND tasks.completed = false", userID).
		Find(&tasks)

	if result.Error != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func CallOneWeekTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	userID := r.URL.Query().Get("userId")

	result := database.DB.
		Preload("Subtasks").
		Where("tasks.user_id = ?", userID).
		Where(`
        TO_DATE(tasks.task_start, 'DD/MM/YYYY') >= ?
        AND
        TO_DATE(tasks.task_start, 'DD/MM/YYYY') <= ?
    `, time.Now(), time.Now().AddDate(0, 0, 7)).
		Find(&tasks)

	if result.Error != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func CallTodaysTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	userID := r.URL.Query().Get("userId")

	today := time.Now()

	result := database.DB.
		Preload("Subtasks").
		Where("tasks.user_id = ?", userID).
		Where(`
			TO_DATE(tasks.task_start, 'DD/MM/YYYY') = ?
			OR
			(tasks.due_date != '' AND TO_DATE(tasks.due_date, 'DD/MM/YYYY') >= ? AND TO_DATE(tasks.task_start, 'DD/MM/YYYY') <= ?)
		`, today, today, today).
		Find(&tasks)

	if result.Error != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func UpdateTaskCompletion(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || userID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid User ID",
		})
		return
	}

	taskId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || taskId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Task ID",
		})
		return
	}

	data, err := strconv.ParseBool(r.URL.Query().Get("data"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid completion value",
		})
		return
	}

	result := database.DB.Model(&models.Task{}).Where("id = ?", taskId).Where("user_id = ? ", userID).Update("completed", data)

	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      result.Error.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Update Task Success",
		"message":      fmt.Sprintf("Task with ID %d has been updated", taskId),
	})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || userID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid User ID",
		})
		return
	}

	taskId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || taskId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Task ID",
		})
		return
	}

	result := database.DB.
		Where("id = ? AND user_id = ?", taskId, userID).
		Delete(&models.Task{})

	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Delete Task Failed",
			"message":      result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]string{
			"mesageTitle": "Delete Task Fail",
			"message":     "Task not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Delete Task Success",
		"message":      fmt.Sprintf("Task with ID %d has been delete", taskId),
	})
}

func UpdateTaskValues(w http.ResponseWriter, r *http.Request) {
	var updateTaskRequest UpdateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&updateTaskRequest)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || userID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid User ID",
		})
		return
	}

	taskId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || taskId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Task ID",
		})
		return
	}

	result := database.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userID).Updates(models.Task{Description: updateTaskRequest.Description, DueDate: updateTaskRequest.DueDate})

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Task not found",
		})
		return
	}

	for _, subtask := range updateTaskRequest.SubTask {
		if subtask.ID == 0 {
			result := database.DB.Create(&models.Subtask{
				TaskID:    uint(taskId),
				Title:     subtask.Title,
				Completed: subtask.Completed,
			})

			if result.Error != nil {
				// handle error
				return
			}
		} else {
			result := database.DB.
				Model(&models.Subtask{}).
				Where("id = ? AND task_id = ?", subtask.ID, uint(taskId)).
				Select("Title", "Completed").
				Updates(models.Subtask{
					Title:     subtask.Title,
					Completed: subtask.Completed,
				})
			if result.Error != nil {
				// handle error
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Update Task Success",
		"message":      fmt.Sprintf("Task with ID %d has been updated", taskId),
	})

}

func DeleteSubtask(w http.ResponseWriter, r *http.Request) {
	subTaskID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || subTaskID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Sub-Task ID",
		})
		return
	}

	taskId, err := strconv.Atoi(r.URL.Query().Get("taskId"))
	if err != nil || taskId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Task ID",
		})
		return
	}

	result := database.DB.Where("id = ? AND task_id = ?", subTaskID, taskId).Delete(&models.Subtask{})

	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Delete Sub-Task Failed",
			"message":      result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]string{
			"mesageTitle": "Delete Sub-Task Fail",
			"message":     "Sub-Task not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Delete Sub-Task Success",
		"message":      fmt.Sprintf("Sub-Task with id of %d has been deleted", taskId),
	})
}

func UpdateTaskStartDate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		StartDate string `json:"task_start"`
	}

	taskId, err := strconv.Atoi(r.URL.Query().Get("taskId"))
	if err != nil || taskId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Task ID",
		})
	}
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid User ID",
		})
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid request body",
		})
		return
	}

	newStartDate, err := time.Parse("02/01/2006 15:04", body.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid start date format",
		})
		return
	}

	var task models.Task

	getExistingData := database.DB.
		Where("id = ? AND user_id = ?", taskId, userId).
		First(&task)

	if getExistingData.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      getExistingData.Error.Error(),
		})
		return
	}

	updateData := map[string]interface{}{
		"task_start": body.StartDate,
	}

	if task.DueDate != "" {
		dueDate, err := time.Parse("02/01/2006 15:04", task.DueDate)

		if err == nil && dueDate.Before(newStartDate) {
			updateData["due_date"] = ""
		}
	}

	result := database.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userId).Updates(updateData)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      result.Error.Error(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Update Task Success",
		"message":      fmt.Sprintf("Task with ID %d has been updated", taskId),
	})
}

// func ChangeTaskToPending(w http.ResponseWriter, r *http.Request) {
// 	taskId, err := strconv.Atoi(r.URL.Query().Get("taskId"))
// 	if err != nil || taskId <= 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"messageTitle": "Update Task Failed",
// 			"message":      "Invalid Task ID",
// 		})
// 	}
// 	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
// 	if err != nil || userId <= 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"messageTitle": "Update Task Failed",
// 			"message":      "Invalid User ID",
// 		})
// 	}

// 	result := database.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userId).Update("due_date", "")

// 	if result.Error != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"messageTitle": "Update Task Failed",
// 			"message":      result.Error.Error(),
// 		})
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	json.NewEncoder(w).Encode(map[string]string{
// 		"messageTitle": "Update Task Success",
// 		"message":      fmt.Sprintf("Task with ID %d has been updated", taskId),
// 	})
// }

func ChangeTaskToActive(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DueDate string `json:"due_date"`
	}

	now := time.Now()
	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	).Format("02/01/2006 15:04")

	taskId, err := strconv.Atoi(r.URL.Query().Get("taskId"))
	if err != nil || taskId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Task ID",
		})
	}
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid User ID",
		})
	}

	json.NewDecoder(r.Body).Decode(&body)

	result := database.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userId).Select("TaskStart", "DueDate").Updates(models.Task{TaskStart: today, DueDate: body.DueDate})

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      result.Error.Error(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Update Task Success",
		"message":      fmt.Sprintf("Task with ID %d has been updated", taskId),
	})
}

func ChangeTaskDueDate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NewDueDate string `json:"date_data"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid request body",
		})
		return
	}

	taskId, err := strconv.Atoi(r.URL.Query().Get("taskId"))
	if err != nil || taskId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid Task ID",
		})
	}
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      "Invalid User ID",
		})
	}

	result := database.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskId, userId).Update("due_date", body.NewDueDate)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"messageTitle": "Update Task Failed",
			"message":      result.Error.Error(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"messageTitle": "Update Task Success",
		"message":      fmt.Sprintf("Task with ID %d has been updated", taskId),
	})
}
