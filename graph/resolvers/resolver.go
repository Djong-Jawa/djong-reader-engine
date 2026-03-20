package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// DBQuerier abstracts the database query methods used by resolvers.
// *pgxpool.Pool satisfies this interface, enabling easy mocking in tests.
type DBQuerier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Resolver is the root GraphQL resolver.
type Resolver struct {
	DB DBQuerier
}
