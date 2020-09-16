package main

import (
	"context"
	"log"
	"net/http"

	resolver "github.com/aeramu/menfess-server/implementation/graphql.resolver"
	"github.com/aeramu/menfess-server/implementation/mongodb/repository"
	"github.com/aeramu/menfess-server/usecase"
	"github.com/friendsofgo/graphiql"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {
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

	log.Println("Server ready at 8000")
	log.Println("Graphiql ready at 8000/graphiql")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
