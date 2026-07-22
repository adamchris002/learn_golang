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

func CallTodaysTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task

	userID := r.URL.Query().Get("userId")

	now := time.Now()

	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	today := time.Now().Format("2006-01-02")

	endOfDay := startOfDay.Add(24 * time.Hour)

	result := database.DB.
		Preload("Subtasks").
		Where("tasks.user_id = ?", userID).
		Where(`
		(tasks.created_at >= ? AND tasks.created_at < ?)
		OR
		(tasks.due_date != '' AND tasks.due_date >= ?)
	`, startOfDay, endOfDay, today).
		Find(&tasks)

	if result.Error != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
