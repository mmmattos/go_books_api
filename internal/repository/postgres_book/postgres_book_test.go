package postgres_book_test

import "testing"

// NOTE:
// The PostgresBookRepo is currently a stub implementation.
// All methods return nil and do not hit a real database.
//
// We explicitly skip these tests to:
// - keep `go test ./...` green
// - avoid false confidence
// - document that DB tests must be added once implemented
func TestPostgresRepository_NotImplemented(t *testing.T) {
	t.Skip("postgres repository not implemented yet; skipping DB tests")
}
