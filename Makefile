# Makefile — local dev + Docker + Cloud Run (immutable, versioned)

APP_NAME ?= books-api

PROJECT ?= gcp-books-api
REGION  ?= us-east1
REPO    ?= book-repo

TAG        := $(shell git rev-parse --short HEAD)
IMAGE_BASE := $(REGION)-docker.pkg.dev/$(PROJECT)/$(REPO)/$(APP_NAME)
IMAGE      := $(IMAGE_BASE):$(TAG)

INSTANCE_CONN ?=
DB_USER ?= user
DB_PASS ?= password
DB_NAME ?= booksdb

.PHONY: deploy deploy-only cloud-build info logs

cloud-build:
	@echo "==> Cloud Build: $(IMAGE)"
	gcloud builds submit \
		--project $(PROJECT) \
		--config cloudbuild.yaml \
		--substitutions=_IMAGE=$(IMAGE),_VERSION=$(TAG),_COMMIT_SHA=$(TAG)

deploy: cloud-build deploy-only

deploy-only:
	@echo "==> Cloud Run deploy: $(IMAGE)"
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
	@gcloud run services logs read $(APP_NAME) \
		--region $(REGION) \
		--limit=100