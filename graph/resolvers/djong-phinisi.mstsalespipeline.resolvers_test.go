package graph

import (
	"context"
	"fmt"
	"testing"
	"time"

	"djong-reader-engine/graph/model"

	"github.com/jackc/pgx/v5"
)

// newSalesPipelineScanFunc returns a ScanFunc that populates a SalesPipeline row.
// Only ID, IsActive and createdAt (dest[5]) are set; nullable ptr fields remain nil.
func newSalesPipelineScanFunc(id string, isActive bool, createdAt time.Time) func(dest ...any) error {
	return func(dest ...any) error {
		// dest[0]: *string   → &slsPln.ID
		*(dest[0].(*string)) = id
		// dest[5]: **time.Time → &createdAt (local *time.Time in resolver)
		*(dest[5].(**time.Time)) = &createdAt
		// dest[9]: *bool     → &slsPln.IsActive
		*(dest[9].(*bool)) = isActive
		return nil
	}
}

// ── SalesPipeline (single row) ────────────────────────────────────────────────

func TestSalesPipelineResolver_Success(t *testing.T) {
	now := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	db := &MockDBQuerier{
		QueryRowFunc: func(_ context.Context, _ string, _ ...any) pgx.Row {
			return &MockRow{
				ScanFunc: func(dest ...any) error {
					*(dest[0].(*string)) = "sp-1"
					*(dest[5].(**time.Time)) = &now
					*(dest[9].(*bool)) = true
					return nil
				},
			}
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	sp, err := qr.SalesPipeline(context.Background(), "sp-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sp == nil {
		t.Fatal("expected non-nil SalesPipeline")
	}
	if sp.ID != "sp-1" {
		t.Errorf("expected ID sp-1, got %s", sp.ID)
	}
	if !sp.IsActive {
		t.Error("expected IsActive=true")
	}
}

func TestSalesPipelineResolver_DBError(t *testing.T) {
	db := &MockDBQuerier{
		QueryRowFunc: func(_ context.Context, _ string, _ ...any) pgx.Row {
			return &MockRow{
				ScanFunc: func(dest ...any) error {
					return fmt.Errorf("no rows in result set")
				},
			}
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	_, err := qr.SalesPipeline(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error when QueryRow.Scan fails, got nil")
	}
}

// ── SalesPipelines (list) ─────────────────────────────────────────────────────

func TestSalesPipelinesResolver_DBError(t *testing.T) {
	// errorQueryFunc is defined in mstlead.resolvers_test.go (same package)
	db := &MockDBQuerier{QueryFunc: errorQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	_, err := qr.SalesPipelines(context.Background(), nil, nil, nil)
	if err == nil {
		t.Error("expected error when DB.Query fails, got nil")
	}
}

func TestSalesPipelinesResolver_EmptyResult(t *testing.T) {
	// emptyQueryFunc is defined in mstlead.resolvers_test.go (same package)
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	conn, err := qr.SalesPipelines(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil SalesPipelineConnection")
	}
	if len(conn.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(conn.Edges))
	}
}

func TestSalesPipelinesResolver_SingleRow(t *testing.T) {
	now := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newSalesPipelineScanFunc("sp-1", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	first := int32(5)
	conn, err := qr.SalesPipelines(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "sp-1" {
		t.Errorf("expected ID sp-1, got %s", conn.Edges[0].Node.ID)
	}
}

func TestSalesPipelinesResolver_MultipleRows_HasNextPage(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newSalesPipelineScanFunc("sp-1", true, now),
					newSalesPipelineScanFunc("sp-2", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	first := int32(2) // limit == returned rows
	conn, err := qr.SalesPipelines(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !conn.PageInfo.HasNextPage {
		t.Error("expected HasNextPage=true when len(rows)==limit")
	}
}

func TestSalesPipelinesResolver_HasNextPageFalse(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					newSalesPipelineScanFunc("sp-1", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	first := int32(10) // limit 10, only 1 row returned
	conn, err := qr.SalesPipelines(context.Background(), &first, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.PageInfo.HasNextPage {
		t.Error("expected HasNextPage=false when len(rows) < limit")
	}
}

func TestSalesPipelinesResolver_WithAfterCursor(t *testing.T) {
	after := "cursor-abc"
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	_, err := qr.SalesPipelines(context.Background(), nil, &after, nil)
	if err != nil {
		t.Fatalf("unexpected error with after cursor: %v", err)
	}
}

func TestSalesPipelinesResolver_OrderByASC(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	asc := model.SortOrderSalesPipelineAsc
	orderBy := &model.SalesPipelineOrderByInput{CreatedAt: &asc}
	_, err := qr.SalesPipelines(context.Background(), nil, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error with orderBy ASC: %v", err)
	}
}

func TestSalesPipelinesResolver_OrderByDESC(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DB: db}
	qr := r.Query()

	desc := model.SortOrderSalesPipelineDesc
	orderBy := &model.SalesPipelineOrderByInput{CreatedAt: &desc}
	_, err := qr.SalesPipelines(context.Background(), nil, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error with orderBy DESC: %v", err)
	}
}

func TestSalesPipelinesResolver_ScanError_RowSkipped(t *testing.T) {
	now := time.Now()
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{
				ScanFuncs: []func(dest ...any) error{
					func(dest ...any) error { return fmt.Errorf("scan error on row 1") },
					newSalesPipelineScanFunc("sp-2", true, now),
				},
			}, nil
		},
	}
	r := &Resolver{DB: db}
	qr := r.Query()

	conn, err := qr.SalesPipelines(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("resolver should continue on scan error, got %v", err)
	}
	if len(conn.Edges) != 1 {
		t.Errorf("expected 1 edge (scan-error row skipped), got %d", len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "sp-2" {
		t.Errorf("expected ID sp-2, got %s", conn.Edges[0].Node.ID)
	}
}
