package graph

import (
	"context"
	"fmt"
	"testing"
	"time"

	"djong-reader-engine/graph/model"

	"github.com/jackc/pgx/v5"
)

// newLeadScanFunc returns a ScanFunc that populates a Lead row for testing.
func newLeadScanFunc(id string, isActive bool, createdAt time.Time) func(dest ...any) error {
	return func(dest ...any) error {
		*(dest[0].(*string)) = id
		*(dest[3].(**time.Time)) = &createdAt
		*(dest[7].(*bool)) = isActive
		return nil
	}
}

func errorQueryFunc(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, fmt.Errorf("db connection refused")
}

func emptyQueryFunc(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return &MockRows{}, nil
}

// ── Lead (not-implemented, panics) ────────────────────────────────────────────

func TestLeadResolver_NotImplemented(t *testing.T) {
	r := &Resolver{DB: &MockDBQuerier{}}
	qr := r.Query()

	defer func() {
		if rec := recover(); rec == nil {
			t.Error("expected panic for unimplemented Lead resolver, got none")
		}
	}()

	qr.Lead(context.Background(), "any-id") //nolint:errcheck
}

// ── Leads ─────────────────────────────────────────────────────────────────────

func TestLeadsResolver_DBError(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: errorQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	_, err := qr.Leads(context.Background(), nil, nil, nil)
	if err == nil {
		t.Error("expected error when DB.Query fails, got nil")
	}
}

func TestLeadsResolver_EmptyResult(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	conn, err := qr.Leads(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil LeadConnection")
	}
	if len(conn.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(conn.Edges))
	}
}

func TestLeadsResolver_SingleRow(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{newLeadScanFunc("lead-1", true, now)},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	first := int32(10)
	conn, err := qr.Leads(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "lead-1" {
		t.Errorf("expected ID lead-1, got %s", conn.Edges[0].Node.ID)
	}
}

func TestLeadsResolver_MultipleRows_HasNextPage(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newLeadScanFunc("lead-1", true, now),
					newLeadScanFunc("lead-2", false, now),
					newLeadScanFunc("lead-3", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	first := int32(3)
	conn, err := qr.Leads(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(conn.Edges))
	}
	if !conn.PageInfo.HasNextPage {
		t.Error("expected HasNextPage=true when len(rows)==limit")
	}
}

func TestLeadsResolver_HasNextPageFalse(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{newLeadScanFunc("lead-1", true, now)},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	first := int32(10)
	conn, err := qr.Leads(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.PageInfo.HasNextPage {
		t.Error("expected HasNextPage=false when len(rows) < limit")
	}
}

func TestLeadsResolver_WithAfterCursor(t *testing.T) {
	after := int32(5)
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	_, err := qr.Leads(context.Background(), nil, &after, nil)
	if err != nil {
		t.Fatalf("unexpected error with after cursor: %v", err)
	}
}

func TestLeadsResolver_OrderByASC(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	asc := model.SortOrderLeadAsc
	orderBy := &model.LeadOrderByInput{CreatedAt: &asc}
	_, err := qr.Leads(context.Background(), nil, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error with orderBy ASC: %v", err)
	}
}

func TestLeadsResolver_OrderByDESC(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	desc := model.SortOrderLeadDesc
	orderBy := &model.LeadOrderByInput{CreatedAt: &desc}
	_, err := qr.Leads(context.Background(), nil, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error with orderBy DESC: %v", err)
	}
}

func TestLeadsResolver_ScanError_RowSkipped(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					func(dest ...any) error { return fmt.Errorf("scan error on row 1") },
					newLeadScanFunc("lead-2", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	conn, err := qr.Leads(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("resolver should continue on scan error, got %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Errorf("expected 1 edge (scan-error row skipped), got %d", len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "lead-2" {
		t.Errorf("expected ID lead-2, got %s", conn.Edges[0].Node.ID)
	}
}

// TestLeadsResolver_OrderByNilCreatedAt covers the branch where orderBy is
// non-nil but orderBy.CreatedAt is nil — the inner sort-field assignment is
// skipped and the default "created_at DESC" is used.
func TestLeadsResolver_OrderByNilCreatedAt(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	// orderBy present but CreatedAt is nil → inner branch NOT entered
	orderBy := &model.LeadOrderByInput{CreatedAt: nil}
	_, err := qr.Leads(context.Background(), nil, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error with nil CreatedAt in orderBy: %v", err)
	}
}

// TestLeadsResolver_NilCreatedAt covers the branch where the DB returns a nil
// createdAt, so model.TimeFromPtr returns nil and the time assignment is skipped.
func TestLeadsResolver_NilCreatedAt(t *testing.T) {
	nilCreatedAtScanFunc := func(dest ...any) error {
		// dest[0]: *string  → ID
		*(dest[0].(*string)) = "lead-nil-ts"
		// dest[3]: **time.Time → leave as nil (don't write)
		*(dest[3].(**time.Time)) = nil
		// dest[7]: *bool → IsActive
		*(dest[7].(*bool)) = true
		return nil
	}

	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{nilCreatedAtScanFunc},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	conn, err := qr.Leads(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error with nil createdAt: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
	// CreatedAt should remain the zero value of time.Time — no panic
	if !conn.Edges[0].Node.CreatedAt.IsZero() {
		t.Errorf("expected zero CreatedAt when DB returns nil, got %v", conn.Edges[0].Node.CreatedAt)
	}
}

// TestLeadsResolver_EdgeCursorAndEndCursor verifies that each edge's Cursor
// equals the node ID, and PageInfo.EndCursor reflects the last row's ID.
func TestLeadsResolver_EdgeCursorAndEndCursor(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newLeadScanFunc("lead-A", true, now),
					newLeadScanFunc("lead-B", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	first := int32(10)
	conn, err := qr.Leads(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Each edge's Cursor must equal its node's ID
	for i, edge := range conn.Edges {
		if edge.Cursor != edge.Node.ID {
			t.Errorf("edge[%d]: Cursor %q != Node.ID %q", i, edge.Cursor, edge.Node.ID)
		}
	}

	// EndCursor must be the last row's ID
	if conn.PageInfo.EndCursor == nil {
		t.Fatal("expected non-nil EndCursor")
	}
	if *conn.PageInfo.EndCursor != "lead-B" {
		t.Errorf("expected EndCursor=lead-B, got %s", *conn.PageInfo.EndCursor)
	}
}
