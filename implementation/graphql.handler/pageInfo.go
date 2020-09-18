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
type pageInfo struct {
	nodeList []interface{ ID() string }
}

//StartCursor get startcursor
func (r pageInfo) StartCursor() *graphql.ID {
	if len(r.nodeList) == 0 {
		return nil
	}
	startCursor := graphql.ID(r.nodeList[0].ID())
	return &startCursor
}

// EndCursor get endcursor
func (r pageInfo) EndCursor() *graphql.ID {
	if len(r.nodeList) == 0 {
		return nil
	}
	endCursor := graphql.ID(r.nodeList[len(r.nodeList)-1].ID())
	return &endCursor
}
