package postgres_book

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/mmmattos/books_api/internal/domain"
)

// TEMPORARY in-memory implementation to satisfy repository contract.
// Keeps SQL-compatible constructor for main.go and tests.

type PostgresBookRepo struct {
	mu    sync.Mutex
	books map[string]*domain.Book
}

// IMPORTANT: keep *sql.DB parameter for compatibility
func NewPostgresBookRepo(_ *sql.DB) *PostgresBookRepo {
	return &PostgresBookRepo{
		books: make(map[string]*domain.Book),
	}
}

func (r *PostgresBookRepo) Create(b *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[b.ID] = b
	return nil
}

func (r *PostgresBookRepo) GetAll() ([]*domain.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]*domain.Book, 0, len(r.books))
	for _, b := range r.books {
		out = append(out, b)
	}
	return out, nil
}

func (r *PostgresBookRepo) GetByID(id string) (*domain.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	b, ok := r.books[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (r *PostgresBookRepo) Update(b *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[b.ID]; !ok {
		return errors.New("not found")
	}
	r.books[b.ID] = b
	return nil
}

func (r *PostgresBookRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return errors.New("not found")
	}
	delete(r.books, id)
	return nil
}
