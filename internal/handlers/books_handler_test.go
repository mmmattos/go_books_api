package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmmattos/books_api/internal/app"
	"github.com/mmmattos/books_api/internal/handlers"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

func setupTestServer() http.Handler {
	repo := memory_book.NewMemoryBookRepo()
	uc := app.NewUsecase(repo)
	return handlers.NewRouter(uc)
}

func TestCreateAndListBooks(t *testing.T) {
	server := setupTestServer()

	payload := map[string]string{
		"title":  "Domain-Driven Design",
		"author": "Eric Evans",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/books", nil)
	rec = httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var books []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&books); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}

	if books[0]["title"] != "Domain-Driven Design" {
		t.Fatalf("unexpected book title: %v", books[0]["title"])
	}
}

func TestGetBookNotFound(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestInvalidMethod(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodPatch, "/books", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}
