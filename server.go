package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"sadewa-portfolio-svc/config"
	"sadewa-portfolio-svc/graph"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	http.Handle("/", handler.Playground("GraphQL Playground", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))

	applicationPORT := os.Getenv("APP_PORT")
	log.Println("Server is running on port : ", applicationPORT)
	log.Fatal(http.ListenAndServe(":"+applicationPORT, nil))
}
