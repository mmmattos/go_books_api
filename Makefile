# Makefile — upgraded for local dev, Docker, and GCP Cloud Run
# Usage:
#   make deploy       # build & deploy to Cloud Run (immutable, versioned)

APP_NAME ?= books-api
PORT ?= 8080

# ------------------------------
# GCP / Cloud Run
# ------------------------------
PROJECT ?= gcp-books-api
REGION  ?= us-east1
REPO    ?= book-repo

# Immutable version (git SHA)
TAG        := $(shell git rev-parse --short HEAD)
IMAGE_BASE := $(REGION)-docker.pkg.dev/$(PROJECT)/$(REPO)/$(APP_NAME)
IMAGE      := $(IMAGE_BASE):$(TAG)

# Cloud SQL instance connection name (optional)
INSTANCE_CONN ?= my-project:us-east1:booksdb
DB_USER ?= user
DB_PASS ?= password
DB_NAME ?= booksdb

.PHONY: deploy deploy-only cloud-build info logs

# ------------------------------
# Cloud Build (with version injection)
# ------------------------------
cloud-build:
	@echo "==> Building image $(IMAGE)"
	gcloud builds submit \
		--project $(PROJECT) \
		--tag $(IMAGE) \
		--build-arg VERSION=$(TAG) \
		--build-arg COMMIT_SHA=$(TAG)

# ------------------------------
# Cloud Run deployment
# ------------------------------
deploy: cloud-build deploy-only

deploy-only:
	@echo "==> Deploying $(IMAGE)"
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
	gcloud run services describe $(APP_NAME) \
		--region $(REGION) \
		--format="value(status.latestReadyRevisionName, spec.template.spec.containers[0].image)"

logs:
	gcloud run logs read $(APP_NAME) --region=$(REGION) --limit=50