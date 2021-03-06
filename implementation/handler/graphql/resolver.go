package resolver

import (
	"context"
	resolver2 "github.com/aeramu/menfess-server/gateway/resolver"

	"github.com/aeramu/menfess-server/post/service"
	"github.com/graph-gophers/graphql-go"
)

//Resolver graphql
type Resolver interface {
	MenfessPost(args struct {
		ID graphql.ID
	}) resolver2.Post
	MenfessPostList(args struct {
		First *int32
		After *graphql.ID
		Sort  *bool
	}) PostConnection
	MenfessPostRooms(args struct {
		IDs   []graphql.ID
		First *int32
		After *graphql.ID
	}) PostConnection
	MenfessRoomList() RoomConnection
	UpvoteMenfessPost(args struct {
		PostID graphql.ID
	}) resolver2.Post
	DownvoteMenfessPost(args struct {
		PostID graphql.ID
	}) resolver2.Post
	MenfessAvatarList() []string
}

type resolver struct {
	Interactor service.Interactor
	Context    context.Context
}

func (r *resolver) MenfessPost(args struct {
	ID graphql.ID
}) resolver2.Post {
	p := r.Interactor.Post(string(args.ID))
	return &resolver2.post{p, r}
}

func (r *resolver) MenfessPostList(args struct {
	First *int32
	After *graphql.ID
	Sort  *bool
}) PostConnection {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := "ffffffffffffffffffffffff"
	if args.After != nil {
		after = string(*args.After)
	}
	postList := r.Interactor.PostFeed(first, after)
	return &postConnection{postList, r}
}

func (r *resolver) MenfessPostRooms(args struct {
	IDs   []graphql.ID
	First *int32
	After *graphql.ID
}) PostConnection {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := "ffffffffffffffffffffffff"
	if args.After != nil {
		after = string(*args.After)
	}
	var roomIDs []string
	for _, id := range args.IDs {
		roomIDs = append(roomIDs, string(id))
	}
	postList := r.Interactor.PostRooms(roomIDs, first, after)
	return &postConnection{postList, r}
}

func (r *resolver) MenfessRoomList() RoomConnection {
	roomList := r.Interactor.RoomList()
	return &roomConnection{roomList, r}
}

func (r *resolver) PostMenfessPost(args struct {
	Name     string
	Avatar   string
	Body     string
	ParentID *graphql.ID
	RepostID *graphql.ID
	RoomID   *graphql.ID
}) resolver2.Post {
	parentID := ""
	if args.ParentID != nil {
		parentID = string(*args.ParentID)
	}
	repostID := ""
	if args.RepostID != nil {
		repostID = string(*args.RepostID)
	}
	roomID := ""
	if args.RoomID != nil {
		roomID = string(*args.RoomID)
	}
	p := r.Interactor.PostPost(args.Name, args.Avatar, args.Body, parentID, repostID, roomID)
	return &resolver2.post{p, r}
}

func (r *resolver) UpvoteMenfessPost(args struct {
	PostID graphql.ID
}) resolver2.Post {
	accountID := r.Context.Value("request").(map[string]string)["id"]
	p := r.Interactor.UpvotePost(accountID, string(args.PostID))
	return &resolver2.post{p, r}
}

func (r *resolver) DownvoteMenfessPost(args struct {
	PostID graphql.ID
}) resolver2.Post {
	accountID := r.Context.Value("request").(map[string]string)["id"]
	p := r.Interactor.DownvotePost(accountID, string(args.PostID))
	return &resolver2.post{p, r}
}

func (r *resolver) MenfessAvatarList() []string {
	avatarList := []string{
		"https://qiup-image.s3.amazonaws.com/avatar/avatar.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/batman.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/spiderman.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/saitama.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/kaonashi.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/mrbean.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/upin.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/ipin.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/einstein.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/monalisa.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/ronald.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/1cokelat.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/2merah.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/3vermilion.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/4oranye.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/5oranye_muda.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/6kuning.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/7hijau.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/8hijau_daun.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/9toska.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/10biru.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/11biru_tua.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/12blue-violet.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/13ungu.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/14red-violet.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/15magenta.jpg",
		"https://qiup-image.s3.amazonaws.com/avatar/16pink.jpg",
	}
	return avatarList
}
