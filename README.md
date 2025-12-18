# Book API – Go REST Service

A clean, extensible REST API for managing books.  
Implements a layered, production-style Go architecture:

- **cmd/api-service** – application entrypoint  
- **internal/domain** – entities + repository interfaces  
- **internal/app** – business logic (usecases)  
- **internal/http** – handlers + routing  
- **repository/** – data stores (in-memory + PostgreSQL)  
- **schema.sql** – DB schema + seed data  
- **Makefile** – complete dev/build/test/deploy workflow  
- **Dockerfile** – multi-stage minimal container build  
- **Cloud Run deploy script** – optional GCP deployment  

---

## Features

### REST API Endpoints

| Method | Path          | Description |
|--------|---------------|-------------|
| GET    | `/healthz`    | Health check |
| GET    | `/books`      | List all books |
| POST   | `/books`      | Create a new book |
| GET    | `/books/{id}` | Retrieve a book |
| PUT    | `/books/{id}` | Update a book |
| DELETE | `/books/{id}` | Delete a book |

### In-memory repository (default)

Perfect for local development and testing.

### PostgreSQL repository (stub included)

Ready to extend with SQL queries.

### Docker support

Multi-stage build → extremely small distroless runtime image.

### Cloud Run–ready

Deployment script, Cloud SQL socket connection, configurable runtime variables.

---

## Project Structure

```
books_api/
├── cmd/api-service/main.go
├── internal/
│   ├── domain/
│   ├── app/
│   ├── http/
│   └── metrics/
├── repository/
│   ├── memory_book/
│   └── postgres_book/
├── test/integration/
├── schema.sql
├── Makefile
├── Dockerfile
└── go.mod
```

---

## Requirements

- Go **1.20+**
- Docker (optional but recommended)
- GCloud CLI (for Cloud Run deployment)
- PostgreSQL (optional: locally via Docker)

---

## Running Locally (Default: In-memory Repository)

### Fastest option:

```bash
make run
```

This:

1. Runs `make build`
2. Executes the compiled binary `./books-api`

Service will run at:

```
http://localhost:8080
```

Test:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/books
```

Create a book:

```bash
curl -X POST http://localhost:8080/books   -H "Content-Type: application/json"   -d '{"title":"Dune","author":"Frank Herbert","year":1965}'
```

### Development mode (no binary):

```bash
make run-dev
```

---

## 🧰 Makefile Commands

| Command | Description |
|--------|-------------|
| `make build` | Build binary into `./books-api` |
| `make run` | Build + run (production-like) |
| `make run-dev` | Run using `go run` (development) |
| `make fmt` | Format with `gofmt -w` |
| `make vet` | Static analysis (`go vet`) |
| `make test` | Run all tests |
| `make db` | Start local Postgres (Docker) + seed |
| `make seed` | Apply schema.sql into Postgres |
| `make clean-db` | Remove Postgres Docker container |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run container locally |
| `make docker-push` | Push container to registry |
| `make cloud-build` | Cloud Build: build + push |
| `make deploy` | Deploy to Cloud Run |
| `make logs` | Tail Cloud Run logs |
| `make clean` | Remove built binaries |

---

## 🧩 Makefile Variables

| Variable | Default | Purpose |
|---------|---------|---------|
| `APP_NAME` | `books-api` | Binary name & Cloud Run service |
| `PORT` | `8080` | Service port |
| `DB_USER` | `user` | Local Postgres user |
| `DB_PASS` | `password` | Local Postgres password |
| `DB_NAME` | `booksdb` | Local Postgres DB name |
| `DB_CONTAINER` | `books-db` | Docker container name |
| `PROJECT` | from gcloud | GCP project |
| `REGION` | `us-central1` | Cloud Run region |
| `IMAGE` | `gcr.io/$(PROJECT)/$(APP_NAME)` | Container image |
| `INSTANCE_CONN` | `my-project:us-central1:booksdb` | Cloud SQL instance |

---

## Examples of Overriding Variables

Run on a different port:

```bash
make run PORT=9090
```

Specify GCP project:

```bash
make deploy PROJECT=my-gcp-project
```

Use a different Cloud SQL instance:

```bash
make deploy INSTANCE_CONN=myproj:us-central1:mydb
```

Run Postgres locally with custom credentials:

```bash
make db DB_USER=app DB_PASS=secret DB_NAME=bookstore
```

---

## Running PostgreSQL Locally (Docker)

Start DB + seed schema:

```bash
make db
```

Apply schema again:

```bash
make seed
```

Remove DB:

```bash
make clean-db
```

---

## 🔌 Connecting the API to PostgreSQL

Edit `cmd/api-service/main.go` and replace:

```go
repo := memory_book.NewMemoryBookRepo()
```

with:

```go
db, _ := sql.Open("postgres", os.Getenv("DB_CONN"))
repo := postgres_book.NewPostgresBookRepo(db)
```

Run API with DB connection:

```bash
DB_CONN="postgres://user:password@localhost:5432/booksdb?sslmode=disable" make run
```

---

## Docker

Build:

```bash
make docker-build
```

Run:

```bash
make docker-run
```

Push:

```bash
make docker-push
```

---

## Deploy to Google Cloud Run

Build via Cloud Build:

```bash
make cloud-build
```

Deploy:

```bash
make deploy
```

Check logs:

```bash
make logs
```

---

## 🧪 Tests

```bash
make test
```

---

## 🩺 Static Analysis

```bash
make vet
```

---

## Formatting

```bash
make fmt
```


---

## License

MIT (or your preferred license)

---

## Contributing

PRs welcome.
