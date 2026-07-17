package models

import "gorm.io/gorm"

type Subtask struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`

	TaskID uint `json:"taskId"`
	Task   Task `json:"task"`
}