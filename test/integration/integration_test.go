package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/bookapi/internal/app"
	"github.com/user/bookapi/internal/domain"
	"github.com/user/bookapi/repository/memory_book"
	"github.com/user/bookapi/internal/http"
)

func TestBooksAPI_Integration(t *testing.T) {
	repo := memory_book.NewMemoryBookRepo()
	uc := app.NewUsecase(repo)
	r := http.NewRouter(tuc)
	ts := httptest.NewServer(r)
	defer ts.Close()

	b := domain.Book{Title: "Test", Author: "Me", Year: 2025}
	buf, _ := json.Marshal(b)
	resp, err := http.Post(ts.URL+"/books", "application/json", bytes.NewBuffer(buf))
	if err != nil { t.Fatalf("post failed: %v", err) }
	if resp.StatusCode != http.StatusCreated { t.Fatalf("expected 201 got %d", resp.StatusCode) }
}
