package routes

import (
	"mmm/reminder/controllers"

	"github.com/gin-gonic/gin"
)

type TaskRouteController struct {
	taskController controllers.TaskController
}

func NewRouteTaskController(taskController controllers.TaskController) TaskRouteController {
	return TaskRouteController{taskController}
}

func (c *TaskRouteController) TaskRoute(rg *gin.RouterGroup) {
	rg.GET("/tasks", c.taskController.ReadTasks)
	rg.POST("/tasks", c.taskController.CreateTask)
	rg.PUT("/tasks/:id", c.taskController.UpdateTask)
	rg.GET("/tasks/:id", c.taskController.ReadTask)
	rg.DELETE("/tasks/:id", c.taskController.DeleteTask)
}
