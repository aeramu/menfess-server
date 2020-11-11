package service

type Post struct {
	ID	string
	Timestamp int
	Body string
	AuthorID string
	UserID string
	ParentID string
	LikeIDs map[string]bool
	RepliesCount int
	ReportsCount int
}

func (p *Post) Like(accountID string) bool {
	voted := p.IsLiked(accountID)
	if !voted {
		p.LikeIDs[accountID] = true
	} else {
		delete(p.LikeIDs, accountID)
	}
	return voted
}
func (p Post) IsLiked(accountID string) bool {
	_, ok := p.LikeIDs[accountID]
	return ok
}
func (p Post) LikesCount() int {
	return len(p.LikeIDs)
}