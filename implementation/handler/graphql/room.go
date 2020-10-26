package resolver

import (
	room2 "github.com/aeramu/menfess-server/room"
	"github.com/graph-gophers/graphql-go"
)

//Room graphql
type Room interface {
	ID() graphql.ID
	Name() string
	Avatar() string
}

type room struct {
	room2.Room
	pr *resolver
}

func (r *room) ID() graphql.ID {
	return graphql.ID(r.Room.ID())
}
