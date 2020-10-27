package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Service interface {
	Get(id string) (*Post, error)
	Feed(first int, after string) (*[]Post, error)
	Child(parentID string, first int, after string) (*[]Post, error)
	Rooms(roomID string, first int, after string) (*[]Post, error)
	Create(name string, avatar string, body string, parentID string, repostID string, roomID string) (*Post, error)
	Upvote(accountID string, postID string) (*Post, error)
	Downvote(accountID string, postID string) (*Post, error)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo Repository
}

func (i *service) Get(id string) (*Post, error) {
	post, err := i.repo.FindByID(id)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return post, nil
}

func (i *service) Feed(first int, after string) (*[]Post, error) {
	if after == ""{
		after = "ffffffffffffffffffffffff"
	}
	postList, err := i.repo.FindByParentID("", first, after, true)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return postList, nil
}

func (i *service) Child(parentID string, first int, after string) (*[]Post, error) {
	postList, err := i.repo.FindByParentID(parentID, first, after, false)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return postList, nil
}

func (i *service) Rooms(roomID string, first int, after string) (*[]Post, error) {
	if after == ""{
		after = "ffffffffffffffffffffffff"
	}
	postList, err := i.repo.FindByRoomID(roomID, first, after, true)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	return postList, nil
}

func (i *service) Create(name string, avatar string, body string, parentID string, repostID string, roomID string) (*Post, error) {
	id := primitive.NewObjectID()
	post := Post{
		ID:           id.Hex(),
		Timestamp:    int(id.Timestamp().Unix()),
		Name:         name,
		Avatar:       avatar,
		Body:         body,
		ParentID:     parentID,
		RepostID:     repostID,
		RoomID:       roomID,
		UpvoterIDs: map[string]bool{},
		DownvoterIDs: map[string]bool{},
		ReplyCount:   0,
	}
	if err := i.repo.Save(post); err != nil{
		log.Println("Repository Error:", err)
		return nil, err
	}

	parent, err := i.repo.FindByID(parentID)
	if err != nil{
		log.Println("Repository Error:", err)
		return nil, err
	}
	if parent != nil{
		parent.ReplyCount++
		i.repo.Save(*parent)
	}

	return &post, nil
}

func (i *service) Upvote(accountID string, postID string) (*Post, error) {
	post, err := i.repo.FindByID(postID)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	if post.IsDownvoted(accountID) {
		post.Downvote(accountID)
	}
	post.Upvote(accountID)
	if err := i.repo.Save(*post); err != nil{
		log.Println("Repository Error:", err)
		return nil, err
	}
	return post, nil
}

func (i *service) Downvote(accountID string, postID string) (*Post, error) {
	post, err := i.repo.FindByID(postID)
	if err != nil {
		log.Println("Repository Error:", err)
		return nil, err
	}
	if post.IsUpvoted(accountID) {
		post.Upvote(accountID)
	}
	post.Downvote(accountID)
	if err := i.repo.Save(*post); err != nil{
		log.Println("Repository Error:", err)
		return nil, err
	}
	return post, nil
}