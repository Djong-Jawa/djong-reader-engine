package graph

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ── MockDBQuerier ─────────────────────────────────────────────────────────────

// MockDBQuerier is a test double for DBQuerier.
// Supply QueryFunc / QueryRowFunc to control what the mock returns.
type MockDBQuerier struct {
	QueryFunc    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
}

func (m *MockDBQuerier) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}
	return &MockRows{}, nil
}

func (m *MockDBQuerier) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}
	return &MockRow{ScanFunc: func(dest ...any) error { return fmt.Errorf("no mock row configured") }}
}

// ── MockRows ──────────────────────────────────────────────────────────────────

// MockRows implements pgx.Rows.
// ScanFuncs is an ordered list of closures, one per row to be returned.
type MockRows struct {
	ScanFuncs []func(dest ...any) error
	CloseErr  error
	pos       int
}

func (r *MockRows) Close() {}

func (r *MockRows) Err() error { return r.CloseErr }

func (r *MockRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }

func (r *MockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }

func (r *MockRows) Next() bool {
	r.pos++
	return r.pos <= len(r.ScanFuncs)
}

func (r *MockRows) Scan(dest ...any) error {
	if r.pos-1 < len(r.ScanFuncs) {
		return r.ScanFuncs[r.pos-1](dest...)
	}
	return fmt.Errorf("no row at position %d", r.pos-1)
}

func (r *MockRows) Values() ([]any, error) { return nil, nil }

func (r *MockRows) RawValues() [][]byte { return nil }

func (r *MockRows) Conn() *pgx.Conn { return nil }

// ── MockRow ───────────────────────────────────────────────────────────────────

// MockRow implements pgx.Row (single-row query result).
type MockRow struct {
	ScanFunc func(dest ...any) error
}

func (r *MockRow) Scan(dest ...any) error {
	if r.ScanFunc != nil {
		return r.ScanFunc(dest...)
	}
	return fmt.Errorf("no scan function configured")
}
