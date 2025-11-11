package services

import (
	"errors"
	"fmt"
	"library_management/models"
	"sync"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ReserveBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type ReservationRequest struct {
	BookID   int
	MemberID int
	Response chan<- error
}

type Library struct {
	Books          map[int]models.Book
	Members        map[int]models.Member
	Mu             sync.Mutex
	ReservationChan chan ReservationRequest
}

func NewLibrary() *Library {
	return &Library{
		Books:          make(map[int]models.Book),
		Members:        make(map[int]models.Member),
		ReservationChan: make(chan ReservationRequest, 100),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	book.Status = "Available"
	book.ReservedBy = 0
	l.Books[book.ID] = book
	fmt.Printf("\nAdded a new book named %s \n", book.Title)
}

func (l *Library) RemoveBook(bookID int) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if _, ok := l.Books[bookID]; ok {
		delete(l.Books, bookID)
		fmt.Println("\nBook removed successfully.")
	} else {
		fmt.Println("\nBook not found.")
	}
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	if _, ok := l.Members[memberID]; !ok {
		return errors.New("\nmember not found")
	}

	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("\nbook not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("\nbook is already borrowed")
	}
	if book.Status == "Reserved" && book.ReservedBy != memberID {
		return errors.New("\nbook reserved by another member")
	}

	book.Status = "Borrowed"
	book.ReservedBy = 0
	l.Books[bookID] = book

	member := l.Members[memberID]
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	fmt.Printf("\n%s borrowed by %s\n", book.Title, member.Name)
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("\nmember not found")
	}

	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("\nbook not found")
	}

	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("\nbook not borrowed by this member")
	}

	book.Status = "Available"
	book.ReservedBy = 0
	l.Books[bookID] = book
	l.Members[memberID] = member

	fmt.Printf("\n%s returned by %s\n", book.Title, member.Name)
	return nil
}

func (l *Library) ReserveBook(bookID int, memberID int) error {
	resp := make(chan error, 1)
	req := ReservationRequest{
		BookID:   bookID,
		MemberID: memberID,
		Response: resp,
	}
	l.ReservationChan <- req
	return <-resp
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	var available []models.Book
	for _, b := range l.Books {
		if b.Status == "Available" {
			available = append(available, b)
		}
	}
	return available
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	if member, ok := l.Members[memberID]; ok {
		books := make([]models.Book, len(member.BorrowedBooks))
		copy(books, member.BorrowedBooks)
		return books
	}
	return []models.Book{}
}
