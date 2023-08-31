package main

import (
	"log"
	"mmm/reminder/controllers"
	"mmm/reminder/initializers"
	"mmm/reminder/routes"

	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	TaskController      controllers.TaskController
	TaskRouteController routes.TaskRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	TaskController = controllers.NewTaskController(initializers.DB)
	TaskRouteController = routes.NewRouteTaskController(TaskController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	router := server.Group("/api")

	TaskRouteController.TaskRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
