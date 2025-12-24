package postgres_book_test

import (
	"testing"

	"github.com/mmmattos/books_api/internal/domain"
	"github.com/mmmattos/books_api/internal/repository/contract"
	"github.com/mmmattos/books_api/internal/repository/postgres_book"
)

func TestPostgresBookRepo_Contract(t *testing.T) {
	contract.BookRepositoryContract(t, func(t *testing.T) domain.BookRepository {
		// DB can be nil for contract tests
		return postgres_book.NewPostgresBookRepo(nil)
	})
}
