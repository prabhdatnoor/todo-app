package main

import (
	"github.com/gofiber/fiber/v2"
	"main/app/auth"
	. "main/app/cache"
	"main/app/controllers"
	. "main/app/database"
	"time"

	"fmt"
	"main/app/models"
	"os"

	"encoding/json"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	var err error
	Db, err = ConnectDB()
	if err != nil {
		fmt.Print(err)
		panic("can't connect to Db")
	}

	Store, err = ConnectStore()
	if err != nil {
		fmt.Print(err)
		panic("can't connect to Redis")
	}

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeZone:     "America/New_York",
		TimeInterval: time.Millisecond,
	}))

	// Or extend your config for customization
	app.Use(favicon.New(favicon.Config{
		File: "static/public/favicon.ico",
	}))

	// Or extend your config for customization
	limiterConfig := limiter.Config{
		Max: 20, // max count of connections
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendFile("static/public/toofast.html") // called when a request hits the limit
		},
	}

	app.Use(
		limiter.New(limiterConfig), // add Limiter middleware with config
		//for some reason, it doesnt want to use it
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/logout", auth.Logout)
	app.Post("/login", auth.Login)
	app.Post("/register", auth.Register)
	//app.Post("/verifySess", auth.VerifySessf)

	app.Get("/users/info/:username", func(c *fiber.Ctx) error {
		var user models.User
		var tasks []models.Task

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
