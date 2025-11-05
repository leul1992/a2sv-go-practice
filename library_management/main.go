package main

import (
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

func main() {
	library := services.NewLibrary()

	library.Members[1] = models.Member{ID: 1, Name: "Abel"}
	library.Members[2] = models.Member{ID: 2, Name: "Meron"}

	controller := controllers.NewLibraryController(library)
	controller.RunConsole()
}
