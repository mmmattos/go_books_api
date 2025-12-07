package http

import (
	"encoding/json"
	"net/http"

	"github.com/user/bookapi/internal/app"
)

// Router wires handlers.
type Router struct {
	mux     *http.ServeMux
	usecase *app.Usecase
}

func NewRouter(u *app.Usecase) *Router {
	r := &Router{mux: http.NewServeMux(), usecase: u}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
	bh := NewBookHandler(r.usecase)

	r.mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.mux.HandleFunc("/books", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			bh.ListBooks(w, req)
		case http.MethodPost:
			bh.CreateBook(w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	r.mux.HandleFunc("/books/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			bh.GetBook(w, req)
		case http.MethodPut:
			bh.UpdateBook(w, req)
		case http.MethodDelete:
			bh.DeleteBook(w, req)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// helpers

func respondJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
}

func decodeJSON(req *http.Request, v interface{}) error {
	dec := json.NewDecoder(req.Body)
	return dec.Decode(v)
}
