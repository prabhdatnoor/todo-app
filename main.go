package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"main/app/controllers"
	. "main/app/controllers"
	. "main/app/database"
	"time"

	"fmt"
	"main/app/models"
	"os"

	"encoding/json"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	var err error
	Db, err = ConnectDB()
	if err != nil {
		panic("can't connect to Db")
	}
	var user models.User
	var tasks []models.Task

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeZone:     "America/New_York",
		TimeInterval: time.Millisecond,
	}))

	// Or extend your config for customization
	app.Use(favicon.New(favicon.Config{
		File: "app/favicon.ico",
	}))

	// Or extend your config for customization
	limiterConfig := limiter.Config{
		Max: 20, // max count of connections
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendFile("app/toofast.html") // called when a request hits the limit
		},
	}

	app.Use(
		limiter.New(limiterConfig), // add Limiter middleware with config
		//for some reason, it doesnt want to use it
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/auth", basicauth.New(basicauth.Config{
		Realm: "Forbidden",
		Authorizer: func(userID, pass string) bool {
			match, err := VerifyUser(Db, userID, pass)
			if err != nil {
				return false
			}
			return match
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}))

	app.Post("auth/login", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&user); err != nil {
			fmt.Print(err)
			return c.SendStatus(500)
		}

		match, err := VerifyUser(Db, user.Username, user.Password)
		if err != nil {
			return c.SendStatus(403)
		}
		if match == true {
			return c.SendStatus(200)
		}

		return c.SendStatus(403)
	})

	app.Post("auth/register", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&user); err != nil {
			fmt.Print(err)
			return c.SendStatus(500)
		}

		if err != nil {
			fmt.Print(err)
			return c.SendStatus(500)
		}
		hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			fmt.Print(err)
			return c.SendStatus(500)
		}

		user.Password = string(hash)

		if erro := Db.Create(&user); err != nil {
			fmt.Print(erro)
			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	})

	app.Get("/users/info/:username", func(c *fiber.Ctx) error {
		if c.Params("username") == "" {
			fmt.Print(err)
			return c.Status(400).SendString("No userID provided!")
		}

		err := Db.First(&user, "Username = ?", c.Params("username"))
		if err != nil {
			fmt.Print(err)
			return c.Status(400).SendString("error in finding user")
		}

		err = Db.Select("Id").Limit(-1).Find(&tasks, "Username = ?", c.Params("username"))
		if err != nil {
			fmt.Print(err)
			return c.Status(400).SendString("error in finding user in tasks table")
		}

		data := map[string]interface{}{
			"pfp":   user.Pfp,
			"tasks": tasks,
			"id":    user.ID,
		}

		out, erro := json.Marshal(data)
		if erro != nil {
			fmt.Print(erro)
			return c.Status(500).SendString("error in stringifying")
		}
		return c.SendString(string(out))
	})

	app.Post("/tasks", controllers.CreateTask)

	app.Get("/tasks/:id", controllers.GetTask)

	app.Patch("/tasks", controllers.PatchTask)

	app.Delete("/tasks", controllers.DeleteTask)

	if app.Listen(":"+os.Getenv("PORT")) != nil {
		fmt.Print("app listening ERROR!")
	}
}
