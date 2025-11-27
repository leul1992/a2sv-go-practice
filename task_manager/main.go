package main

import (
	"log"

	"task_manager/data"
	"task_manager/router"
)

func main() {
	if err := data.InitDB(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	r := router.SetupRouter()
	log.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Println("Failed to start server:", err)
	}
}