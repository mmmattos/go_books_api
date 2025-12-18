#!/usr/bin/env bash
# Deploy to Cloud Run in project gcp-books-api. Configurable via env or args.

# Cloud Run deploy for project gcp-books-api)

# Preconditions: you said the GCP project gcp-books-api exists. This script assumes:
# 	•	gcloud is installed and authenticated,
# 	•	Cloud SQL instance may or may not exist (if it does, pass its connection name as INSTANCE_CONN),
# 	•	schema.sql exists for seeding (script seeds if instance exists and gcloud sql connect available).

set -euo pipefail
. "$(dirname "$0")/_helpers.sh" || true

PROJECT="${PROJECT:-gcp-books-api}"
REGION="${REGION:-us-central1}"
APP_NAME="${APP_NAME:-books-api}"
IMAGE="gcr.io/${PROJECT}/${APP_NAME}"
INSTANCE_CONN="${INSTANCE_CONN:-}"  # expected like project:region:instance
DB_USER="${DB_USER:-postgres}"
DB_PASS="${DB_PASS:-password}"
DB_NAME="${DB_NAME:-booksdb}"

info "Using project: ${PROJECT}, region: ${REGION}, service: ${APP_NAME}"

info "Enabling required APIs..."
gcloud services enable run.googleapis.com sqladmin.googleapis.com cloudbuild.googleapis.com --project "${PROJECT}"

info "Building and pushing image to ${IMAGE} using Cloud Build..."
gcloud builds submit --tag "${IMAGE}" --project "${PROJECT}"

if [ -n "${INSTANCE_CONN}" ]; then
  info "Deploying to Cloud Run and attaching Cloud SQL instance ${INSTANCE_CONN}..."
  # DB_CONN uses unix socket path for Cloud Run
  CONN="postgres://${DB_USER}:${DB_PASS}@/${DB_NAME}?host=/cloudsql/${INSTANCE_CONN}"
  gcloud run deploy "${APP_NAME}" \
    --image "${IMAGE}" \
    --region "${REGION}" \
    --platform managed \
    --project "${PROJECT}" \
    --allow-unauthenticated \
    --set-env-vars "PORT=8080,DB_CONN=${CONN}" \
    --add-cloudsql-instances "${INSTANCE_CONN}"
else
  info "Deploying to Cloud Run without Cloud SQL instance (no DB_CONN injected)..."
  gcloud run deploy "${APP_NAME}" \
    --image "${IMAGE}" \
    --region "${REGION}" \
    --platform managed \
    --project "${PROJECT}" \
    --allow-unauthenticated \
    --set-env-vars "PORT=8080"
fi

info "Deployment finished. Use 'gcloud run services describe ${APP_NAME} --region ${REGION} --project ${PROJECT}' to see the URL."