package routes

import (
	"errors"
	db "mmm/reminder/db"
	model "mmm/reminder/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateTaskInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type UpdateTaskInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

func ReadTask(c *gin.Context) {
	DB, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	var tasks []model.Task
	r := DB.Find(&tasks)

	if r.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errors.New("Failed to retrieve Tasks."),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}

func CreateTask(c *gin.Context) {
	DB, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	var createTaskInput CreateTaskInput
	if err := c.ShouldBindJSON(&createTaskInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	task := model.Task{Title: createTaskInput.Title, Description: createTaskInput.Description, Timestamp: createTaskInput.Timestamp}
	r := DB.Create(&task)

	if r.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errors.New("Failed to insert new Task."),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

func UpdateTask(c *gin.Context) {
	DB, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	id := c.Param("id")
	var updateTaskInput UpdateTaskInput
	var task model.Task

	if err := DB.Where(id).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": errors.New("Task not found."),
		})
		return
	}

	if err := c.ShouldBindJSON(&updateTaskInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	toUpdateTask := model.Task{Title: updateTaskInput.Title, Description: updateTaskInput.Description, Timestamp: updateTaskInput.Timestamp}
	r := DB.Model(&task).Updates(toUpdateTask)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errors.New("Failed to updated task."),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

func DeleteTask(c *gin.Context) {
	DB, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	id := c.Param("id")
	var task model.Task

	if err := DB.Where(id).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": errors.New("Task not found."),
		})
		return
	}

	r := DB.Delete(&task)
	if r.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errors.New("Failed to delete Task."),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "Task deleted successfully.",
	})
}
