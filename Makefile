# Makefile — upgraded for local dev, Docker, and GCP Cloud Run
# Usage:
#   make build        # build production binary
#   make run          # build then run the binary (prod-style)
#   make run-dev      # run with `go run` (fast dev cycle)
#   make db           # start local Postgres (Docker) and seed it
#   make seed         # seed schema.sql into running DB container
#   make clean-db     # remove local Postgres container
#   make fmt vet test # code hygiene & tests
#   make docker-build # build docker image
#   make docker-run   # run docker image locally
#   make deploy       # build & deploy to Cloud Run (immutable image)
# Environment: override any variable on CLI, e.g. `make run PORT=9090`

APP_NAME ?= books-api
PORT ?= 8080

# ------------------------------
# Local Postgres container
# ------------------------------
DB_USER ?= user
DB_PASS ?= password
DB_NAME ?= booksdb
DB_CONTAINER ?= books-db

# ------------------------------
# GCP / Cloud Run
# ------------------------------
PROJECT ?= gcp-books-api
REGION  ?= us-east1
REPO    ?= book-repo

# Immutable image (git SHA)
TAG        := $(shell git rev-parse --short HEAD)
IMAGE_BASE := $(REGION)-docker.pkg.dev/$(PROJECT)/$(REPO)/$(APP_NAME)
IMAGE      := $(IMAGE_BASE):$(TAG)

# Cloud SQL instance connection name (project:region:instance)
INSTANCE_CONN ?= my-project:us-east1:booksdb

# Helpers
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*")

.PHONY: help build run run-dev fmt vet test \
        db seed clean-db \
        docker-build docker-run docker-push cloud-build \
        deploy deploy-only logs clean info

help:
	@echo "Makefile targets:"
	@echo "  build         Build production binary"
	@echo "  run           Build + run binary (reads PORT env)"
	@echo "  run-dev       Run using 'go run' (fast dev)"
	@echo "  fmt           gofmt -w on repo"
	@echo "  vet           go vet ./..."
	@echo "  test          go test ./..."
	@echo "  db            Start local Postgres Docker container and seed schema"
	@echo "  seed          Load schema.sql into local Postgres container"
	@echo "  clean-db      Stop & remove local Postgres container"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Run Docker image locally"
	@echo "  cloud-build   Build & push image using Cloud Build (immutable)"
	@echo "  deploy        Build + deploy to Cloud Run (immutable)"
	@echo "  deploy-only   Deploy previously built image"
	@echo "  info          Show active Cloud Run revision & image"
	@echo "  logs          Read Cloud Run logs"
	@echo "  clean         Remove local build artifacts"

# ------------------------------
# Build / Run
# ------------------------------
build:
	@echo "=> go mod tidy && go build -o $(APP_NAME) ./cmd/api-service"
	@go mod tidy
	@go build -o $(APP_NAME) ./cmd/api-service

run: build
	@echo "=> Running $(APP_NAME) on :$(PORT)"
	@PORT=$(PORT) ./$(APP_NAME)

run-dev:
	@echo "=> go run ./cmd/api-service (PORT=$(PORT))"
	@PORT=$(PORT) go run ./cmd/api-service

# ------------------------------
# Formatting / Tests
# ------------------------------
fmt:
	@gofmt -w .

vet:
	@echo "=> running go vet"
	@go vet ./...

test:
	@echo "=> running go test ./..."
	@go test ./...

# ------------------------------
# Local Postgres (Docker)
# ------------------------------
db:
	@echo "=> Starting local Postgres container '$(DB_CONTAINER)'"
	@docker run --name $(DB_CONTAINER) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p 5432:5432 -d postgres:16-alpine || true
	@sleep 2
	@$(MAKE) seed

seed:
	@echo "=> Seeding database using schema.sql"
	@docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < schema.sql

clean-db:
	@echo "=> Removing local Postgres container '$(DB_CONTAINER)'"
	@docker rm -f $(DB_CONTAINER) || true

# ------------------------------
# Docker (local only)
# ------------------------------
docker-build:
	@echo "=> docker build -t $(IMAGE) ."
	docker build -t $(IMAGE) .

docker-run:
	@echo "=> docker run -p $(PORT):8080 --rm $(IMAGE)"
	docker run -p $(PORT):8080 --rm $(IMAGE)

docker-push:
	@echo "=> docker push $(IMAGE)"
	docker push $(IMAGE)

# ------------------------------
# Cloud Build (immutable)
# ------------------------------
cloud-build:
	@echo "=> Building image $(IMAGE) with Cloud Build"
	gcloud builds submit \
		--project $(PROJECT) \
		--tag $(IMAGE)

# ------------------------------
# Cloud Run deployment (IMMUTABLE)
# ------------------------------
deploy: cloud-build deploy-only

deploy-only:
	@echo "=> Deploying $(IMAGE) to Cloud Run"
	@if [ -n "$(INSTANCE_CONN)" ]; then \
		CONN="postgres://$(DB_USER):$(DB_PASS)@/$(DB_NAME)?host=/cloudsql/$(INSTANCE_CONN)"; \
		gcloud run deploy $(APP_NAME) \
			--image $(IMAGE) \
			--region $(REGION) \
			--allow-unauthenticated \
			--set-env-vars DB_CONN="$$CONN" \
			--add-cloudsql-instances $(INSTANCE_CONN); \
	else \
		gcloud run deploy $(APP_NAME) \
			--image $(IMAGE) \
			--region $(REGION) \
			--allow-unauthenticated; \
	fi

info:
	@gcloud run services describe $(APP_NAME) \
		--region $(REGION) \
		--format="value(status.latestReadyRevisionName, spec.template.spec.containers[0].image)"

logs:
	@gcloud run logs read $(APP_NAME) --region=$(REGION) --limit=100

# ------------------------------
# Cleanup
# ------------------------------
clean:
	@echo "=> cleaning $(APP_NAME)"
	@rm -f $(APP_NAME)