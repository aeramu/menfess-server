package resolver

import (
	"github.com/aeramu/menfess-server/entity"
	"github.com/graph-gophers/graphql-go"
)

//Room graphql
type Room interface {
	ID() graphql.ID
	Name() string
	Avatar() string
}

type room struct {
	entity.Room
	pr *resolver
}

func (r *room) ID() graphql.ID {
	return graphql.ID(r.Room.ID())
}
