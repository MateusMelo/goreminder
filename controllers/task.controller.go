package controllers

import (
	"fmt"
	"mmm/reminder/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

func NewTaskController(DB *gorm.DB) TaskController {
	return TaskController{DB}
}

func (c *TaskController) ReadTasks(ctx *gin.Context) {
	var tasks []models.Task

	results := c.DB.Find(&tasks)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": results.Error,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	var payload *models.CreateTaskRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	task := models.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Timestamp:   payload.Timestamp,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := c.DB.Create(&task)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": task})
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	fmt.Printf(id)

	var payload *models.UpdateTask
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}

	var updatedTask models.Task
	result := c.DB.First(&updatedTask, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})
		return
	}
	now := time.Now()
	taskToUpdate := models.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Timestamp:   payload.Timestamp,
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   now,
	}

	c.DB.Model(&updatedTask).Updates(taskToUpdate)

	ctx.JSON(http.StatusOK, gin.H{
		"data": updatedTask,
	})
}

func (c *TaskController) ReadTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var task models.Task
	result := c.DB.First(&task, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": task})
}

func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Delete(&models.Task{}, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
