package memory_book_test

import (
	"testing"

	"github.com/mmmattos/books_api/internal/domain"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

func setupRepo() domain.BookRepository {
	return memory_book.NewMemoryBookRepo()
}

func TestCreateAndGetBook(t *testing.T) {
	repo := setupRepo()

	book := &domain.Book{
		Title:  "Domain-Driven Design",
		Author: "Eric Evans",
	}

	if err := repo.Create(book); err != nil {
		t.Fatalf("unexpected error creating book: %v", err)
	}

	if book.ID == "" {
		t.Fatalf("expected book ID to be set")
	}

	found, err := repo.GetByID(book.ID)
	if err != nil {
		t.Fatalf("unexpected error fetching book: %v", err)
	}

	if found.Title != book.Title {
		t.Fatalf("expected title %q, got %q", book.Title, found.Title)
	}
}

func TestGetAllBooks(t *testing.T) {
	repo := setupRepo()

	_ = repo.Create(&domain.Book{Title: "Book 1", Author: "A"})
	_ = repo.Create(&domain.Book{Title: "Book 2", Author: "B"})

	books, err := repo.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}

func TestUpdateBook(t *testing.T) {
	repo := setupRepo()

	book := &domain.Book{
		Title:  "Old Title",
		Author: "Someone",
	}
	_ = repo.Create(book)

	book.Title = "New Title"

	if err := repo.Update(book); err != nil {
		t.Fatalf("unexpected error updating book: %v", err)
	}

	updated, _ := repo.GetByID(book.ID)
	if updated.Title != "New Title" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}
}

func TestDeleteBook(t *testing.T) {
	repo := setupRepo()

	book := &domain.Book{
		Title:  "To Be Deleted",
		Author: "Anon",
	}
	_ = repo.Create(book)

	if err := repo.Delete(book.ID); err != nil {
		t.Fatalf("unexpected error deleting book: %v", err)
	}

	_, err := repo.GetByID(book.ID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := setupRepo()

	_, err := repo.GetByID("non-existent-id")
	if err == nil {
		t.Fatalf("expected error for missing book")
	}
}
