package main

import (
	"log"

	"task-manager/Delivery/controllers"
	"task-manager/Delivery/routers"
	"task-manager/Infrastructure"
	"task-manager/Repositories"
	"task-manager/Usecases"
)

func main() {
	if err := Infrastructure.InitDB(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	taskRepo := Repositories.NewTaskRepository()
	userRepo := Repositories.NewUserRepository()

	jwtService := Infrastructure.NewJWTService()
	passwordService := Infrastructure.NewPasswordService()

	taskUseCase := Usecases.NewTaskUseCase(taskRepo)
	userUseCase := Usecases.NewUserUseCase(userRepo, passwordService, jwtService)

	controllers.InitControllers(taskUseCase, userUseCase)

	r := routers.SetupRouter()

	log.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Println("Failed to start server:", err)
	}
}
