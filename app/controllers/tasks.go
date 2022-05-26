package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"main/app/auth"
	"main/app/database"
	"main/app/models"
	"strconv"
)

func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	//Create

	if err := c.BodyParser(&task); err != nil {
		fmt.Print(err)
		return c.SendStatus(500)
	}

	goodCreds, err := auth.CheckAuth(c, task.ID, 'r')
	if err != nil || !goodCreds {
		return c.SendStatus(401)
	}

	//stmt := database.Db.Session(&gorm.Session{DryRun: true}).Create(&task).Statement.SQL.String()
	//fmt.Print(stmt)

	database.Db.Create(&task)

	return c.SendStatus(200)
}

func GetTask(c *fiber.Ctx) error {
	//READ
	var tasks []models.Task

	//var payload models.User

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
	}

	goodCreds, err := auth.CheckAuth(c, uint(id), 'r')
	if err != nil || !goodCreds {
		fmt.Print(err)
		return c.SendStatus(401)
	}

	var count int64

	res := database.Db.Limit(-1).Where("id=?", id).Find(&tasks).Count(&count)
	//fmt.Print(db.Session(&gorm.Session{DryRun: true}).Limit(-1).Where(models.Task{}, &task).Find(&tasks).Count(&count).Statement.SQL.String())
	if res.Error != nil {
		fmt.Print(res.Error)
		return c.SendStatus(400)
	}

	if len(tasks) == 0 {
		fmt.Print("no found")
		return c.Status(404).SendString("none found")
	}

	return c.Status(200).JSON(&tasks)
}

func PatchTask(c *fiber.Ctx) error {
	//UPDATE
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.SendStatus(500)
	}
	goodCreds, err := auth.CheckAuth(c, task.ID, 'w')
	if err != nil || !goodCreds {
		return c.SendStatus(401)
	}

	if err := database.Db.Model(&task).Updates(&task).Error; err != nil {
		fmt.Print(err)
		return c.Status(fiber.StatusInternalServerError).SendString("error in finding task: " + c.Params("id"))
	}

	return c.SendStatus(200)
}

func DeleteTask(c *fiber.Ctx) error {
	//delete
	var task models.Task

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		fmt.Print(err)
		return c.SendStatus(400)
	}

	goodCreds, err := auth.CheckAuth(c, uint(id), 'r')
	if err != nil || !goodCreds {
		return c.SendStatus(401)
	}

	if err := database.Db.Where("id = ?", id).Delete(&task); err != nil {
		fmt.Print(err)
		return c.SendStatus(400)
	}

	return c.SendStatus(200)
}
