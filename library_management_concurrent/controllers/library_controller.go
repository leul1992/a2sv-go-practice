package controllers

import (
	"fmt"
	"library_management/models"
	"library_management/services"
)

type LibraryController struct {
	service services.LibraryManager
}

func NewLibraryController(service services.LibraryManager) *LibraryController {
	return &LibraryController{service: service}
}

func (lc *LibraryController) RunConsole() {
	for {
		fmt.Println("\nLibrary Management System")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Reserve Book")
		fmt.Println("8. Simulate Concurrent Reservations")
		fmt.Println("9. Exit")

		var choice int
		fmt.Print("Enter choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var id int
			var title, author string
			fmt.Print("Enter book ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter title: ")
			fmt.Scan(&title)
			fmt.Print("Enter author: ")
			fmt.Scan(&author)
			lc.service.AddBook(models.Book{ID: id, Title: title, Author: author})

		case 2:
			var id int
			fmt.Print("Enter book ID to remove: ")
			fmt.Scan(&id)
			lc.service.RemoveBook(id)

		case 3:
			var bookID, memberID int
			fmt.Print("Enter book ID: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter member ID: ")
			fmt.Scan(&memberID)
			err := lc.service.BorrowBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			}

		case 4:
			var bookID, memberID int
			fmt.Print("Enter book ID: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter member ID: ")
			fmt.Scan(&memberID)
			err := lc.service.ReturnBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			}

		case 5:
			books := lc.service.ListAvailableBooks()
			fmt.Println("\nAvailable Books:")
			for _, b := range books {
				fmt.Printf("ID: %d | %s by %s (Reserved by: %d)\n", b.ID, b.Title, b.Author, b.ReservedBy)
			}

		case 6:
			var memberID int
			fmt.Print("Enter member ID: ")
			fmt.Scan(&memberID)
			books := lc.service.ListBorrowedBooks(memberID)
			fmt.Println("\nBorrowed Books:")
			for _, b := range books {
				fmt.Printf("ID: %d | %s by %s\n", b.ID, b.Title, b.Author)
			}

		case 7:
			var bookID, memberID int
			fmt.Print("Enter book ID: ")
			fmt.Scan(&bookID)
			fmt.Print("Enter member ID: ")
			fmt.Scan(&memberID)
			err := lc.service.ReserveBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book reserved successfully. You have 5 seconds to borrow it.")
			}

		case 8:
			var bookID int
			fmt.Print("Enter book ID for simulation: ")
			fmt.Scan(&bookID)
			fmt.Println("Simulating 5 concurrent reservation attempts by members 1-5")
			for i := 1; i <= 5; i++ {
				memberID := i
				go func(m int) {
					err := lc.service.ReserveBook(bookID, m)
					if err != nil {
						fmt.Printf("Member %d: %v\n", m, err)
					} else {
						fmt.Printf("Member %d reserved book %d successfully\n", m, bookID)
					}
				}(memberID)
			}
			fmt.Println("Launched concurrent simulations. Check console for async outputs.")

		case 9:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option. Try again.")
		}
	}
}
