package main

import (
	routes "mmm/reminder/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/tasks", routes.ReadTask)
	r.POST("/tasks", routes.CreateTask)
	r.PUT("/tasks/:id", routes.UpdateTask)
	r.DELETE("/tasks/:id", routes.DeleteTask)

	r.Run()
}
