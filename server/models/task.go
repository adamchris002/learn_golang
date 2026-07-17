package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Completed   bool   `json:"completed"`

	UserID uint `json:"userId"`
	User   User `json:"user"`

	Subtasks []Subtask `json:"subtasks"`
}