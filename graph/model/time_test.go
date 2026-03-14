package model

import (
	"bytes"
	"testing"
	"time"
)

func TestTime_UnmarshalGQL_ValidString(t *testing.T) {
	var mt Time
	err := mt.UnmarshalGQL("2024-01-15T10:30:00Z")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if time.Time(mt) != expected {
		t.Errorf("expected %v, got %v", expected, time.Time(mt))
	}
}

func TestTime_UnmarshalGQL_NonStringInput(t *testing.T) {
	var mt Time
	err := mt.UnmarshalGQL(12345)
	if err == nil {
		t.Error("expected error for non-string input, got nil")
	}
}

func TestTime_UnmarshalGQL_InvalidFormat(t *testing.T) {
	var mt Time
	err := mt.UnmarshalGQL("not-a-valid-time")
	if err == nil {
		t.Error("expected error for invalid time format, got nil")
	}
}

func TestTime_MarshalGQL(t *testing.T) {
	tt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	mt := Time(tt)

	var buf bytes.Buffer
	mt.MarshalGQL(&buf)

	expected := `"2024-01-15T10:30:00Z"`
	if buf.String() != expected {
		t.Errorf("expected %s, got %s", expected, buf.String())
	}
}

func TestTime_ToTime(t *testing.T) {
	tt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	mt := Time(tt)
	result := mt.ToTime()

	if result != tt {
		t.Errorf("expected %v, got %v", tt, result)
	}
}

func TestToModelTime(t *testing.T) {
	tt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	mt := ToModelTime(tt)

	if time.Time(mt) != tt {
		t.Errorf("expected %v, got %v", tt, time.Time(mt))
	}
}

func TestMarshalTime(t *testing.T) {
	tt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	marshaler := MarshalTime(tt)

	var buf bytes.Buffer
	marshaler.MarshalGQL(&buf)

	if buf.String() == "" {
		t.Error("expected non-empty output from MarshalTime")
	}
}

func TestUnmarshalTime_String(t *testing.T) {
	result, err := UnmarshalTime("2024-01-15T10:30:00Z")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestUnmarshalTime_StringPointer(t *testing.T) {
	s := "2024-01-15T10:30:00Z"
	result, err := UnmarshalTime(&s)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestUnmarshalTime_NilStringPointer(t *testing.T) {
	var s *string
	result, err := UnmarshalTime(s)
	if err != nil {
		t.Errorf("expected no error for nil pointer, got: %v", err)
	}
	if !result.IsZero() {
		t.Errorf("expected zero time for nil pointer, got %v", result)
	}
}

func TestUnmarshalTime_UnknownType(t *testing.T) {
	_, err := UnmarshalTime(99999)
	if err == nil {
		t.Error("expected error for unknown type, got nil")
	}
}

func TestTimeFromTime(t *testing.T) {
	tt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	mt := TimeFromTime(tt)

	if time.Time(mt) != tt {
		t.Errorf("expected %v, got %v", tt, time.Time(mt))
	}
}

func TestTimeFromPtr_WithValue(t *testing.T) {
	tt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	result := TimeFromPtr(&tt)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if time.Time(*result) != tt {
		t.Errorf("expected %v, got %v", tt, time.Time(*result))
	}
}

func TestTimeFromPtr_WithNil(t *testing.T) {
	result := TimeFromPtr(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}
