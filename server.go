package main

import (
	"log"
	"net/http"

	"djong-reader-engine/config"
	"djong-reader-engine/rest/controllers"
	"djong-reader-engine/graph/resolvers"
	"djong-reader-engine/rest/services"

	"github.com/joho/godotenv"
	"os"
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	
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
	config.ConnectDBJukung()
	
	// http.Handle("/query", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}})))
	// Set up GraphQL server
    // Initialize the GraphQL server
	srv := handler.NewDefaultServer(
    graph.NewExecutableSchema(
			graph.Config{Resolvers: &graph.Resolver{DB: config.DB, DBJukung: config.DBJukung}},
		),
	)

	srv.Use(extension.Introspection{})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	// srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(transport.Websocket{})

	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	// REST API routes
	// Initialize service and controller
	salesPipelineService := services.NewSalesPipelineService(config.DB)
	salesPipelineController := controllers.NewSalesPipelineController(salesPipelineService)
	
	// Register REST API endpoint
	http.HandleFunc("/api/mst/sales-pipeline", salesPipelineController.HandleSalesPipeline)



	applicationPORT := os.Getenv("APP_PORT")
	log.Println("Server runs on port : ", applicationPORT)
	log.Fatal(http.ListenAndServe(":"+applicationPORT, nil))
}
