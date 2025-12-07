package app

import (
	"fmt"

	"github.com/user/bookapi/internal/domain"
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
	return u.Repo.Create(b)
}

func (u *Usecase) GetAllBooks() ([]*domain.Book, error) {
	return u.Repo.GetAll()
}

func (u *Usecase) GetBookByID(id string) (*domain.Book, error) {
	if id == "" {
		return nil, fmt.Errorf("id required")
	}
	return u.Repo.GetByID(id)
}

func (u *Usecase) UpdateBook(b *domain.Book) error {
	if b == nil || b.ID == "" {
		return fmt.Errorf("invalid book")
	}
	return u.Repo.Update(b)
}

func (u *Usecase) DeleteBook(id string) error {
	if id == "" {
		return fmt.Errorf("id required")
	}
	return u.Repo.Delete(id)
}
