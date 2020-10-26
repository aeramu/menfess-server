package main

import (
	"context"
	"github.com/aeramu/menfess-server/gateway/server"
	"github.com/aeramu/menfess-server/post/repository"
	"log"
	"net/http"
	"os"

	"github.com/aeramu/menfess-server/post/service"
	"github.com/friendsofgo/graphiql"
)

func main() {
	//TODO: catch header, not test header
	ctx := context.WithValue(context.Background(), "request", map[string]string{
		"id": "5ef89baaec8ff2af8b9934c1",
	})

	repo := repository.NewRepository()
	//defer mongodb.Disconnect()
	service := service.NewService(repo)
	handler := server.NewServer(ctx, service)
	http.Handle("/", handler)

	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/")
	if err != nil {
		panic(err)
	}
	http.Handle("/graphiql", graphiqlHandler)

	port := getPort()
	log.Println("Server ready at " + port)
	log.Println("Graphiql ready at " + port + "/graphiql")
	log.Fatal(http.ListenAndServe(port, nil))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return ":" + port
}
