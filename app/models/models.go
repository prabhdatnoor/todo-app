package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique"`
	Pfp       string
	Password  string
	IsAdmin   bool      `gorm:"default:false"`
	LastLogin time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Task struct {
	gorm.Model
	Assignee    uint `gorm:"foreignKey:users_Username"`
	Creator     uint `gorm:"foreignKey:users_Username"`
	Description string
	Name        string
	Status      int
}

type Search struct {
	//Assignee []string
	//NotAssignee bool //apply NOT to users
	/*Time1        time.Time
	Time2        time.Time
	TimeInterval string //one of 'before', 'after', 'in interval', 'out of interval'*/
	//Status      string
	ID string `json:"id"`
	/*DeletedAt   gorm.DeletedAt `gorm:"index"`
	Description string
	Name        string*/
}

type StoreVal struct {
	IsAdmin  bool
	Username string
	ID       uint
	Pfp      string
}
