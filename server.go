package main

import (
	"log"
	"net/http"

	"sadewa-portfolio-svc/config"
	"sadewa-portfolio-svc/graph/resolvers"

	"github.com/joho/godotenv"
	"os"
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	
)

func main() {

	// checking .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// profile log
	fmt.Println(`
					  _           _   _                _                  _ 
  __ _ _ __ __ _ _ __ | |__   __ _| | | |__   __ _  ___| | _____ _ __   __| |
 / _' | '__/ _' | '_ \| '_ \ / _' | | | '_ \ / _' |/ __| |/ / _ \ '_ \ / _' |
| (_| | | | (_| | |_) | | | | (_| | | | |_) | (_| | (__|   <  __/ | | | (_| |
 \__, |_|  \__,_| .__/|_| |_|\__, |_| |_.__/ \__,_|\___|_|\_\___|_| |_|\__,_|
 |___/          |_|             |_|                                          
 `)
	
	// connect to db
	config.ConnectDB()
	
	// http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))
	// Set up GraphQL server
    // Initialize the GraphQL server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Websocket{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(transport.UrlEncodedForm{})
	srv.AddTransport(transport.GRAPHQL{})

	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))



	applicationPORT := os.Getenv("APP_PORT")
	log.Println("Server is running on port : ", applicationPORT)
	log.Fatal(http.ListenAndServe(":"+applicationPORT, nil))
}
