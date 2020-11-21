package server

import (
	"context"
	auth "github.com/aeramu/menfess-server/internal/auth/service"
	"github.com/aeramu/menfess-server/internal/gateway/resolver"
	user "github.com/aeramu/menfess-server/internal/user/service"
	"net/http"

	post "github.com/aeramu/menfess-server/internal/post/service"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

//Handler interface
type Handler interface {
	http.Handler
	Response(ctx context.Context, query string, operation string, variables map[string]interface{}) *graphql.Response
}

//New Handler for graphql
func NewServer(ctx context.Context, post post.Service, auth auth.Service, user user.Service) Handler {
	schema := graphql.MustParseSchema(schemaString, resolver.NewResolver(ctx, post, auth, user))
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

		me: User
		post(id: ID!): Post
		posts(first: Int, after: ID): PostConnection!
		menfess: UserConnection!
		avatars: [String!]!
	}
	type Mutation{
		register(email: String!, password: String!, pushToken: String!): String!
		login(email: String!, password: String!, pushToken: String!): String!
		logout(id: ID!, pushToken: String!): String!
		updateProfile(name: String!, avatar: String!, bio: String!): User
		createPost(body: String!, authorID: ID!, parentID: ID): Post
		likePost(id: ID!): Post
	}
	type Post{
		id: ID!
		timestamp: Int!
		body: String!
		author: User
		likesCount: Int!
		liked: Boolean!
		repliesCount: Int!
		replies(first: Int, after: ID): PostConnection!
	}
	type User{
		id: ID!
		name: String!
		avatar: String!
		bio: String!
	}
	type PostConnection{
		edges: [Post!]!
		pageInfo: PageInfo!
	}
	type UserConnection{
		edges: [User!]!
		pageInfo: PageInfo!
	}
	type PageInfo{
		startCursor: ID
		endCursor: ID
	}
`

