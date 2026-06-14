package graph

import (
	"context"
	"fmt"
	"testing"
	"time"

	"djong-reader-engine/graph/model"

	"github.com/jackc/pgx/v5"
)

// newPricingTierScanFunc returns a ScanFunc that fills the eleven columns returned by
// the JukungPricingTiers query:
//
//	dest[0] *string      → &pt.ID
//	dest[1] **int32      → &pt.ProductID   (nullable)
//	dest[2] **int32      → &pt.PrValidMin  (nullable)
//	dest[3] **int32      → &pt.PrValidMax  (nullable)
//	dest[4] **time.Time  → &createdAt      (local *time.Time in resolver)
//	dest[5] **string     → &pt.CreatedBy   (nullable)
//	dest[6] **time.Time  → &pt.UpdatedAt   (nullable)
//	dest[7] **string     → &pt.UpdatedBy   (nullable)
//	dest[8] *bool        → &pt.IsActive
//	dest[9] **float64    → &pt.Rate        (nullable)
//	dest[10] **int32     → &pt.CurrencyID  (nullable)
func newPricingTierScanFunc(id string, isActive bool, createdAt time.Time) func(dest ...any) error {
	return func(dest ...any) error {
		*(dest[0].(*string)) = id
		*(dest[4].(**time.Time)) = &createdAt
		*(dest[8].(*bool)) = isActive
		return nil
	}
}

func errorQueryFunc(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, fmt.Errorf("db connection refused")
}

func emptyQueryFunc(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return &MockRows{}, nil
}

// ── JukungPricingTier (not-implemented, panics) ───────────────────────────────

// TestJukungPricingTierResolver_Panics verifies that JukungPricingTier() panics with the expected message
// and that the panic value is a non-nil error.
func TestJukungPricingTierResolver_Panics(t *testing.T) {
	r := &Resolver{DBJukung: &MockDBQuerier{}}
	qr := r.Query()

	var panicVal any
	func() {
		defer func() { panicVal = recover() }()
		qr.JukungPricingTier(context.Background(), "any-id") //nolint:errcheck
	}()

	if panicVal == nil {
		t.Fatal("JukungPricingTier() must panic — got nil recover value")
	}
	panicErr, ok := panicVal.(error)
	if !ok {
		t.Fatalf("JukungPricingTier() panicked with non-error value %T: %v", panicVal, panicVal)
	}
	expected := "not implemented: JukungPricingTier - jukungPricingTier"
	if panicErr.Error() != expected {
		t.Errorf("panic message: got %q, want %q", panicErr.Error(), expected)
	}
}

// ── JukungPricingTiers ────────────────────────────────────────────────────────

func TestJukungPricingTiersResolver_DBError(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: errorQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	_, err := qr.JukungPricingTiers(context.Background(), nil, nil, nil)
	if err == nil {
		t.Error("expected error when DBJukung.Query fails, got nil")
	}
}

func TestJukungPricingTiersResolver_EmptyResult(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	conn, err := qr.JukungPricingTiers(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil PricingTierConnection")
	}
	if len(conn.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(conn.Edges))
	}
}

func TestJukungPricingTiersResolver_SingleRow(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{newPricingTierScanFunc("pt-1", true, now)},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	first := int32(10)
	conn, err := qr.JukungPricingTiers(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "pt-1" {
		t.Errorf("expected ID pt-1, got %s", conn.Edges[0].Node.ID)
	}
}

func TestJukungPricingTiersResolver_MultipleRows_HasNextPage(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newPricingTierScanFunc("pt-1", true, now),
					newPricingTierScanFunc("pt-2", false, now),
					newPricingTierScanFunc("pt-3", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	first := int32(3)
	conn, err := qr.JukungPricingTiers(context.Background(), &first, nil, nil)
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

func TestJukungPricingTiersResolver_HasNextPageFalse(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{newPricingTierScanFunc("pt-1", true, now)},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	first := int32(10)
	conn, err := qr.JukungPricingTiers(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.PageInfo.HasNextPage {
		t.Error("expected HasNextPage=false when len(rows) < limit")
	}
}

func TestJukungPricingTiersResolver_WithAfterCursor(t *testing.T) {
	after := "5"
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	_, err := qr.JukungPricingTiers(context.Background(), nil, &after, nil)
	if err != nil {
		t.Fatalf("unexpected error with after cursor: %v", err)
	}
}

func TestJukungPricingTiersResolver_WithOrderBy(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{newPricingTierScanFunc("pt-1", true, now)},
			}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	asc := model.SortOrderPricingTierAsc
	orderBy := &model.PricingTierOrderByInput{CreatedAt: &asc}
	first := int32(10)

	conn, err := qr.JukungPricingTiers(context.Background(), &first, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error with orderBy: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
}
