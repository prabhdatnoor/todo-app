package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Pfp      string
}
type Task struct {
	gorm.Model
	Assignee    uint `gorm:"foreignKey:users_Username"`
	Creator     uint `gorm:"foreignKey:users_Username"`
	Description string
	Name        string
	Status      string
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
