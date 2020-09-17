package usecase

import "github.com/aeramu/menfess-server/entity"

func (i *interactor) Post(id string) entity.Post {
	//TODO: sekalian ngambil child
	post := i.repo.GetPostByID(id)
	return post
}

func (i *interactor) PostFeed(first int, after string) []entity.Post {
	//TODO: first after default value in interactor, not implementation
	postList := i.repo.GetPostListByParentID("", first, after, false)
	return postList
}

func (i *interactor) PostChild(parentID string, first int, after string) []entity.Post {
	postList := i.repo.GetPostListByParentID(parentID, first, after, true)
	return postList
}

func (i *interactor) PostRooms(roomIDs []string, first int, after string) []entity.Post {
	postList := i.repo.GetPostListByRoomIDs(roomIDs, first, after, false)
	return postList
}

func (i *interactor) PostPost(name string, avatar string, body string, parentID string, repostID string, roomID string) entity.Post {
	post := i.repo.PutPost(name, avatar, body, parentID, repostID, roomID)
	return post
}

func (i *interactor) UpvotePost(accountID string, postID string) entity.Post {
	post := i.repo.GetPostByID(postID)
	if post.IsDownvoted(accountID) {
		exist := post.Downvote(accountID)
		i.repo.UpdateDownvoterIDs(postID, accountID, exist)
	}
	exist := post.Upvote(accountID)
	i.repo.UpdateUpvoterIDs(postID, accountID, exist)
	return post
}

func (i *interactor) DownvotePost(accountID string, postID string) entity.Post {
	post := i.repo.GetPostByID(postID)
	if post.IsUpvoted(accountID) {
		exist := post.Upvote(accountID)
		i.repo.UpdateUpvoterIDs(postID, accountID, exist)
	}
	exist := post.Downvote(accountID)
	i.repo.UpdateDownvoterIDs(postID, accountID, exist)
	return post
}
