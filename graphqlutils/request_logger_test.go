package graphqlutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/graphql"
)

// errMarshal is a custom type whose MarshalJSON always returns an error,
// guaranteeing the ResponseLogger error-path is covered regardless of runtime behaviour.
type errMarshal struct{}

func (errMarshal) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("intentional marshal error")
}

func TestFetchOperationContext_NoPanic(t *testing.T) {
	// Valid context — should return opCtx, nil error
	ctx := graphql.WithOperationContext(context.Background(), &graphql.OperationContext{
		RawQuery: "{ salesPipelines { id } }",
	})
	opCtx, err := fetchOperationContext(ctx)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if opCtx == nil {
		t.Error("expected non-nil opCtx")
	}
}

func TestFetchOperationContext_Panic(t *testing.T) {
	// Empty context — gqlgen panics, fetchOperationContext should recover and return an error
	_, err := fetchOperationContext(context.Background())
	if err == nil {
		t.Error("expected error when operation context is missing, got nil")
	}
}

func TestRequestLogger_NoOperationContext(t *testing.T) {
	// Exercises the err != nil branch (no op-ctx → gqlgen panics → recovered → logs error)
	RequestLogger(context.Background(), "TestFunction")
}

func TestRequestLogger_WithOperationContext(t *testing.T) {
	// Exercises the happy-path branch (op-ctx present → logs raw query)
	ctx := graphql.WithOperationContext(context.Background(), &graphql.OperationContext{
		RawQuery: "{ salesPipelines { id } }",
	})
	RequestLogger(ctx, "Query salesPipelines")
}

func TestRequestLogger_WithEmptyFunctionName(t *testing.T) {
	ctx := graphql.WithOperationContext(context.Background(), &graphql.OperationContext{
		RawQuery: "",
	})
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
	// errMarshal always returns a JSON error — guarantees the error-handling branch is covered
	ResponseLogger(errMarshal{})
}
