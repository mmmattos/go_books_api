FROM golang:1.23-alpine AS builder
WORKDIR /src

# Accept build arguments (provided by Cloud Build / Makefile)
ARG VERSION=dev
ARG COMMIT_SHA=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with version injected
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -ldflags "\
        -X main.Version=${VERSION} \
        -X main.CommitSHA=${COMMIT_SHA}" \
      -o /out/books-api \
      ./cmd/api-service

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/books-api /books-api
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/books-api"]