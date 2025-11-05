package services

import (
	"errors"
	"fmt"
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	book.Status = "Available"
	l.Books[book.ID] = book
	fmt.Printf("\nAdded a new book named %s \n", book.Title)
}

func (l *Library) RemoveBook(bookID int) {
	if _, ok := l.Books[bookID]; ok {
		delete(l.Books, bookID)
		fmt.Println("\nBook removed successfully.")
	} else {
		fmt.Println("\nBook not found.")
	}
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("\nbook not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("\nbook is already borrowed")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("\nmember not found")
	}

	book.Status = "Borrowed"
	l.Books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	fmt.Printf("\n%s borrowed by %s\n", book.Title, member.Name)
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("\nmember not found")
	}

	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("\nbook not found")
	}

	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			book.Status = "Available"
			l.Books[bookID] = book
			l.Members[memberID] = member
			fmt.Printf("\n%s returned by %s\n", book.Title, member.Name)
			return nil
		}
	}
	return errors.New("\nbook not borrowed by this member")
}

func (l *Library) ListAvailableBooks() []models.Book {
	var available []models.Book
	for _, b := range l.Books {
		if b.Status == "Available" {
			available = append(available, b)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	if member, ok := l.Members[memberID]; ok {
		return member.BorrowedBooks
	}
	return []models.Book{}
}
