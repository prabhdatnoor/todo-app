package controllers

import (
	"fmt"
	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
	"main/app/models"
)

var (
	Def_passhash string = "$argon2id$v=19$m=65536,t=1,p=2$tlc5gPMfmKK5HRToLjpxCA$Pq0RBqRoZJdL07r/wpQnwWY91nK1zrjxIkOF8aB005g"
)

func CheckAuth(db *gorm.DB, uname string, taskID uint, role rune) bool {
	return true
}

func VerifyUser(Db *gorm.DB, uname, password string) (bool, error) {
	var passhash string
	//fmt.Print(uname, password)
	var user models.User
	err := Db.Model(&user).Select("password").Where("username = ?", uname).Find(&passhash).Error
	if err != nil {
		passhash = Def_passhash
		fmt.Print(err)
	}

	match, erro := argon2id.ComparePasswordAndHash(password, passhash)
	if erro != nil {
		fmt.Print(erro)
		return false, erro
	}

	return match, nil
}
