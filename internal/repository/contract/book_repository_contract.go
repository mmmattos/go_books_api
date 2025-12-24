package contract

import (
	"testing"

	"github.com/mmmattos/books_api/internal/domain"
)

// BookRepositoryContract is a shared test helper.
// It must live in a non _test.go file so other packages can import it.
func BookRepositoryContract(
	t *testing.T,
	newRepo func(t *testing.T) domain.BookRepository,
) {
	t.Helper()

	t.Run("create and list books", func(t *testing.T) {
		repo := newRepo(t)

		book := &domain.Book{
			ID:     "1",
			Title:  "Domain-Driven Design",
			Author: "Eric Evans",
		}

		if err := repo.Create(book); err != nil {
			t.Fatalf("create failed: %v", err)
		}

		books, err := repo.GetAll()
		if err != nil {
			t.Fatalf("get all failed: %v", err)
		}

		if len(books) != 1 {
			t.Fatalf("expected 1 book, got %d", len(books))
		}
	})

	t.Run("get by id", func(t *testing.T) {
		repo := newRepo(t)

		book := &domain.Book{
			ID:     "42",
			Title:  "Clean Architecture",
			Author: "Robert C. Martin",
		}

		_ = repo.Create(book)

		got, err := repo.GetByID("42")
		if err != nil {
			t.Fatalf("get by id failed: %v", err)
		}

		if got.ID != "42" {
			t.Fatalf("unexpected book id: %s", got.ID)
		}
	})

	t.Run("delete removes book", func(t *testing.T) {
		repo := newRepo(t)

		book := &domain.Book{
			ID:     "99",
			Title:  "Refactoring",
			Author: "Martin Fowler",
		}

		_ = repo.Create(book)

		if err := repo.Delete("99"); err != nil {
			t.Fatalf("delete failed: %v", err)
		}

		if _, err := repo.GetByID("99"); err == nil {
			t.Fatalf("expected error after delete")
		}
	})
}
