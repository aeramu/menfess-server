package main

import (
	"context"
	"log"
	"net/http"
	"os"

	resolver "github.com/aeramu/menfess-server/implementation/graphql.resolver"
	"github.com/aeramu/menfess-server/implementation/mongodb/repository"
	"github.com/aeramu/menfess-server/usecase"
	"github.com/friendsofgo/graphiql"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
	//TODO: catch header, not test header
	context := context.WithValue(context.Background(), "request", map[string]string{
		"id": "5ef89baaec8ff2af8b9934c1",
	})

	schema := graphql.MustParseSchema(resolver.Schema, &resolver.Resolver{
		Context: context,
		Interactor: usecase.InteractorConstructor{
			Repository: repository.New(),
		}.New(),
	})
	defer repository.Disconnect()
	http.Handle("/", &relay.Handler{Schema: schema})

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
