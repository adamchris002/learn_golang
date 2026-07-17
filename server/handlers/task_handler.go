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

	endOfDay := startOfDay.Add(24 * time.Hour)

	result := database.DB.Where("user_id = ?", userID).Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).Find(&tasks)

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
