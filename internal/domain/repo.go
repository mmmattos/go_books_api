package domain

// Repository defines persistence operations for books.
type BookRepository interface {
	Create(b *Book) error
	GetAll() ([]*Book, error)
	GetByID(id string) (*Book, error)
	Update(b *Book) error
	Delete(id string) error
}
