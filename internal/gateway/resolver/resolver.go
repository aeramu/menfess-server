package resolver

import (
	"context"
	"log"

	auth "github.com/aeramu/menfess-server/internal/auth/service"
	post "github.com/aeramu/menfess-server/internal/post/service"
	user "github.com/aeramu/menfess-server/internal/user/service"
)

type Resolver interface {
	// Auth
	Register(req RegisterReq) (string, error)
	Login(req LoginReq) (string, error)
	Logout(ctx context.Context) (string, error)

	// User
	UpdateProfile(ctx context.Context, req UpdateProfileReq) (*User, error)
	Me(ctx context.Context) *User

	Post(req PostReq) *Post
	Posts(req ConnectionReq) *PostConnection
	Menfess() *UserConnection
	CreatePost(ctx context.Context, req CreatePostReq) *Post
	DeletePost(ctx context.Context, req DeletePostReq) string
	LikePost(ctx context.Context, req LikePostReq) *Post
	Avatars() []string
}

func NewResolver(post post.Service, auth auth.Service, user user.Service) Resolver {
	return &resolver{
		post: post,
		auth: auth,
		user: user,
	}
}

type resolver struct {
	post post.Service
	auth auth.Service
	user user.Service
}

func (r *resolver) Post(req PostReq) *Post {
	p, err := r.post.Get(post.GetReq{ID: string(req.ID)})
	if err != nil {
		log.Println("Post Service Error:", err)
		return nil
	}
	if p == nil {
		return nil
	}
	return &Post{*p, r}
}

func (r *resolver) Posts(req ConnectionReq) *PostConnection {
	first := 20
	if req.First != nil {
		first = int(*req.First)
	}
	after := ""
	if req.After != nil {
		after = string(*req.After)
	}
	posts, err := r.post.Feed(post.FeedReq{
		First: first,
		After: after,
	})
	if err != nil {
		log.Println("Post Service Error:", err)
		return nil
	}
	var postList []Post
	for _, elem := range posts {
		postList = append(postList, Post{elem, r})
	}
	return &PostConnection{postList, r}
}

func (r *resolver) CreatePost(ctx context.Context, req CreatePostReq) *Post {
	jwt := ctx.Value("Authorization").(string)
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil {
		log.Println("Auth Service Error:", err)
		return nil
	}
	authorID := string(req.AuthorID)
	if authorID == "" {
		authorID = payload.ID
	}
	parentID := ""
	if req.ParentID != nil {
		parentID = string(*req.ParentID)
	}
	p, err := r.post.Create(post.CreateReq{
		Body:     req.Body,
		AuthorID: authorID,
		UserID:   payload.ID,
		ParentID: parentID,
	})
	if err != nil {
		log.Println("Post Service Error:", err)
		return nil
	}
	return &Post{
		Post: *p,
		root: r,
	}
}

func (r *resolver) DeletePost(ctx context.Context, req DeletePostReq) string {
	jwt := ctx.Value("Authorization").(string)
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil {
		log.Println("Auth Service Error:", err)
		return "Failed"
	}
	p, err := r.post.Get(post.GetReq{ID: string(req.ID)})
	if err != nil {
		log.Println("Post Service Error:", err)
		return "Failed"
	}
	if payload.ID != p.UserID {
		return "Not Authorized"
	}
	if err := r.post.Delete(post.DeleteReq{PostID: string(req.ID)}); err != nil {
		log.Println("Post Service Error:", err)
		return "Failed"
	}
	return "Success"
}

func (r *resolver) LikePost(ctx context.Context, req LikePostReq) *Post {
	jwt := ctx.Value("Authorization").(string)
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil {
		log.Println("Auth Service Error:", err)
	}
	p, err := r.post.Like(post.LikeReq{
		PostID: string(req.ID),
		UserID: payload.ID,
	})
	if err != nil {
		log.Println("Post Service Error:", err)
		return nil
	}
	if p == nil {
		return nil
	}
	return &Post{Post: *p, root: r}
}

func (r *resolver) Menfess() *UserConnection {
	users, err := r.user.GetMenfess(user.GetMenfessReq{})
	if err != nil {
		log.Println("User Service Error:", err)
		return nil
	}
	var userList []User
	for _, elem := range users {
		userList = append(userList, User{elem, r})
	}
	return &UserConnection{
		users: userList,
		root:  nil,
	}
}

func (r *resolver) Avatars() []string {
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
