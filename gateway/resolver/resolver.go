package resolver

import (
	"context"
	post "github.com/aeramu/menfess-server/post/service"
	room "github.com/aeramu/menfess-server/room/service"
	"github.com/graph-gophers/graphql-go"
)

type Resolver interface {
	MenfessPost(args MenfessPostRequest) *Post
	MenfessPostList(args ConnectionRequest) *PostConnection
	MenfessPostRooms(args MenfessPostRoomsRequest) *PostConnection
	UpvoteMenfessPost(args UpvoteMenfessPostRequest) *Post
	DownvoteMenfessPost(args DownvoteMenfessPostRequest) *Post
	// MenfessRoomList() RoomConnection
	MenfessAvatarList() []string
}

func NewResolver(ctx context.Context, post post.Service, room room.Service) Resolver {
	return &resolver{
		post:    post,
		room: room,
		Context: ctx,
	}
}

type resolver struct{
	post    post.Service
	room room.Service
	Context context.Context
}

func (r *resolver) MenfessPost(args MenfessPostRequest) *Post {
	p, err := r.post.Get(string(args.ID))
	if err != nil{
		return nil
	}
	return &Post{*p, r}
}

func (r *resolver) MenfessPostList(args ConnectionRequest) *PostConnection {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := ""
	if args.After != nil {
		after = string(*args.After)
	}
	postList, err := r.post.Feed(first, after)
	if err != nil{
		return nil
	}
	var posts []Post
	for _, elem := range *postList {
		posts = append(posts, Post{elem, r})
	}
	return &PostConnection{posts, r}
}

func (r *resolver) MenfessPostRooms(args MenfessPostRoomsRequest) *PostConnection {
	first := 20
	if args.First != nil {
		first = int(*args.First)
	}
	after := ""
	if args.After != nil {
		after = string(*args.After)
	}
	var roomIDs []string
	for _, id := range args.IDs {
		roomIDs = append(roomIDs, string(id))
	}
	postList, err := r.post.Rooms(roomIDs[0], first, after)
	if err != nil{
		return nil
	}
	var posts []Post
	for _, elem := range *postList {
		posts = append(posts, Post{elem, r})
	}
	return &PostConnection{posts, r}
}

func (r *resolver) MenfessRoomList() *RoomConnection {
	roomList, err := r.room.GetList(room.GetListReq{})
	if err != nil{
		return nil
	}
	var rooms []Room
	for _, elem := range *roomList {
		rooms = append(rooms, Room{elem, r})
	}
	return &RoomConnection{rooms, r}
}

func (r *resolver) PostMenfessPost(args struct {
	Name     string
	Avatar   string
	Body     string
	ParentID *graphql.ID
	RepostID *graphql.ID
	RoomID   *graphql.ID
}) *Post {
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
	p, err := r.post.Create(args.Name, args.Avatar, args.Body, parentID, repostID, roomID)
	if err != nil{
		return nil
	}
	return &Post{*p, r}
}

func (r *resolver) UpvoteMenfessPost(args UpvoteMenfessPostRequest) *Post {
	accountID := r.Context.Value("request").(map[string]string)["id"]
	p, err := r.post.Upvote(accountID, string(args.PostID))
	if err != nil{
		return nil
	}
	return &Post{*p, r}
}

func (r *resolver) DownvoteMenfessPost(args DownvoteMenfessPostRequest) *Post {
	accountID := r.Context.Value("request").(map[string]string)["id"]
	p, err := r.post.Downvote(accountID, string(args.PostID))
	if err != nil {
		return nil
	}
	return &Post{*p, r}
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
