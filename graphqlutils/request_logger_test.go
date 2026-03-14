package graphqlutils

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
)

func TestRequestLogger_NoOperationContext(t *testing.T) {
	// Test with empty context (no GraphQL operation context attached) — exercises the recover() path
	ctx := context.Background()
	// Should not panic
	RequestLogger(ctx, "TestFunction")
}

func TestRequestLogger_WithOperationContext(t *testing.T) {
	// Inject a real gqlgen operation context — exercises the opCtx != nil branch
	ctx := graphql.WithOperationContext(context.Background(), &graphql.OperationContext{
		RawQuery: "{ salesPipelines { id } }",
	})
	// Should not panic and should log the raw query
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
	// Channels cannot be marshalled to JSON — tests the error-handling path
	ch := make(chan int)
	// Should not panic even with unmarshalable data
	ResponseLogger(ch)
}
