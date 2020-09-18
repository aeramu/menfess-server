package main

import (
	"context"
	"log"
	"net/http"
	"os"

	graphql "github.com/aeramu/menfess-server/implementation/graphql.handler"
	mongodb "github.com/aeramu/menfess-server/implementation/mongodb/repository"
	"github.com/aeramu/menfess-server/usecase"
	"github.com/friendsofgo/graphiql"
)

func main() {
	//TODO: catch header, not test header
	ctx := context.WithValue(context.Background(), "request", map[string]string{
		"id": "5ef89baaec8ff2af8b9934c1",
	})

	repository := mongodb.New()
	defer mongodb.Disconnect()
	interactor := usecase.InteractorConstructor{
		Repository: repository,
	}.New()
	handler := graphql.New(ctx, interactor)
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
