package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"main/app/database"
	"main/app/models"
	"strconv"
)

func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	//Create
	headers := c.GetReqHeaders()

	if !CheckAuth(database.Db, headers["username"], 12, 'r') {
		return c.SendStatus(401)
	}

	if err := c.BodyParser(&task); err != nil {
		fmt.Print(err)
		return c.SendStatus(500)
	}

	//stmt := database.Db.Session(&gorm.Session{DryRun: true}).Create(&task).Statement.SQL.String()
	//fmt.Print(stmt)

	database.Db.Create(&task)

	return c.SendStatus(200)
}

func GetTask(c *fiber.Ctx) error {
	//READ
	var tasks []models.Task
	headers := c.GetReqHeaders()
	if !CheckAuth(database.Db, headers["username"], 12, 'r') {
		return c.SendStatus(401)
	}

	//var payload models.User

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return c.SendStatus(400)
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
	headers := c.GetReqHeaders()

	if !CheckAuth(database.Db, headers["username"], 12, 'r') {
		return c.SendStatus(401)
	}

	if err := c.BodyParser(&task); err != nil {
		return c.SendStatus(500)
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
	headers := c.GetReqHeaders()

	if !CheckAuth(database.Db, headers["username"], 12, 'w') {
		return c.SendStatus(401)
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		fmt.Print(err)
		return c.SendStatus(400)
	}

	if err := database.Db.Where("id = ?", id).Delete(&task); err != nil {
		fmt.Print(err)
		return c.SendStatus(400)
	}

	return c.SendStatus(200)
}
