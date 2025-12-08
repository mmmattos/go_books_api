package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmmattos/books_api/internal/app"
	"github.com/mmmattos/books_api/internal/domain"
	"github.com/mmmattos/books_api/internal/handlers"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

func TestBooksAPI_Integration(t *testing.T) {
	repo := memory_book.NewMemoryBookRepo()
	tuc := app.NewUsecase(repo)
	r := handlers.NewRouter(tuc)
	ts := httptest.NewServer(r)
	defer ts.Close()

	b := domain.Book{Title: "Test", Author: "Me", Year: 2025}
	buf, _ := json.Marshal(b)
	resp, err := stdhttp.Post(ts.URL+"/books", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 got %d", resp.StatusCode)
	}
}
