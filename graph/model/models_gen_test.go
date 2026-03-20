package model

import (
	"bytes"
	"testing"
)

func TestSortOrderSalesPipeline_IsValid_True(t *testing.T) {
	if !SortOrderSalesPipelineAsc.IsValid() {
		t.Error("expected ASC to be valid")
	}
	if !SortOrderSalesPipelineDesc.IsValid() {
		t.Error("expected DESC to be valid")
	}
}

func TestSortOrderSalesPipeline_IsValid_False(t *testing.T) {
	invalid := SortOrderSalesPipeline("INVALID")
	if invalid.IsValid() {
		t.Error("expected INVALID to not be valid")
	}
}

func TestSortOrderSalesPipeline_String(t *testing.T) {
	if SortOrderSalesPipelineAsc.String() != "ASC" {
		t.Errorf("expected ASC, got %s", SortOrderSalesPipelineAsc.String())
	}
	if SortOrderSalesPipelineDesc.String() != "DESC" {
		t.Errorf("expected DESC, got %s", SortOrderSalesPipelineDesc.String())
	}
}

func TestSortOrderSalesPipeline_UnmarshalGQL_Valid(t *testing.T) {
	var e SortOrderSalesPipeline
	if err := e.UnmarshalGQL("ASC"); err != nil {
		t.Errorf("expected no error for ASC, got: %v", err)
	}
	if e != SortOrderSalesPipelineAsc {
		t.Errorf("expected ASC, got %s", e)
	}

	if err := e.UnmarshalGQL("DESC"); err != nil {
		t.Errorf("expected no error for DESC, got: %v", err)
	}
	if e != SortOrderSalesPipelineDesc {
		t.Errorf("expected DESC, got %s", e)
	}
}

func TestSortOrderSalesPipeline_UnmarshalGQL_InvalidValue(t *testing.T) {
	var e SortOrderSalesPipeline
	if err := e.UnmarshalGQL("UNKNOWN"); err == nil {
		t.Error("expected error for invalid value, got nil")
	}
}

func TestSortOrderSalesPipeline_UnmarshalGQL_NonString(t *testing.T) {
	var e SortOrderSalesPipeline
	if err := e.UnmarshalGQL(123); err == nil {
		t.Error("expected error for non-string input, got nil")
	}
}

func TestSortOrderSalesPipeline_MarshalGQL(t *testing.T) {
	var buf bytes.Buffer
	SortOrderSalesPipelineAsc.MarshalGQL(&buf)
	if buf.String() != `"ASC"` {
		t.Errorf("expected \"ASC\", got %s", buf.String())
	}

	buf.Reset()
	SortOrderSalesPipelineDesc.MarshalGQL(&buf)
	if buf.String() != `"DESC"` {
		t.Errorf("expected \"DESC\", got %s", buf.String())
	}
}

func TestAllSortOrderSalesPipeline(t *testing.T) {
	if len(AllSortOrderSalesPipeline) != 2 {
		t.Errorf("expected 2 sort orders, got %d", len(AllSortOrderSalesPipeline))
	}
}

// ── SortOrderLead ──────────────────────────────────────────────────────────────

func TestSortOrderLead_IsValid_True(t *testing.T) {
	if !SortOrderLeadAsc.IsValid() {
		t.Error("expected ASC to be valid")
	}
	if !SortOrderLeadDesc.IsValid() {
		t.Error("expected DESC to be valid")
	}
}

func TestSortOrderLead_IsValid_False(t *testing.T) {
	invalid := SortOrderLead("INVALID")
	if invalid.IsValid() {
		t.Error("expected INVALID to not be valid")
	}
}

func TestSortOrderLead_String(t *testing.T) {
	if SortOrderLeadAsc.String() != "ASC" {
		t.Errorf("expected ASC, got %s", SortOrderLeadAsc.String())
	}
	if SortOrderLeadDesc.String() != "DESC" {
		t.Errorf("expected DESC, got %s", SortOrderLeadDesc.String())
	}
}

func TestSortOrderLead_UnmarshalGQL_Valid(t *testing.T) {
	var e SortOrderLead
	if err := e.UnmarshalGQL("ASC"); err != nil {
		t.Errorf("expected no error for ASC, got: %v", err)
	}
	if e != SortOrderLeadAsc {
		t.Errorf("expected ASC, got %s", e)
	}

	if err := e.UnmarshalGQL("DESC"); err != nil {
		t.Errorf("expected no error for DESC, got: %v", err)
	}
	if e != SortOrderLeadDesc {
		t.Errorf("expected DESC, got %s", e)
	}
}

func TestSortOrderLead_UnmarshalGQL_InvalidValue(t *testing.T) {
	var e SortOrderLead
	if err := e.UnmarshalGQL("UNKNOWN"); err == nil {
		t.Error("expected error for invalid value, got nil")
	}
}

func TestSortOrderLead_UnmarshalGQL_NonString(t *testing.T) {
	var e SortOrderLead
	if err := e.UnmarshalGQL(123); err == nil {
		t.Error("expected error for non-string input, got nil")
	}
}

func TestSortOrderLead_MarshalGQL(t *testing.T) {
	var buf bytes.Buffer
	SortOrderLeadAsc.MarshalGQL(&buf)
	if buf.String() != `"ASC"` {
		t.Errorf("expected \"ASC\", got %s", buf.String())
	}

	buf.Reset()
	SortOrderLeadDesc.MarshalGQL(&buf)
	if buf.String() != `"DESC"` {
		t.Errorf("expected \"DESC\", got %s", buf.String())
	}
}

func TestAllSortOrderLead(t *testing.T) {
	if len(AllSortOrderLead) != 2 {
		t.Errorf("expected 2 sort orders, got %d", len(AllSortOrderLead))
	}
}
