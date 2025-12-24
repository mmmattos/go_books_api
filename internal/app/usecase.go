package app

import (
	"fmt"
	"time"

	"github.com/mmmattos/books_api/internal/domain"
)

// Usecase coordinates domain operations.
type Usecase struct {
	Repo domain.BookRepository
}

func NewUsecase(r domain.BookRepository) *Usecase {
	return &Usecase{Repo: r}
}

func (u *Usecase) CreateBook(b *domain.Book) error {
	if b == nil {
		return fmt.Errorf("book is nil")
	}
	if b.Title == "" || b.Author == "" {
		return fmt.Errorf("title and author are required")
	}

	// FIX: assign ID without external dependency
	if b.ID == "" {
		b.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return u.Repo.Create(b)
}

func (u *Usecase) GetAllBooks() ([]*domain.Book, error) {
	return u.Repo.GetAll()
}

func (u *Usecase) GetBookByID(id string) (*domain.Book, error) {
	return u.Repo.GetByID(id)
}

func (u *Usecase) UpdateBook(b *domain.Book) error {
	if b == nil {
		return fmt.Errorf("book is nil")
	}
	return u.Repo.Update(b)
}

func (u *Usecase) DeleteBook(id string) error {
	return u.Repo.Delete(id)
}
