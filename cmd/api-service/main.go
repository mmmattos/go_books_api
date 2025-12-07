package main

import (
	"log"
	"net/http"
	"os"

	"github.com/user/bookapi/internal/app"
	httpint "github.com/user/bookapi/internal/http"
	"github.com/user/bookapi/repository/memory_book"
)

func main() {
	repo := memory_book.NewMemoryBookRepo()
	tuc := app.NewUsecase(repo)
	router := httpint.NewRouter(tuc)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
