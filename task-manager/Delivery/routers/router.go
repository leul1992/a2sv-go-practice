package routers

import (
	"task-manager/Delivery/controllers"
	"task-manager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authMiddleware := Infrastructure.AuthMiddleware()
	adminMiddleware := Infrastructure.AdminMiddleware()

	// Task routes
	r.GET("/tasks", authMiddleware, controllers.GetTasks)
	r.GET("/tasks/:id", authMiddleware, controllers.GetTask)
	r.POST("/tasks", authMiddleware, adminMiddleware, controllers.CreateTask)
	r.PUT("/tasks/:id", authMiddleware, adminMiddleware, controllers.UpdateTask)
	r.DELETE("/tasks/:id", authMiddleware, adminMiddleware, controllers.DeleteTask)

	// User routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/promote/:id", authMiddleware, adminMiddleware, controllers.PromoteUser)

	return r
}
