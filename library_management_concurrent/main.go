package main

import (
	"library_management/controllers"
	"library_management/concurrency"
	"library_management/models"
	"library_management/services"
)

func main() {
	library := services.NewLibrary()

	library.Members[1] = models.Member{ID: 1, Name: "Abel"}
	library.Members[2] = models.Member{ID: 2, Name: "Meron"}
	library.Members[3] = models.Member{ID: 3, Name: "Alice"}
	library.Members[4] = models.Member{ID: 4, Name: "Bob"}
	library.Members[5] = models.Member{ID: 5, Name: "Charlie"}

	testBook := models.Book{ID: 1, Title: "Concurrent Programming in Go", Author: "Test Author"}
	library.AddBook(testBook)

	concurrency.StartWorkers(library, 4)

	controller := controllers.NewLibraryController(library)
	controller.RunConsole()
}