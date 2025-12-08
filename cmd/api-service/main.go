package main

import (
	"log"
	stdhttp "net/http"
	"os"

	"github.com/mmmattos/books_api/internal/app"
	"github.com/mmmattos/books_api/internal/handlers"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
)

func main() {
	repo := memory_book.NewMemoryBookRepo()
	tuc := app.NewUsecase(repo)
	router := handlers.NewRouter(tuc)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("starting server on %s", addr)
	if err := stdhttp.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
