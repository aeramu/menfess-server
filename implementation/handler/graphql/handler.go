package resolver

import (
	"context"
	"net/http"

	"github.com/aeramu/menfess-server/post/service"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

//Handler interface
type Handler interface {
	http.Handler
	Response(ctx context.Context, query string, operation string, variables map[string]interface{}) *graphql.Response
}

//New Handler for graphql
func New(ctx context.Context, i service.Interactor) Handler {
	schema := graphql.MustParseSchema(schemaString, &resolver{
		Context:    ctx,
		Interactor: i,
	})
	return &handler{&relay.Handler{Schema: schema}}
}

type handler struct {
	*relay.Handler
}

func (h handler) Response(ctx context.Context, query string, operation string, variables map[string]interface{}) *graphql.Response {
	return h.Schema.Exec(ctx, query, operation, variables)
}

var schemaString = `
  	schema{
		query: Query
		mutation: Mutation
  	}
  	type Query{
		menfessPost(id: ID!): MenfessPost!
		menfessPostList(first: Int, after: ID, sort: Int): MenfessPostConnection!
		menfessPostRooms(ids: [ID!]!, first: Int, after: ID): MenfessPostConnection!
		menfessRoomList: MenfessRoomConnection!
		menfessAvatarList: [String!]!
	}
	type Mutation{
		postMenfessPost(name: String!, avatar: String!, body: String!, parentID: ID, repostID: ID, roomID: ID): MenfessPost!
		upvoteMenfessPost(postID: ID!): MenfessPost!
		downvoteMenfessPost(postID: ID!): MenfessPost!
	}
	type MenfessPost{
		id: ID!
		timestamp: Int!
		name: String!
		avatar: String!
		body: String!
		replyCount: Int!
		upvoteCount: Int!
		downvoteCount: Int!
		upvoted: Boolean!
		downvoted: Boolean!
		parent: MenfessPost
		repost: MenfessPost
		child(first: Int, after: ID, before: ID, sort: Int): MenfessPostConnection!
		room: String!
	}
	type MenfessRoom{
		id: ID!
		name: String!
		avatar: String!
	}
	type MenfessPostConnection{
		edges: [MenfessPost]!
		pageInfo: PageInfo!
	}
	type MenfessRoomConnection{
		edges: [MenfessRoom]!
		pageInfo: PageInfo!
	}
	type PageInfo{
		startCursor: ID
		endCursor: ID
	}
`
