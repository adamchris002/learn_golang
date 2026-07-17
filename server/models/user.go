package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Tasks []Task `json:"tasks"`
}