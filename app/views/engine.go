package views

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"main/app/auth"
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
