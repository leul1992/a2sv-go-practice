package router

import (
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
	"task_manager/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", middleware.AuthMiddleWare(), controllers.GetTasks)
	r.GET("/tasks/:id", middleware.AuthMiddleWare(), controllers.GetTask)
	r.POST("/tasks", middleware.AuthMiddleWare(), middleware.IsAdmin(), controllers.CreateTask)
	r.PUT("/tasks/:id", middleware.AuthMiddleWare(), middleware.IsAdmin(), controllers.UpdateTask)
	r.DELETE("/tasks/:id", middleware.AuthMiddleWare(), middleware.IsAdmin(), controllers.DeleteTask)
	// User authentication routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/promote/:id", middleware.AuthMiddleWare(), middleware.IsAdmin(), controllers.PromoteUser)
	return r
}