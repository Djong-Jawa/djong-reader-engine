package graphqlutils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
)

// fetchOperationContext safely retrieves the GraphQL operation context,
// returning an error instead of panicking when the context is missing.
func fetchOperationContext(ctx context.Context) (opCtx *graphql.OperationContext, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return graphql.GetOperationContext(ctx), nil
}

func RequestLogger(ctx context.Context, functionName string) {
	log.Printf("==== %s ==== ", functionName)

	// Retreive the GraphQL operation context from the provided context
	opCtx, err := fetchOperationContext(ctx)
	if err != nil {
		log.Printf("No operation context found in context")
		return
	}
	// log the raw query string
	log.Printf("GraphQL Query : \n%s", opCtx.RawQuery)
}

func ResponseLogger(data interface{}) {
	// Assuming 'data' is your slice of *Experience
	log.Println("Response: ")
    expJSON, err := json.MarshalIndent(data, "", "    ")
    if err != nil {
        log.Println("Error marshalling log response to JSON:", err)       
    }
    log.Println(string(expJSON))
}