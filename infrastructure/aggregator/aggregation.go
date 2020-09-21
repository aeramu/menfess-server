package aggregator

import (
	"context"
	"net/http"

	graphql "github.com/aeramu/menfess-server/implementation/handler/graphql"
	mongodb "github.com/aeramu/menfess-server/implementation/mongodb/repository"
	"github.com/aeramu/menfess-server/usecase"
)

func NewHandler() http.Handler {
	ctx := context.WithValue(context.Background(), "request", map[string]string{
		"id": "5ef89baaec8ff2af8b9934c1",
	})
	repository := mongodb.New()
	defer mongodb.Disconnect()
	interactor := usecase.InteractorConstructor{
		Repository: repository,
	}.New()
	handler := graphql.New(ctx, interactor)

	return handler
}
