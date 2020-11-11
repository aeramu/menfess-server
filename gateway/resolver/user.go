package resolver

import (
	user "github.com/aeramu/menfess-server/user/service"
	"github.com/graph-gophers/graphql-go"
)

type User struct {
	user.User
	root *resolver
}

func (r User) ID() graphql.ID {
	return graphql.ID(r.User.ID)
}
func (r User) Name() string {
	return r.User.Name
}
func (r User) Avatar() string {
	return r.User.Avatar
}
func (r User) Bio() string {
	return r.User.Bio
}

type UserConnection struct {
	users []User
	root  *resolver
}

func (r UserConnection) Edges() []User {
	return r.users
}
func (r UserConnection) PageInfo() PageInfo {
	var nodeList PageInfo
	for _, node := range r.users {
		nodeList = append(nodeList, node)
	}
	return nodeList
}
