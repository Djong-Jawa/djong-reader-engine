package graphqlutils

import (
	"context"
	"testing"
)

func TestRequestLogger_NoOperationContext(t *testing.T) {
	// Test with empty context (no GraphQL operation context attached)
	ctx := context.Background()
	// Should not panic
	RequestLogger(ctx, "TestFunction")
}

func TestRequestLogger_WithFunctionName(t *testing.T) {
	ctx := context.Background()
	// Various function name values should not panic
	RequestLogger(ctx, "Query salesPipelines")
	RequestLogger(ctx, "")
}

func TestResponseLogger_WithMap(t *testing.T) {
	data := map[string]string{"key": "value"}
	// Should not panic
	ResponseLogger(data)
}

func TestResponseLogger_WithNil(t *testing.T) {
	// Should not panic with nil input
	ResponseLogger(nil)
}

func TestResponseLogger_WithSlice(t *testing.T) {
	data := []string{"a", "b", "c"}
	// Should not panic
	ResponseLogger(data)
}

func TestResponseLogger_WithStruct(t *testing.T) {
	type testStruct struct {
		Name  string
		Value int
	}
	data := testStruct{Name: "test", Value: 42}
	// Should not panic
	ResponseLogger(data)
}

func TestResponseLogger_WithUnmarshalableData(t *testing.T) {
	// Channels cannot be marshalled to JSON — tests the error-handling path
	ch := make(chan int)
	// Should not panic even with unmarshalable data
	ResponseLogger(ch)
}
