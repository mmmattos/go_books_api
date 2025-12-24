package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/mmmattos/books_api/internal/app"
	"github.com/mmmattos/books_api/internal/domain"
	"github.com/mmmattos/books_api/internal/handlers"
	"github.com/mmmattos/books_api/internal/repository/memory_book"
	"github.com/mmmattos/books_api/internal/repository/postgres_book"
)

func main() {
	var repo domain.BookRepository

	dbConn := os.Getenv("DB_CONN")

	if dbConn == "" {
		log.Println("BOOT: using in-memory repository")
		repo = memory_book.NewMemoryBookRepo()
	} else {
		log.Println("BOOT: using postgres repository")

		db, err := sql.Open("postgres", dbConn)
		if err != nil {
			log.Fatalf("db open failed: %v", err)
		}
		defer db.Close()

		repo = postgres_book.NewPostgresBookRepo(db)
	}

	uc := app.NewUsecase(repo)
	router := handlers.NewRouter(uc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("BOOT: listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
