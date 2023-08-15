package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTask struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type UpdateTask struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
	)

	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_paulo",
		host,
		user,
		password,
		dbname,
		port,
	)

	db, err := gorm.Open(postgres.Open(conn))
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(Task{})

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	})

	r.GET("/tasks", func(c *gin.Context) {
		var tasks []Task
		r := db.Find(&tasks)

		if r.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "Failed to retrieve Tasks.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": tasks,
		})
	})

	r.POST("/tasks", func(c *gin.Context) {
		var createTask CreateTask
		if err := c.ShouldBindJSON(&createTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		task := Task{Title: createTask.Title, Description: createTask.Description, Timestamp: createTask.Timestamp}
		r := db.Create(&task)

		if r.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "Failed to insert new Task.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": task,
		})
	})

	r.PUT("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updateTask UpdateTask
		var task Task

		if err := db.Where(id).First(&task).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"err": "Task not found",
			})
			return
		}

		if err := c.ShouldBindJSON(&updateTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		toUpdateTask := Task{Title: updateTask.Title, Description: updateTask.Description, Timestamp: updateTask.Timestamp}
		r := db.Model(&task).Updates(toUpdateTask)
		if r.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": task,
		})
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		var task Task

		if err := db.Where(id).First(&task).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"err": "Task not found",
			})
			return
		}

		r := db.Delete(&task)
		if r.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "Failed to delete task",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": "Task deleted successfully.",
		})
	})

	r.Run()
}
