package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/user/bookapi/internal/app"
	"github.com/user/bookapi/internal/domain"
)

// BookHandler provides HTTP handlers.
type BookHandler struct {
	uc *app.Usecase
}

// NewBookHandler constructs a handler.
func NewBookHandler(uc *app.Usecase) *BookHandler {
	return &BookHandler{uc: uc}
}

func extractID(path string) (string, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[1] == "" {
		return "", fmt.Errorf("invalid id")
	}
	return parts[1], nil
}

func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.uc.GetAllBooks()
	if err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, books)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b domain.Book
	if err := decodeJSON(r, &b); err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.uc.CreateBook(&b); err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	respondJSON(w, http.StatusCreated, b)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	b, err := h.uc.GetBookByID(id)
	if err != nil {
		httpError(w, http.StatusNotFound, err)
		return
	}
	respondJSON(w, http.StatusOK, b)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	var b domain.Book
	if err := decodeJSON(r, &b); err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	b.ID = id
	if err := h.uc.UpdateBook(&b); err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	respondJSON(w, http.StatusOK, b)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		httpError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.uc.DeleteBook(id); err != nil {
		httpError(w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
