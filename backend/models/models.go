package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       uint   `gorm:"primaryKey;autoIncrement:true"`
	Username string `gorm:"unique"`
}
type Task struct {
	gorm.Model
	Id          uint `gorm:"primaryKey;autoIncrement:true"`
	Assignee    uint `gorm:"foreignKey:users_Id"`
	Description string
	Name        string
	Status      string
}
