package main

import (
	"backend/models"
	"strconv"
	"strings"
	"time"

	"fmt"
	"os"

	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//type DBINFO struct {
//	port     uint64
//	host     string
//	user     string
//	password string
//	name     string
//}

func checkAuth(db *gorm.DB, uname string, taskID uint, role rune) bool {
	return true
}

func main() {
	fmt.Print(os.Getenv("POSTGRES_PASSWORD_FILE"))
	dat, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		panic("failed to read password file")
	}

	port := "port=5432"
	host := "host=" + os.Getenv("DATABASE_HOST")
	userq := "user=" + os.Getenv("DATABASE_USER")
	password := "password=" + strings.TrimSuffix(string(dat), "\n")
	password = "password=postgres"
	name := "database=" + os.Getenv("DATABASE_DB")

	dsn := port + " " + host + " " + userq + " " + password + " " + name + " sslmode=disable TimeZone=America/New_York"
	fmt.Print(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema

	if db.AutoMigrate(&models.User{}) != nil {
		fmt.Print("failed to migrate users")
	}
	if db.AutoMigrate(&models.Task{}) != nil {
		fmt.Print("failed to migrate tasks")
	}

	// Create
	db.Create(&models.User{Username: "communitech"})
	var user models.User
	var task models.Task
	var tasks []models.Task
	//var search models.Search

	db.First(&user, "Username = ?", "communi")
	//fmt.Print(user)
	db.First(&user, "Username = ?", "carmenwinstead")
	//fmt.Print(user)

	//db.Create(&models.Task{Assignee: user.ID, Name: "what da dog doin?", Status: "not started", Description: "hello there, quandale dingle here"})

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeZone:     "America/New_York",
		TimeInterval: time.Millisecond,
	}))

	// Provide a minimal config
	app.Use(favicon.New())

	// Or extend your config for customization
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))

	// Default middleware config
	app.Use(limiter.New())

	// Or extend your config for customization
	limiterConfig := limiter.Config{
		Max: 20, // max count of connections
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendFile("./toofast.html") // called when a request hits the limit
		},
	}

	app.Use(
		limiter.New(limiterConfig), // add Limiter middleware with config
		//for some reason, it doesnt want to use it
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	/*app.Get("/tasks/:userID?/:taskID?", func(c *fiber.Ctx) error {
		if c.Params("userID") == "" {
			return fiber.NewError(401, "No userID provided!")
		}
		return c.SendString("Hello " + c.Params("userID") + " " + c.Params("taskID") + "!")
		// => Hello john

		//return c.SendString("Where is john?")
	})*/

	app.Get("/users/info/:username", func(c *fiber.Ctx) error {
		if c.Params("username") == "" {
			fmt.Print(err)
			return fiber.NewError(400, "No userID provided!")
		}

		err := db.First(&user, "Username = ?", c.Params("username"))
		if err != nil {
			fmt.Print(err)
			return fiber.NewError(400, "error in finding user")
		}

		err = db.Select("Id").Limit(-1).Find(&tasks, "Username = ?", c.Params("username"))
		if err != nil {
			fmt.Print(err)
			return fiber.NewError(400, "error in finding user in tasks table")
		}
		fmt.Print("bruh")

		data := map[string]interface{}{
			"pfp":   user.Pfp,
			"tasks": tasks,
			"id":    user.ID,
		}

		out, erro := json.Marshal(data)
		if erro != nil {
			fmt.Print(erro)
			return fiber.NewError(500, "error in stringifying")
		}
		return c.SendString(string(out))
	})

	app.Post("/tasks/create", func(c *fiber.Ctx) error {
		//Create
		headers := c.GetReqHeaders()

		if !checkAuth(db, headers["username"], 12, 'r') {
			return fiber.NewError(401)
		}

		if err := c.BodyParser(&task); err != nil {
			return fiber.NewError(500)
		}

		stmt := db.Session(&gorm.Session{DryRun: true}).Create(&task).Statement.SQL.String()
		fmt.Print(stmt)

		db.Create(&task)

		return c.SendStatus(200)
	})

	app.Post("/tasks/read/:all?", func(c *fiber.Ctx) error {
		//READ
		headers := c.GetReqHeaders()

		if !checkAuth(db, headers["username"], 12, 'r') {
			return fiber.NewError(401)
		}

		//payload := new(models.Search)

		if err := c.BodyParser(&task); err != nil {
			fmt.Print(err)
			return fiber.NewError(500)
		}

		//c.BodyParser(&task)

		var toSearch []string
		toSearch = append(toSearch, "id")
		if !(task.CreatedAt.IsZero()) {
			toSearch = append(toSearch, "CreatedAt")
		}

		if !(task.UpdatedAt.IsZero()) {
			toSearch = append(toSearch, "UpdatedAt")
		}

		if !(task.Assignee == 0) {
			toSearch = append(toSearch, "Assignee")
		}

		if !(task.Creator == 0) {
			toSearch = append(toSearch, "Creator")
		}

		var count int64
		fmt.Print(task.ID, "\n")

		res := db.Limit(-1).Where(models.Task{}, &task).Find(&tasks)
		fmt.Print(db.Session(&gorm.Session{DryRun: true}).Limit(-1).Where(models.Task{}, &task).Find(&tasks).Count(&count).Statement.SQL.String())
		if res.Error != nil {
			fmt.Print(res.Error)
			return fiber.NewError(400)
		}
		if len(tasks) == 0 {
			fmt.Print("no found")
			return fiber.NewError(404, "none found")
		}

		return c.Status(200).JSON(&tasks)
	})

	app.Put("/tasks/edit", func(c *fiber.Ctx) error {
		//UPDATE

		headers := c.GetReqHeaders()

		if !checkAuth(db, headers["username"], 12, 'r') {
			return fiber.NewError(401)
		}

		if err := c.BodyParser(&task); err != nil {
			return fiber.NewError(500)
		}

		if err := db.Model(&task).Updates(&task); err != nil {
			return fiber.NewError(400)
		}

		return c.SendStatus(200)
	})

	app.Delete("/tasks/delete/:id", func(c *fiber.Ctx) error {
		//delete
		headers := c.GetReqHeaders()

		if !checkAuth(db, headers["username"], 12, 'w') {
			return fiber.NewError(401)
		}

		id, err := strconv.ParseUint(c.Params("id"), 10, 32)

		if err != nil {
			return fiber.NewError(400)
		}

		//var toDel models.Task
		//stmt := db.Session(&gorm.Session{DryRun: true}).Where("id = ?", id).Delete(&task).Statement.SQL.String()
		//fmt.Print(stmt)
		if err := db.Where("id = ?", id).Delete(&task); err != nil {
			return fiber.NewError(400)
		}

		return c.SendStatus(200)
	})

	if app.Listen(":"+os.Getenv("PORT")) != nil {
		fmt.Print("app listening ERROR!")
	}
}
