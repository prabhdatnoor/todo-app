package views

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"main/app/auth"
	"main/app/controllers"
	"main/app/database"
	"main/app/models"
	"main/app/utils"
)

var (
	Engine *html.Engine
)

func SetupEngine() {
	Engine = html.New("resources", ".html").Reload(true).Debug(true)
}

func isZero(n int) bool {
	return n == 0
}

func RenderUser(c *fiber.Ctx) error {
	var tasks []models.Task

	goodCreds, err := auth.VerifySess(c)
	if err != nil || goodCreds.Username == "guest" {
		return c.SendStatus(401)
	}
	id := goodCreds.ID

	if erro := database.Db.Table("tasks").Limit(20).Where("assignee=?", id).Or("creator=?", id).Find(&tasks).Error; erro != nil {
		fmt.Print(erro)
	}

	role := "User"
	if goodCreds.IsAdmin {
		role = "Admin"
	}

	username := goodCreds.Username

	return c.Render(
		"home", fiber.Map{
			"Username": username,
			"Pfp":      utils.GetPfp(username),
			"Role":     role,
			"Tasks":    tasks,
		})
}

/*func RenderTasks(c *fiber.Ctx) error {
	var tasks *[]models.Task

	if err := controllers.Tasks(c, tasks); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err,
			"success": false,
		})
	}

	goodCreds, err := auth.VerifySess(c)
	if err != nil || goodCreds.Username == "guest" {
		return c.SendStatus(401)
	}
	id := goodCreds.ID

	role := "User"
	if goodCreds.IsAdmin {
		role = "Admin"
	}

	return c.Render(
		"tasks", fiber.Map{
			"Username": goodCreds.Username,
			"Role":     role,
			"ID":       id,
			"Tasks":    *tasks,
		})
}*/

func RenderTask(c *fiber.Ctx) error {
	var task models.Task

	if err := controllers.TaskbyID(c, &task); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err,
			"success": false,
		})
	}

	goodCreds, err := auth.VerifySess(c)
	if err != nil || goodCreds.Username == "guest" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err,
			"success": false,
		})
	}
	id := goodCreds.ID
	role := "User"
	if goodCreds.IsAdmin {
		role = "Admin"
	}

	return c.Render(
		"task", fiber.Map{
			"Username": goodCreds.Username,
			"Role":     role,
			"ID":       id,
			"Task":     task,
		})

}

func RenderEditTask(c *fiber.Ctx) error {
	var task models.Task

	goodCreds, err := auth.VerifySess(c)
	if err != nil || goodCreds.Username == "guest" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err,
			"success": false,
		})
	}
	id := goodCreds.ID
	role := "User"
	if goodCreds.IsAdmin {
		role = "Admin"
	}
	f := c.AllParams()

	fmt.Print(f)

	if c.Params("id") == "new" {
		return c.Render(
			"edit_task", fiber.Map{
				"Username": goodCreds.Username,
				"Role":     role,
				"ID":       id,
				"Task":     task,
				"isPost":   true,
			})
	} else {
		if err := controllers.TaskbyID(c, &task); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": err.Error(),
				"success": false,
			})
		}
	}

	return c.Render(
		"edit_task", fiber.Map{
			"Username": goodCreds.Username,
			"Role":     role,
			"ID":       id,
			"Task":     task,
			"isPost":   false,
		})
}
