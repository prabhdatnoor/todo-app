package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	creds, err := auth.CheckAuth(c, task.ID, 'r')
	if err != nil || !creds.Success {
		return c.SendStatus(401)
	}

	//stmt := database.Db.Session(&gorm.Session{DryRun: true}).Create(&task).Statement.SQL.String()
	//fmt.Print(stmt)
	task.Creator = creds.Creds.ID
	database.Db.Create(&task)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "task created",
		"task":    task,
	})
}

func AllTasks(c *fiber.Ctx, tasks *[]models.Task) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return err
	}

	creds, err := auth.CheckAuth(c, uint(id), 'r')

	if err != nil || !creds.Success {
		return err
	}

	var count int64

	res := database.Db.Limit(-1).Find(tasks).Count(&count)
	//fmt.Print(db.Session(&gorm.Session{DryRun: true}).Limit(-1).Where(models.Task{}, &task).Find(&tasks).Count(&count).Statement.SQL.String())
	if res.Error != nil {
		return res.Error
	}

	if len(*tasks) == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func Tasks(c *fiber.Ctx, tasks *[]models.Task) error {

	//var payload models.User

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)

	if err != nil {
		return err
	}

	creds, err := auth.CheckAuth(c, uint(id), 'r')
	if err != nil || !creds.Success {
		return err
	}

	var count int64

	res := database.Db.Limit(-1).Where("id=?", id).Find(tasks).Count(&count)
	//fmt.Print(db.Session(&gorm.Session{DryRun: true}).Limit(-1).Where(models.Task{}, &task).Find(&tasks).Count(&count).Statement.SQL.String())
	if res.Error != nil {
		return res.Error
	}

	if len(*tasks) == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func TaskbyID(c *fiber.Ctx, task *models.Task) error {

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	var def_task models.Task

	if err != nil {
		return err
	}

	creds, err := auth.CheckAuth(c, uint(id), 'r')
	if err != nil || !creds.Success {
		return err
	}

	res := database.Db.Table("tasks").Limit(-1).Where("id=?", id).Find(&task)
	//fmt.Print(db.Session(&gorm.Session{DryRun: true}).Limit(-1).Where(models.Task{}, &task).Find(&tasks).Count(&count).Statement.SQL.String())
	if res.Error != nil {
		return res.Error
	}

	if *task == def_task {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func GetTask(c *fiber.Ctx) error {
	var tasks []models.Task

	if c.Params("id") == "all" {
		err := AllTasks(c, &tasks)

		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"tasks":   tasks,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}

	if err := Tasks(c, &tasks); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err,
			"success": false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"tasks":   tasks,
	})
}

func PatchTask(c *fiber.Ctx) error {
	//UPDATE
	var task models.TaskSave

	if err := c.BodyParser(&task); err != nil {
		return c.SendStatus(500)
	}
	creds, err := auth.CheckAuth(c, task.ID, 'w')
	if err != nil || !creds.Success {
		return c.SendStatus(401)
	}

	if err := database.Db.Table("tasks").Model(&task).Updates(&task).Error; err != nil {
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

	creds, err := auth.CheckAuth(c, uint(id), 'r')
	if err != nil || !creds.Success {
		return c.SendStatus(401)
	}

	if err := database.Db.Where("id = ?", id).Delete(&task); err != nil {
		fmt.Print(err)
		return c.SendStatus(400)
	}

	return c.SendStatus(200)
}
