package memory_book

import (
	"errors"
	"fmt"
	"sync"

	"github.com/user/bookapi/internal/domain"
)

// In-memory threadsafe repository.
type MemoryBookRepo struct {
	mu     sync.RWMutex
	data   map[string]*domain.Book
	nextID int
}

func NewMemoryBookRepo() *MemoryBookRepo {
	return &MemoryBookRepo{data: make(map[string]*domain.Book), nextID: 1}
}

func (r *MemoryBookRepo) nextIDStr() string {
	id := r.nextID
	r.nextID++
	return fmt.Sprintf("%d", id)
}

func (r *MemoryBookRepo) Create(b *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b.ID == "" {
		b.ID = r.nextIDStr()
	}
	cp := *b
	r.data[b.ID] = &cp
	return nil
}

func (r *MemoryBookRepo) GetAll() ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Book, 0, len(r.data))
	for _, v := range r.data {
		cpy := *v
		out = append(out, &cpy)
	}
	return out, nil
}

func (r *MemoryBookRepo) GetByID(id string) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cpy := *b
	return &cpy, nil
}

func (r *MemoryBookRepo) Update(b *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.data[b.ID]
	if !ok {
		return errors.New("not found")
	}
	cpy := *b
	r.data[b.ID] = &cpy
	return nil
}

func (r *MemoryBookRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.data[id]
	if !ok {
		return errors.New("not found")
	}
	delete(r.data, id)
	return nil
}
