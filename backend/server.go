package main

import (
	"strings"
	"time"

	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type User struct {
	gorm.Model
	id       uint   `gorm:"primaryKey;autoIncrement:true"`
	username string `gorm:"unique"`
}
type Task struct {
	gorm.Model
	id          uint `gorm:"primaryKey;autoIncrement:true"`
	assignee    User
	description string
	name        string
	createdTime int64 `gorm:"serializer:unixtime;type:time"`
	status      string
}

//type DBINFO struct {
//	port     uint64
//	host     string
//	user     string
//	password string
//	name     string
//}

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
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Task{})

	// Create
	db.Create(&User{username: "communi"})
	var user User

	db.First(&user, "username = ?", "communi")

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

	app.Get("/tasks/:userID/:taskID?", func(c *fiber.Ctx) error {
		if c.Params("userID") == "" {
			return fiber.NewError(401, "No userID provided!")
		}
		return c.SendString("Hello " + c.Params("userID") + " " + c.Params("taskID") + "!")
		// => Hello john

		//return c.SendString("Where is john?")
	})

	app.Get("pfp/:userID/", func(c *fiber.Ctx) error {
		if c.Params("userID") == "" {
			return fiber.NewError(401, "No userID provided!")
		}

		return c.SendFile(os.Getenv("LOC") + "./bruh.png")
	})

	app.Listen(":" + os.Getenv("PORT"))
}
