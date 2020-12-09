package main

import (
	router "github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "todo/docs"
	h "todo/handlers"
)

// @title Swagger TODO
// @version 1.0
// @description This is an api for todo apps.
// @termsOfService http://swagger.io/terms/

// @contact.name Kurian
// @contact.url http://www.shopalyst.com/support
// @contact.email kurianc@shopalyst.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	app := router.Default()                                                //gin mux
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //swagger documentation

	app.GET("/task", h.GetAllHandler)            // get task by id
	app.GET("/task/:id", h.GetTaskHandler)       // get task by id
	app.POST("/task", h.CreateTaskHandler)       // create task
	app.PUT("/task/:id", h.UpdateTaskHandler)    // update task
	app.DELETE("/task/:id", h.DeleteTaskHandler) // delete task
	app.POST("/subtask", h.SubTaskHandler)       // create subtask
	app.Run(":7070")
}
