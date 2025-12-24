package memory_book

import (
	"errors"

	"github.com/mmmattos/books_api/internal/domain"
)

type MemoryBookRepo struct {
	data map[string]*domain.Book
}

func NewMemoryBookRepo() *MemoryBookRepo {
	return &MemoryBookRepo{
		data: make(map[string]*domain.Book),
	}
}

// Alias for tests / convenience
func New() *MemoryBookRepo {
	return NewMemoryBookRepo()
}

func (r *MemoryBookRepo) Create(b *domain.Book) error {
	r.data[b.ID] = b
	return nil
}

func (r *MemoryBookRepo) GetAll() ([]*domain.Book, error) {
	var out []*domain.Book
	for _, b := range r.data {
		out = append(out, b)
	}
	return out, nil
}

func (r *MemoryBookRepo) GetByID(id string) (*domain.Book, error) {
	b, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (r *MemoryBookRepo) Update(b *domain.Book) error {
	if _, ok := r.data[b.ID]; !ok {
		return errors.New("not found")
	}
	r.data[b.ID] = b
	return nil
}

func (r *MemoryBookRepo) Delete(id string) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("not found")
	}
	delete(r.data, id)
	return nil
}
