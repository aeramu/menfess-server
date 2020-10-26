package service

type Post struct {
	ID				string
	Timestamp    int
	Name         string
	Avatar       string
	Body         string
	ParentID 	string
	RepostID       string
	RoomID         string
	UpvoterIDs		map[string]bool
	DownvoterIDs	map[string]bool
	ReplyCount   int
}

func (p *Post) Upvote(accountID string) bool {
	voted := p.IsUpvoted(accountID)
	if !voted {
		p.UpvoterIDs[accountID] = true
	} else {
		delete(p.UpvoterIDs, accountID)
	}
	return voted
}
func (p *Post) Downvote(accountID string) bool {
	voted := p.IsDownvoted(accountID)
	if !voted {
		p.DownvoterIDs[accountID] = true
	} else {
		delete(p.DownvoterIDs, accountID)
	}
	return voted
}
func (p Post) IsUpvoted(accountID string) bool {
	_, ok := p.UpvoterIDs[accountID]
	return ok
}
func (p Post) IsDownvoted(accountID string) bool {
	_, ok := p.DownvoterIDs[accountID]
	return ok
}
func (p Post) DownvoteCount() int {
	return len(p.DownvoterIDs)
}
func (p Post) UpvoteCount() int {
	return len(p.DownvoterIDs)
}