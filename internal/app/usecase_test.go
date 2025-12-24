package app_test

import (
	"testing"

	"github.com/mmmattos/books_api/internal/app"
	"github.com/mmmattos/books_api/internal/domain"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

func newUsecase() *app.Usecase {
	repo := memory_book.NewMemoryBookRepo()
	return app.NewUsecase(repo)
}

func TestCreateBook_Success(t *testing.T) {
	uc := newUsecase()

	book := &domain.Book{
		ID:     "1",
		Title:  "DDD",
		Author: "Evans",
	}

	if err := uc.CreateBook(book); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAllBooks(t *testing.T) {
	uc := newUsecase()

	_ = uc.CreateBook(&domain.Book{
		ID:     "1",
		Title:  "DDD",
		Author: "Evans",
	})

	books, err := uc.GetAllBooks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
}

func TestGetBookByID(t *testing.T) {
	uc := newUsecase()

	book := &domain.Book{
		ID:     "1",
		Title:  "DDD",
		Author: "Evans",
	}
	_ = uc.CreateBook(book)

	got, err := uc.GetBookByID("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != "1" {
		t.Fatalf("expected ID 1, got %s", got.ID)
	}
}

func TestUpdateBook(t *testing.T) {
	uc := newUsecase()

	book := &domain.Book{
		ID:     "1",
		Title:  "Old",
		Author: "A",
	}
	_ = uc.CreateBook(book)

	book.Title = "New"
	if err := uc.UpdateBook(book); err != nil {
		t.Fatalf("unexpected error updating book: %v", err)
	}
}

func TestDeleteBook(t *testing.T) {
	uc := newUsecase()

	_ = uc.CreateBook(&domain.Book{
		ID:     "1",
		Title:  "DDD",
		Author: "Evans",
	})

	if err := uc.DeleteBook("1"); err != nil {
		t.Fatalf("unexpected error deleting book: %v", err)
	}
}
