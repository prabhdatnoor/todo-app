package auth

import (
	"encoding/json"
	_ "encoding/json"
	"errors"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/prabhdatnoor/todo-app/app/cache"
	"github.com/prabhdatnoor/todo-app/app/database"
	"github.com/prabhdatnoor/todo-app/app/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"unicode"
)

type APIUser struct {
	Password string
	ID       uint
}

var (
	Def_passhash string = `$argon2id$v=19$m=65536,t=1,p=2$tlc5gPMfmKK5HRToLjpxCA$Pq0RBqRoZJdL07r/wpQnwWY91nK1zrjxIkOF8aB005g`
	//output of argon2id.CreateHash("guest", argon2id.DefaultParams)
)

func Bool2String(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
func StoreVal2Json(user *models.User) string {
	fmt.Print(user)
	return "{Username:" + user.Username + ", IsAdmin:" + Bool2String(user.IsAdmin) + ", ID:" + fmt.Sprint(user.ID) + "}"

}

type isAdmin struct {
	isAdmin bool
}

func CheckAuth(c *fiber.Ctx, taskID uint, r rune) (models.AuthCheck, error) {
	var auth models.AuthCheck

	deets, err := VerifySess(c)
	auth.Creds = deets

	if err != nil {
		auth.Success = false
		return auth, err
	}

	if deets.Username == "guest" {
		auth.Success = false
		return auth, nil
	}

	var is isAdmin
	if erro := database.Db.Table("users").Where("username=?", deets.Username).Find(&is).Error; erro != nil {
		auth.Success = false
		return auth, erro
	}

	if is.isAdmin {
		auth.Success = false
		return auth, nil
	}
	var count int64

	if erro := database.Db.Table("tasks").Where("creator=?", deets.ID).Where("id=?", taskID).Count(&count).Error; err != nil {
		auth.Success = false
		return auth, erro
	}

	if count > 0 {
		auth.Success = true
		return auth, nil
	}

	if r == 'r' {
		auth.Success = true
		return auth, nil
	}

	auth.Success = false

	return auth, nil

}

func VerifySess(c *fiber.Ctx) (models.StoreVal, error) {
	sess, err := cache.Store.Get(c)

	if err != nil {
		panic(err)
	}

	val := sess.Get("session:" + sess.ID())

	s := fmt.Sprintf("%s", val)

	var details models.StoreVal
	if err := json.Unmarshal([]byte(s), &details); err != nil {
		return models.StoreVal{Username: "guest"}, err
	}

	return details, nil
}

func VerifyUser(uname, password string) (uint, bool, error) {
	var passhash string
	//var user models.User

	var pid APIUser

	var count int64

	err := database.Db.Table("users").Where("username = ?", uname).Find(&pid).Count(&count).Error
	//fmt.Print(db.Session(&gorm.Session{DryRun: true}).Model(&user).Select("password, id").Where("username = ?", uname).Find(&user).Statement.SQL.String())
	if err != nil || count == 0 {
		passhash = Def_passhash
		fmt.Print(err)
	} else {
		//fmt.Print(pid)

		passhash = pid.Password
	}

	//fmt.Print(pid)

	match, erro := argon2id.ComparePasswordAndHash(password, passhash)
	if erro != nil {
		fmt.Print(erro)
		if count == 0 {
			return 0, false, gorm.ErrRecordNotFound
		}
		return 0, false, erro
	}

	return pid.ID, match, nil
}

func Login(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		goodCreds, erro := VerifySess(c)

		if goodCreds.Username != "guest" && erro == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"User": fiber.Map{
					"username": goodCreds.Username,
					"id":       goodCreds.ID,
					"pfp":      goodCreds.Pfp,
				},
				"success": true,
				"message": "Content de te revoir, " + goodCreds.Username + "!",
			})
		}

		fmt.Print(erro)
		return c.Status(fiber.StatusInternalServerError).SendString("server failure in login")

	}

	id, match, err := VerifyUser(user.Username, user.Password)

	if !match {
		fmt.Print(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": true,
				"message": "username doesn't exist",
			})
		}
		return c.Status(fiber.StatusForbidden).SendString("failure")
	}

	sess, err := cache.Store.Get(c)
	if err != nil {
		panic(err)
	}

	//fmt.Print(user)

	//var storeval StoreVal
	//storeval.IsAdmin = user.IsAdmin
	//storeval.Username = user.Username
	user.ID = id
	b, err := json.Marshal(&user)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("failure in create session")
	}
	sess.Set("session:"+sess.ID(), string(b))

	if err := sess.Save(); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"User": fiber.Map{
			"username": user.Username,
			"id":       id,
			"pfp":      user.Pfp,
		},
		"success": true,
		"message": "Bienvenue, " + user.Username + "!",
	})
}

func Logout(c *fiber.Ctx) error {
	sess, err := cache.Store.Get(c)

	if err != nil {
		panic(err)
	}

	if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "logout success",
	})
}

func CreateUser(user *models.User) error {
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		fmt.Print(err)
		return fiber.ErrInternalServerError

	}

	user.Password = hash

	if err := database.Db.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Create(user).Error; err != nil {
		fmt.Print(err)
		return fiber.ErrBadRequest
	}
	return nil
}

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		fmt.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"User": fiber.Map{
				"username": user.Username,
			},
			"success": false,
			"message": "Error in parsing body, likely bad request",
		})
	}

	if err := CreateUser(&user); err != nil {
		if errors.Is(err, fiber.ErrInternalServerError) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"User": fiber.Map{
					"username": user.Username,
				},
				"success": false,
				"message": "Error in password hashing",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"User": fiber.Map{
				"username": user.Username,
				"id":       user.ID,
			},
			"success": false,
			"message": "Username already exists or invalid",
		})

	}

	firstChar := []rune(user.Username[0:1])[0]
	isLetter := unicode.IsLetter(firstChar)
	if user.Username == "guest" || !isLetter {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"User": fiber.Map{
				"username": user.Username,
			},
			"success": false,
			"message": `Can't use 'guest' as username and username must start with an alphabet`,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"User": fiber.Map{
			"username": user.Username,
			"id":       user.ID,
		},
		"success": true,
		"message": "User , " + user.Username + "is created. Welcome!",
	})
}
