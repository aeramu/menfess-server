package resolver

import (
	"context"
	auth "github.com/aeramu/menfess-server/auth/service"
	post "github.com/aeramu/menfess-server/post/service"
	user "github.com/aeramu/menfess-server/user/service"
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

func (r *resolver) Me() *User{
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
