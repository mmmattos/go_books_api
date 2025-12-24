package memory_book_test

import (
	"testing"

	"github.com/mmmattos/books_api/internal/domain"
	"github.com/mmmattos/books_api/internal/repository/contract"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

func TestMemoryBookRepo_Contract(t *testing.T) {
	contract.BookRepositoryContract(t, func(t *testing.T) domain.BookRepository {
		return memory_book.NewMemoryBookRepo()
	})
}
