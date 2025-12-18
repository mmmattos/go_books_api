FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ MUST match where main.go lives
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/books-api ./cmd/api-service

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/books-api /books-api
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/books-api"]