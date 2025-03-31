package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"sadewa-portfolio-svc/config"
	"sadewa-portfolio-svc/graph"
	"github.com/joho/godotenv"
	"os"
	"fmt"
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
	http.Handle("/", handler.Playground("GraphQL Playground", "/query"))
	http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))

	applicationPORT := os.Getenv("APP_PORT")
	log.Println("Server is running on port : ", applicationPORT)
	log.Fatal(http.ListenAndServe(":"+applicationPORT, nil))
}
