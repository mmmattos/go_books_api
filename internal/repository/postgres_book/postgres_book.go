package postgres_book

import (
	"database/sql"

	"github.com/mmmattos/books_api/internal/domain"
)

// PostgresBookRepo placeholder for real DB implementation.
type PostgresBookRepo struct {
	DB *sql.DB
}

func NewPostgresBookRepo(db *sql.DB) *PostgresBookRepo {
	return &PostgresBookRepo{DB: db}
}

func (r *PostgresBookRepo) Create(b *domain.Book) error {
	// implement INSERT
	return nil
}

func (r *PostgresBookRepo) GetAll() ([]*domain.Book, error) {
	// implement SELECT
	return nil, nil
}

func (r *PostgresBookRepo) GetByID(id string) (*domain.Book, error) {
	// implement SELECT WHERE id=$1
	return nil, nil
}

func (r *PostgresBookRepo) Update(b *domain.Book) error {
	// implement UPDATE
	return nil
}

func (r *PostgresBookRepo) Delete(id string) error {
	// implement DELETE
	return nil
}
