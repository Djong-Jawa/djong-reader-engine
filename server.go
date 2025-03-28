package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"sadewa-portfolio-svc/config"
	"sadewa-portfolio-svc/graph"
)

func main() {
	config.ConnectDB()
	http.Handle("/", handler.Playground("GraphQL Playground", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))

	log.Println("Server is running on port 8089")
	log.Fatal(http.ListenAndServe(":8089", nil))
}
