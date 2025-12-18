package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mmmattos/books_api/internal/app"
	"github.com/mmmattos/books_api/internal/handlers"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

// These are set at build time via -ldflags.
// Defaults are safe for local dev.
var (
	Version   = "dev"
	CommitSHA = "unknown"
)

func main() {
	log.Println("BOOT MARKER: starting books-api")

	repo := memory_book.NewMemoryBookRepo()
	uc := app.NewUsecase(repo)
	api := handlers.NewRouter(uc)

	// Root mux (infra-level)
	root := http.NewServeMux()

	// Health endpoint (Cloud Run–safe)
	root.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST HIT: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Debug/version endpoint
	root.HandleFunc("/debug/version", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST HIT: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(
			`{"service":"books-api","version":"` + Version + `","commit":"` + CommitSHA + `"}`,
		))
	})

	// API routes
	root.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST HIT: %s %s", r.Method, r.URL.Path)
		api.ServeHTTP(w, r)
	}))

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("BOOT MARKER: listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, root))
}
