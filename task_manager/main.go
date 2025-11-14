package main

import "task_manager/router"

func main() {
	r := router.SetupRouter()
	println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		println("Failed to start server:", err)
	}
}
