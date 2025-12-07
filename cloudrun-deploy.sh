#!/bin/bash
set -e
PROJECT=$(gcloud config get-value project)
gcloud builds submit --tag gcr.io/${PROJECT}/books-api
gcloud run deploy books-api --image gcr.io/${PROJECT}/books-api --region=us-central1 --platform managed --allow-unauthenticated
