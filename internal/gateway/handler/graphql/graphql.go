package graphql

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	auth "github.com/aeramu/menfess-server/internal/auth/service"
	"github.com/aeramu/menfess-server/internal/gateway/resolver"
	post "github.com/aeramu/menfess-server/internal/post/service"
	user "github.com/aeramu/menfess-server/internal/user/service"

	"github.com/graph-gophers/graphql-go"
)

//New Handler for graphql
func NewHandler(post post.Service, auth auth.Service, user user.Service) *handler {
	f, err := os.Open("api/graphql/schema.graphql")
	if err != nil {
		log.WithError(err).Errorln("[NewHandler] Failed open graphql schema file")
		return nil
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.WithError(err).Errorln("[NewHandler] Failed read graphql schema")
		return nil
	}

	schema, err := graphql.ParseSchema(string(b), resolver.NewResolver(post, auth, user))
	if err != nil {
		log.WithError(err).Errorln("[NewHandler] Failed parse graphql schema")
		return nil
	}

	return &handler{
		Schema: schema,
	}
}

type handler struct {
	*graphql.Schema
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization"))

	response := h.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
