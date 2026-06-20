package graph

import (
	"context"
	"fmt"
	"testing"
	"time"

	"djong-reader-engine/graph/model"

	"github.com/jackc/pgx/v5"
)

// newProductScanFunc returns a ScanFunc that fills the nine columns returned by
// the MstProduct query:
//
//	dest[0] *string      → &p.ID
//	dest[1] **int32      → &p.DestinationID (nullable, left nil)
//	dest[2] **string     → &p.Description (nullable, left nil)
//	dest[3] **int32      → &p.TotalDuration (nullable, left nil)
//	dest[4] **string     → &p.CreatedBy (nullable, left nil)
//	dest[5] **time.Time  → &createdAt (local *time.Time in resolver)
//	dest[6] **string     → &p.UpdatedBy (nullable, left nil)
//	dest[7] **time.Time  → &p.UpdatedAt (nullable, left nil)
//	dest[8] *bool        → &p.IsActive
func newProductScanFunc(id string, isActive bool, createdAt time.Time) func(dest ...any) error {
	return func(dest ...any) error {
		*(dest[0].(*string)) = id
		*(dest[5].(**time.Time)) = &createdAt
		*(dest[8].(*bool)) = isActive
		return nil
	}
}

// ── MstProduct ────────────────────────────────────────────────────────────────

func TestMstProductResolver_DBError(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: errorQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	_, err := qr.MstProduct(context.Background(), nil, nil, nil)
	if err == nil {
		t.Error("expected error when DBJukung.Query fails, got nil")
	}
}

func TestMstProductResolver_EmptyResult(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	conn, err := qr.MstProduct(context.Background(), nil, nil, nil)
	if err != nil {
		t.Errorf("expected no error with empty result set; got %v", err)
		return
	}

	if conn == nil {
		t.Error("connection must not be nil")
		return
	}

	if len(conn.Edges) != 0 {
		t.Errorf("expected 0 edges; got %d", len(conn.Edges))
	}

	if conn.PageInfo == nil {
		t.Error("pageInfo must not be nil")
		return
	}

	if conn.PageInfo.HasNextPage {
		t.Error("hasNextPage should be false for empty result")
	}
}

func TestMstProductResolver_SingleRow(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	rows := &MockRows{
		data: []func(dest ...any) error{
			newProductScanFunc("prod-1", true, now),
		},
	}

	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return rows, nil
		},
	}

	r := &Resolver{DBJukung: db}
	qr := r.Query()

	conn, err := qr.MstProduct(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conn == nil || len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge; got %v", conn)
	}

	edge := conn.Edges[0]
	if edge.Cursor != "prod-1" {
		t.Errorf("cursor: got %q, want %q", edge.Cursor, "prod-1")
	}

	if edge.Node.ID != "prod-1" {
		t.Errorf("node.ID: got %q, want %q", edge.Node.ID, "prod-1")
	}

	if !edge.Node.IsActive {
		t.Error("node.IsActive should be true")
	}

	expectedTime := now
	if edge.Node.CreatedAt != expectedTime {
		t.Errorf("node.CreatedAt: got %v, want %v", edge.Node.CreatedAt, expectedTime)
	}

	if conn.PageInfo == nil {
		t.Fatal("pageInfo must not be nil")
	}

	if conn.PageInfo.EndCursor == nil || *conn.PageInfo.EndCursor != "prod-1" {
		t.Errorf("endCursor: got %v, want %q", conn.PageInfo.EndCursor, "prod-1")
	}

	if conn.PageInfo.HasNextPage {
		t.Error("hasNextPage should be false when result size < limit")
	}
}

func TestMstProductResolver_MultipleRows(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	rows := &MockRows{
		data: []func(dest ...any) error{
			newProductScanFunc("prod-1", true, now),
			newProductScanFunc("prod-2", true, now.Add(1*time.Hour)),
			newProductScanFunc("prod-3", false, now.Add(2*time.Hour)),
		},
	}

	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return rows, nil
		},
	}

	r := &Resolver{DBJukung: db}
	qr := r.Query()

	limit := int32(3)
	conn, err := qr.MstProduct(context.Background(), &limit, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conn == nil || len(conn.Edges) != 3 {
		t.Fatalf("expected 3 edges; got %v", conn)
	}

	// Verify first edge
	if conn.Edges[0].Node.ID != "prod-1" {
		t.Errorf("edge[0].node.ID: got %q, want %q", conn.Edges[0].Node.ID, "prod-1")
	}

	// Verify second edge
	if conn.Edges[1].Node.ID != "prod-2" {
		t.Errorf("edge[1].node.ID: got %q, want %q", conn.Edges[1].Node.ID, "prod-2")
	}

	// Verify third edge
	if conn.Edges[2].Node.ID != "prod-3" {
		t.Errorf("edge[2].node.ID: got %q, want %q", conn.Edges[2].Node.ID, "prod-3")
	}

	// Verify third edge is not active
	if conn.Edges[2].Node.IsActive {
		t.Error("edge[2].node.IsActive should be false")
	}

	// Verify page info
	if conn.PageInfo == nil {
		t.Fatal("pageInfo must not be nil")
	}

	if conn.PageInfo.EndCursor == nil || *conn.PageInfo.EndCursor != "prod-3" {
		t.Errorf("endCursor: got %v, want %q", conn.PageInfo.EndCursor, "prod-3")
	}

	if !conn.PageInfo.HasNextPage {
		t.Error("hasNextPage should be true when result size == limit")
	}
}

func TestMstProductResolver_WithCursor(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	rows := &MockRows{
		data: []func(dest ...any) error{
			newProductScanFunc("prod-11", true, now),
			newProductScanFunc("prod-12", true, now.Add(1*time.Hour)),
		},
	}

	var capturedOffset int32
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, args ...any) (pgx.Rows, error) {
			if len(args) > 0 {
				if offset, ok := args[0].(int32); ok {
					capturedOffset = offset
				}
			}
			return rows, nil
		},
	}

	r := &Resolver{DBJukung: db}
	qr := r.Query()

	cursor := "10"
	limit := int32(2)
	_, err := qr.MstProduct(context.Background(), &limit, &cursor, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedOffset != 10 {
		t.Errorf("expected offset 10; got %d", capturedOffset)
	}
}

func TestMstProductResolver_WithOrderBy(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	rows := &MockRows{
		data: []func(dest ...any) error{
			newProductScanFunc("prod-1", true, now),
		},
	}

	var capturedQuery string
	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, query string, _ ...any) (pgx.Rows, error) {
			capturedQuery = query
			return rows, nil
		},
	}

	r := &Resolver{DBJukung: db}
	qr := r.Query()

	asc := model.SortOrderProductAsc
	orderBy := &model.ProductOrderByInput{
		CreatedAt: &asc,
	}

	_, err := qr.MstProduct(context.Background(), nil, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedQuery == "" {
		t.Error("query should have been captured")
	}

	// Verify that the query contains ASC ordering
	expectedFragment := "ORDER BY created_at ASC"
	if !contains(capturedQuery, expectedFragment) {
		t.Errorf("query should contain %q; got:\n%s", expectedFragment, capturedQuery)
	}
}

func TestMstProductResolver_ScanError(t *testing.T) {
	rows := &MockRows{
		data: []func(dest ...any) error{
			func(_ ...any) error {
				return fmt.Errorf("scan error")
			},
			newProductScanFunc("prod-ok", true, time.Now()),
		},
	}

	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return rows, nil
		},
	}

	r := &Resolver{DBJukung: db}
	qr := r.Query()

	conn, err := qr.MstProduct(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("resolver should not return error on scan failure; got %v", err)
	}

	if len(conn.Edges) != 1 {
		t.Errorf("expected 1 edge (scan error skipped); got %d", len(conn.Edges))
	}

	if conn.Edges[0].Node.ID != "prod-ok" {
		t.Errorf("edge[0].node.ID: got %q, want %q", conn.Edges[0].Node.ID, "prod-ok")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
