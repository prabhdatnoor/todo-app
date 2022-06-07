package main

import (
	"github.com/gofiber/fiber/v2"
	"main/app/auth"
	. "main/app/cache"
	"main/app/controllers"
	. "main/app/database"
	"main/app/views"
	"time"

	"fmt"
	"os"

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

	views.SetupEngine()
	app := fiber.New(fiber.Config{Views: views.Engine})
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

	app.Get("/logout", auth.Logout)
	app.Post("/login", auth.Login)
	app.Post("/api/register", auth.Register)

	app.Post("/tasks", controllers.CreateTask)

	app.Get("/tasks/:id", controllers.GetTask)

	app.Patch("/tasks", controllers.PatchTask)

	app.Delete("/tasks", controllers.DeleteTask)

	app.Static("/", "./static/public")
	app.Static("/pfps", "./static/pfps")
	app.Static("/register", "./static/public/register.html")
	app.Get("/task/:id/edit", views.RenderEditTask)
	app.Get("/task/:id", views.RenderTask)

	app.Get("/home", views.RenderUser)

	if app.Listen(":"+os.Getenv("PORT")) != nil {
		fmt.Print("app listening ERROR!")
	}
}
