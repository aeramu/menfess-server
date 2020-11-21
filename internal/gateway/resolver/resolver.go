package resolver

import (
	"context"
	auth "github.com/aeramu/menfess-server/internal/auth/service"
	post "github.com/aeramu/menfess-server/internal/post/service"
	user "github.com/aeramu/menfess-server/internal/user/service"
	"log"
)

type Resolver interface {
	Post(req PostReq) *Post
	Posts(req ConnectionReq) *PostConnection
	Menfess() *UserConnection
	CreatePost(req CreatePostReq) *Post
	DeletePost(req DeletePostReq) string
	LikePost(req LikePostReq) *Post
	Register(req RegisterReq) string
	Login(req LoginReq) string
	Logout(req LogoutReq) string
	Me() *User
	UpdateProfile(req UpdateProfileReq) *User
	Avatars() []string
}

func NewResolver(ctx context.Context, post post.Service, auth auth.Service, user user.Service) Resolver {
	return &resolver{
		post:    post,
		auth:    auth,
		user:    user,
		Context: ctx,
	}
}

type resolver struct {
	post    post.Service
	auth    auth.Service
	user    user.Service
	Context context.Context
}

func (r *resolver) Me() *User {
	jwt := r.Context.Value("request").(map[string]string)["Authorization"]
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil{
		log.Println("Auth Service Error:", err)
		return nil
	}
	u, err := r.user.Get(user.GetReq{ID: payload.ID})
	if err != nil{
		log.Println("User Service Error:", err)
		return nil
	}
	if u == nil{
		return nil
	}
	return &User{*u, r}
}

func (r *resolver) Post(req PostReq) *Post {
	p, err := r.post.Get(post.GetReq{ID: string(req.ID)})
	if err != nil{
		log.Println("Post Service Error:", err)
		return nil
	}
	if p == nil {
		return nil
	}
	return &Post{*p, r }
}

func (r *resolver) Posts(req ConnectionReq) *PostConnection {
	first := 20
	if req.First != nil {
		first = int(*req.First)
	}
	after := ""
	if req.After != nil{
		after = string(*req.After)
	}
	posts, err := r.post.Feed(post.FeedReq{
		First: first,
		After: after,
	})
	if err != nil{
		log.Println("Post Service Error:", err)
		return nil
	}
	var postList []Post
	for _, elem := range *posts {
		postList = append(postList, Post{elem, r})
	}
	return &PostConnection{postList, r}
}

func (r *resolver) CreatePost(req CreatePostReq) *Post {
	jwt := r.Context.Value("request").(map[string]string)["Authorization"]
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil{
		log.Println("Auth Service Error:", err)
		return nil
	}
	authorID := string(req.AuthorID)
	if authorID == ""{
		authorID = payload.ID
	}
	parentID := ""
	if req.ParentID != nil{
		parentID = string(*req.ParentID)
	}
	p, err := r.post.Create(post.CreateReq{
		Body:     req.Body,
		AuthorID: authorID,
		UserID:   payload.ID,
		ParentID: parentID,
	})
	if err != nil{
		log.Println("Post Service Error:", err)
		return nil
	}
	return &Post{
		Post: *p,
		root: r,
	}
}

func (r *resolver) DeletePost(req DeletePostReq) string{
	jwt := r.Context.Value("request").(map[string]string)["Authorization"]
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil{
		log.Println("Auth Service Error:", err)
		return "Failed"
	}
	p, err := r.post.Get(post.GetReq{ID: string(req.ID)})
	if err != nil{
		log.Println("Post Service Error:", err)
		return "Failed"
	}
	if payload.ID != p.UserID{
		return "Not Authorized"
	}
	if err := r.post.Delete(post.DeleteReq{PostID: string(req.ID)}); err != nil{
		log.Println("Post Service Error:", err)
		return "Failed"
	}
	return "Success"
}

func (r *resolver) LikePost(req LikePostReq) *Post {
	jwt := r.Context.Value("request").(map[string]string)["Authorization"]
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil{
		log.Println("Auth Service Error:", err)
	}
	p, err := r.post.Like(post.LikeReq{
		PostID: string(req.ID),
		UserID: payload.ID,
	})
	if err != nil{
		log.Println("Post Service Error:", err)
		return nil
	}
	if p == nil {
		return nil
	}
	return &Post{Post: *p, root: r}
}

func (r *resolver) Register(req RegisterReq) string {
	jwt, err := r.auth.Register(auth.RegisterReq{
		Email:     req.Email,
		Password:  req.Password,
		PushToken: req.PushToken,
	})
	if err != nil{
		log.Println("Auth Service Error:", err)
	}
	return jwt
}

func (r *resolver) Login(req LoginReq) string {
	jwt, err := r.auth.Login(auth.LoginReq{
		Email:     req.Email,
		Password:  req.Password,
		PushToken: req.PushToken,
	})
	if err != nil{
		log.Println("Auth Service Error:", err)
	}
	return jwt
}

func (r *resolver) Logout(req LogoutReq) string{
	jwt := r.Context.Value("request").(map[string]string)["Authorization"]
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil{
		log.Println("Auth Service Error:", err)
		return "Failed"
	}
	if err := r.auth.Logout(auth.LogoutReq{
		ID:        payload.ID,
		PushToken: req.PushToken,
	}); err != nil{
		log.Println("Auth Service Error:", err)
		return "Failed"
	}
	return "Success"
}

func (r *resolver) Menfess() *UserConnection {
	users, err := r.user.GetMenfess(user.GetMenfessReq{})
	if err != nil {
		log.Println("User Service Error:", err)
		return nil
	}
	var userList []User
	for _, elem := range *users {
		userList = append(userList, User{elem, r})
	}
	return &UserConnection{
		users: userList,
		root:  nil,
	}
}

func (r *resolver) UpdateProfile(req UpdateProfileReq) *User {
	jwt := r.Context.Value("request").(map[string]string)["Authorization"]
	payload, err := r.auth.Auth(auth.AuthReq{Token: jwt})
	if err != nil{
		log.Println("Auth Service Error:", err)
	}
	u, err := r.user.UpdateProfile(user.UpdateProfileReq{
		ID:     payload.ID,
		Name:   req.Name,
		Avatar: req.Avatar,
		Bio:    req.Bio,
	})
	if err != nil {
		log.Println("User Service Error:", err)
		return nil
	}
	if u == nil{
		return nil
	}
	return &User{User: *u, root: r}
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

