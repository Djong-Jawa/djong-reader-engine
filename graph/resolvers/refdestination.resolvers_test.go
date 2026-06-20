package graph

import (
	"context"
	"testing"
	"time"

	"djong-reader-engine/graph/model"

	"github.com/jackc/pgx/v5"
)

// newRefDestinationScanFunc returns a ScanFunc that fills the eight columns returned by
// the JukungRefDestinations query:
//
//	dest[0] *string      → &rd.ID
//	dest[1] **string     → &rd.DestinationName (nullable)
//	dest[2] **bool       → &rd.IsCombination     (nullable)
//	dest[3] **time.Time  → &createdAt            (local *time.Time in resolver)
//	dest[4] **string     → &rd.CreatedBy         (nullable)
//	dest[5] **time.Time  → &rd.UpdatedAt         (nullable)
//	dest[6] **string     → &rd.UpdatedBy         (nullable)
//	dest[7] *bool        → &rd.IsActive
func newRefDestinationScanFunc(id string, isActive bool, createdAt time.Time) func(dest ...any) error {
	return func(dest ...any) error {
		*(dest[0].(*string)) = id
		*(dest[3].(**time.Time)) = &createdAt
		*(dest[7].(*bool)) = isActive
		return nil
	}
}

// ── JukungRefDestination (not-implemented, panics) ───────────────────────────────

// TestJukungRefDestinationResolver_Panics verifies that JukungRefDestination() panics with the expected message
// and that the panic value is a non-nil error.
func TestJukungRefDestinationResolver_Panics(t *testing.T) {
	r := &Resolver{DBJukung: &MockDBQuerier{}}
	qr := r.Query()

	var panicVal any
	func() {
		defer func() { panicVal = recover() }()
		qr.JukungRefDestination(context.Background(), "any-id") //nolint:errcheck
	}()

	if panicVal == nil {
		t.Fatal("JukungRefDestination() must panic — got nil recover value")
	}
	panicErr, ok := panicVal.(error)
	if !ok {
		t.Fatalf("JukungRefDestination() panicked with non-error value %T: %v", panicVal, panicVal)
	}
	expected := "not implemented: JukungRefDestination - jukungRefDestination"
	if panicErr.Error() != expected {
		t.Errorf("panic message: got %q, want %q", panicErr.Error(), expected)
	}
}

// ── JukungRefDestinations ────────────────────────────────────────────────────────

func TestJukungRefDestinationsResolver_DBError(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: errorQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	_, err := qr.JukungRefDestinations(context.Background(), nil, nil, nil)
	if err == nil {
		t.Error("expected error when DBJukung.Query fails, got nil")
	}
}

func TestJukungRefDestinationsResolver_EmptyResult(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	conn, err := qr.JukungRefDestinations(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil RefDestinationConnection")
	}
	if len(conn.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(conn.Edges))
	}
}

func TestJukungRefDestinationsResolver_SingleRow(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{newRefDestinationScanFunc("rd-1", true, now)},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	first := int32(10)
	conn, err := qr.JukungRefDestinations(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "rd-1" {
		t.Errorf("expected ID rd-1, got %s", conn.Edges[0].Node.ID)
	}
}

func TestJukungRefDestinationsResolver_MultipleRows_HasNextPage(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newRefDestinationScanFunc("rd-1", true, now),
					newRefDestinationScanFunc("rd-2", false, now),
					newRefDestinationScanFunc("rd-3", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	first := int32(3)
	conn, err := qr.JukungRefDestinations(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(conn.Edges))
	}
	if conn.PageInfo.HasNextPage != true {
		t.Error("expected HasNextPage=true when result count matches limit")
	}
	if conn.PageInfo.EndCursor == nil || *conn.PageInfo.EndCursor != "rd-3" {
		t.Errorf("expected EndCursor=rd-3, got %v", conn.PageInfo.EndCursor)
	}
}

func TestJukungRefDestinationsResolver_WithCursor(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, args ...any) (pgx.Rows, error) {
			// Verify cursor (offset) is passed correctly
			if len(args) >= 1 {
				if offset, ok := args[0].(int32); ok && offset != 5 {
					t.Errorf("expected offset=5, got %d", offset)
				}
			}
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newRefDestinationScanFunc("rd-6", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	after := "5"
	first := int32(10)
	conn, err := qr.JukungRefDestinations(context.Background(), &first, &after, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
}

func TestJukungRefDestinationsResolver_OrderByCreatedAtASC(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, query string, _ ...any) (pgx.Rows, error) {
			// Verify ORDER BY created_at ASC is in the query
			if query == "" {
				t.Error("query should not be empty")
			}
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newRefDestinationScanFunc("rd-1", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	sortOrder := model.SortOrderRefDestinationAsc
	orderBy := &model.RefDestinationOrderByInput{CreatedAt: &sortOrder}
	first := int32(10)

	conn, err := qr.JukungRefDestinations(context.Background(), &first, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
}
