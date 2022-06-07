package main

import (
	"bytes"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/prabhdatnoor/todo-app/app/auth"
	"github.com/prabhdatnoor/todo-app/app/cache"
	"github.com/prabhdatnoor/todo-app/app/controllers"
	"github.com/prabhdatnoor/todo-app/app/database"
	"github.com/prabhdatnoor/todo-app/app/models"
	"github.com/prabhdatnoor/todo-app/app/utils"
	"github.com/prabhdatnoor/todo-app/app/views"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

var username, password string
var app fiber.App
var store *session.Store
var db *gorm.DB

func setup() {
	var err error
	db, err = database.ConnectDB()
	if err != nil {
		fmt.Print(err)
		panic("can't connect to Db")
	}

	store, err = cache.ConnectStore()
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

	passPlain, err := utils.GenSalt(12)
	if err != nil {
		panic(err)
	}

	password, err = argon2id.CreateHash(string(passPlain), argon2id.DefaultParams)
	if err != nil {
		panic(err)
	}

	username = "testsubject_" + strconv.Itoa(rand.Intn(100000-2)+2)
	var testUser models.User

	testUser.Username = username
	testUser.Password = password

	if err := auth.CreateUser(&testUser); err != nil {
		log.Printf("error creating test user: %s", err.Error())
	}
}

func shutdown() {
	var user models.User
	if err := database.Db.Delete(&user); err.Error != nil {
		fmt.Print("error deleting test user:")
		fmt.Print(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestWrongUsernamePassword(t *testing.T) {
	data := url.Values{}
	data.Set("Username", username)
	data.Add("password", password)

	req, err := http.NewRequest(
		"POST",
		"/api/login",
		bytes.NewBufferString(data.Encode()),
	)

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}

	res, err := app.Test(req, -1)

	assert.Equalf(t, false, err != nil, "Failed to create correct username password request")
	assert.Equalf(t, 200, res.StatusCode, "Failed checking correct username password")

}

func TestAuthAfterLogout(t *testing.T) {
	req, err := http.NewRequest(
		"POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	res, err := app.Test(req, -1)
	assert.Equalf(t, false, err != nil, "Failed to create logout request")
	assert.Equalf(t, 200, res.StatusCode, "Failed checking correct username password")

	req, err = http.NewRequest(
		"GET", "/home", nil)

	// Perform the request plain with the app.
	// The -1 disables request latency.
	res, err = app.Test(req, -1)

	// verify that no error occured, that is not expected
	assert.Equalf(t, false, err != nil, "Failed to create post logout request")

	// Verify if the status code is as expected
	assert.Equalf(t, 401, res.StatusCode, "Failed to request auth from logged out user")
}
