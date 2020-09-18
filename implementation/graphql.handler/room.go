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
	room entity.Room
	pr   *resolver
}

func (r *room) ID() graphql.ID {
	return graphql.ID(r.room.ID())
}

func (r *room) Name() string {
	return r.room.Name()
}

func (r *room) Avatar() string {
	return r.room.Avatar()
}
