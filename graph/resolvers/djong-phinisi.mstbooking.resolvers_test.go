package graph

import (
	"context"
	"fmt"
	"testing"
	"time"

	"djong-reader-engine/graph/model"

	"github.com/jackc/pgx/v5"
)

// newBookingScanFunc returns a ScanFunc that fills the eight columns returned by
// the MstBookings query:
//
//	dest[0] *string      → &bk.ID
//	dest[1] **string     → &bk.BookingCode    (nullable)
//	dest[2] **time.Time  → &bookingDate       (nullable)
//	dest[3] **time.Time  → &createdAt         (local *time.Time in resolver)
//	dest[4] **string     → &bk.CreatedBy      (nullable)
//	dest[5] **time.Time  → &bk.UpdatedAt      (nullable)
//	dest[6] **string     → &bk.UpdatedBy      (nullable)
//	dest[7] *bool        → &bk.IsActive
func newBookingScanFunc(id string, isActive bool, createdAt time.Time) func(dest ...any) error {
	return func(dest ...any) error {
		*(dest[0].(*string)) = id
		*(dest[3].(**time.Time)) = &createdAt
		*(dest[7].(*bool)) = isActive
		return nil
	}
}

// ── MstBooking (not-implemented, panics) ──────────────────────────────────────

// TestMstBookingResolver_Panics verifies that MstBooking() panics with the expected message
// and that the panic value is a non-nil error.
func TestMstBookingResolver_Panics(t *testing.T) {
	r := &Resolver{DBJukung: &MockDBQuerier{}}
	qr := r.Query()

	var panicVal any
	func() {
		defer func() { panicVal = recover() }()
		qr.MstBooking(context.Background(), "any-id") //nolint:errcheck
	}()

	if panicVal == nil {
		t.Fatal("MstBooking() must panic — got nil recover value")
	}
	panicErr, ok := panicVal.(error)
	if !ok {
		t.Fatalf("MstBooking() panicked with non-error value %T: %v", panicVal, panicVal)
	}
	expected := "not implemented: MstBooking - mstBooking"
	if panicErr.Error() != expected {
		t.Errorf("panic message: got %q, want %q", panicErr.Error(), expected)
	}
}

// ── MstBookings ───────────────────────────────────────────────────────────────

func TestMstBookingsResolver_DBError(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: errorQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	_, err := qr.MstBookings(context.Background(), nil, nil, nil)
	if err == nil {
		t.Error("expected error when DBJukung.Query fails, got nil")
	}
}

func TestMstBookingsResolver_EmptyResult(t *testing.T) {
	db := &MockDBQuerier{QueryFunc: emptyQueryFunc}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	conn, err := qr.MstBookings(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil BookingConnection")
	}
	if len(conn.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(conn.Edges))
	}
}

func TestMstBookingsResolver_Success(t *testing.T) {
	now := time.Now()
	scanFuncs := []func(dest ...any) error{
		newBookingScanFunc("1", true, now),
		newBookingScanFunc("2", true, now.Add(-time.Hour)),
	}

	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{ScanFuncs: scanFuncs}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	conn, err := qr.MstBookings(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil BookingConnection")
	}
	if len(conn.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "1" {
		t.Errorf("expected first booking ID=1, got %s", conn.Edges[0].Node.ID)
	}
	if conn.Edges[1].Node.ID != "2" {
		t.Errorf("expected second booking ID=2, got %s", conn.Edges[1].Node.ID)
	}
}

func TestMstBookingsResolver_WithPagination(t *testing.T) {
	now := time.Now()
	scanFuncs := []func(dest ...any) error{
		newBookingScanFunc("3", true, now),
	}

	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{ScanFuncs: scanFuncs}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	first := int32(1)
	after := "2"
	conn, err := qr.MstBookings(context.Background(), &first, &after, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil BookingConnection")
	}
	if len(conn.Edges) != 1 {
		t.Errorf("expected 1 edge, got %d", len(conn.Edges))
	}
	if conn.PageInfo.HasNextPage != true {
		t.Error("expected hasNextPage=true when limit reached")
	}
}

func TestMstBookingsResolver_WithOrderBy(t *testing.T) {
	now := time.Now()
	scanFuncs := []func(dest ...any) error{
		newBookingScanFunc("1", true, now.Add(-2*time.Hour)),
		newBookingScanFunc("2", true, now.Add(-time.Hour)),
	}

	db := &MockDBQuerier{
		QueryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &MockRows{ScanFuncs: scanFuncs}, nil
		},
	}
	r := &Resolver{DBJukung: db}
	qr := r.Query()

	asc := model.SortOrderBookingAsc
	orderBy := &model.BookingOrderByInput{CreatedAt: &asc}
	conn, err := qr.MstBookings(context.Background(), nil, nil, orderBy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil BookingConnection")
	}
	if len(conn.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(conn.Edges))
	}
}
