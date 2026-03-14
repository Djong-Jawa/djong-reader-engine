package graphqlutils

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"log"
	"encoding/json"
)

func RequestLogger(ctx context.Context, functionName string) {
	log.Printf("==== %s ==== ", functionName)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("No operation context found in context")
		}
	}()

	// Retreive the GraphQL operation context from the provided context
	opCtx := graphql.GetOperationContext(ctx)
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