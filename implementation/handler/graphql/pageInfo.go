package resolver

import (
	"github.com/graph-gophers/graphql-go"
)

//PageInfo graphql
type PageInfo interface {
	StartCursor() *graphql.ID
	EndCursor() *graphql.ID
}

type node interface {
	ID() string
}

//PageInfoResolver graphql
type pageInfo []interface{ ID() string }

//StartCursor get startcursor
func (r pageInfo) StartCursor() *graphql.ID {
	if len(r) == 0 {
		return nil
	}
	startCursor := graphql.ID(r[0].ID())
	return &startCursor
}

// EndCursor get endcursor
func (r pageInfo) EndCursor() *graphql.ID {
	if len(r) == 0 {
		return nil
	}
	endCursor := graphql.ID(r[len(r)-1].ID())
	return &endCursor
}
