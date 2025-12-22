package app_test

import (
	"testing"

	"github.com/mmmattos/books_api/internal/app"
	"github.com/mmmattos/books_api/internal/domain"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

func setupUsecase() *app.Usecase {
	repo := memory_book.NewMemoryBookRepo()
	return app.NewUsecase(repo)
}

func TestCreateBook_Success(t *testing.T) {
	uc := setupUsecase()

	book := &domain.Book{
		Title:  "Clean Architecture",
		Author: "Robert C. Martin",
	}

	if err := uc.CreateBook(book); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if book.ID == "" {
		t.Fatalf("expected book ID to be set")
	}
}

func TestCreateBook_ValidationError(t *testing.T) {
	uc := setupUsecase()

	book := &domain.Book{
		Title: "",
	}

	if err := uc.CreateBook(book); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestGetAllBooks(t *testing.T) {
	uc := setupUsecase()

	book := &domain.Book{
		Title:  "DDD",
		Author: "Eric Evans",
	}
	_ = uc.CreateBook(book)

	books, err := uc.GetAllBooks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
}

func TestGetBookByID(t *testing.T) {
	uc := setupUsecase()

	book := &domain.Book{
		Title:  "Refactoring",
		Author: "Martin Fowler",
	}
	_ = uc.CreateBook(book)

	found, err := uc.GetBookByID(book.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Title != book.Title {
		t.Fatalf("expected title %q, got %q", book.Title, found.Title)
	}
}

func TestGetBookByID_NotFound(t *testing.T) {
	uc := setupUsecase()

	_, err := uc.GetBookByID("non-existent-id")
	if err == nil {
		t.Fatalf("expected error for missing book")
	}
}

func TestUpdateBook(t *testing.T) {
	uc := setupUsecase()

	book := &domain.Book{
		Title:  "Old Title",
		Author: "Someone",
	}
	_ = uc.CreateBook(book)

	book.Title = "New Title"

	if err := uc.UpdateBook(book); err != nil {
		t.Fatalf("unexpected error updating book: %v", err)
	}

	updated, _ := uc.GetBookByID(book.ID)
	if updated.Title != "New Title" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}
}

func TestDeleteBook(t *testing.T) {
	uc := setupUsecase()

	book := &domain.Book{
		Title:  "To Be Deleted",
		Author: "Anon",
	}
	_ = uc.CreateBook(book)

	if err := uc.DeleteBook(book.ID); err != nil {
		t.Fatalf("unexpected error deleting book: %v", err)
	}

	_, err := uc.GetBookByID(book.ID)
	if err == nil {
		t.Fatalf("expected error after deletion")
	}
}
