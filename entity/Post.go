package entity

//Post entity interface
type Post interface {
	ID() string
	Timestamp() int
	Name() string
	Avatar() string
	Body() string
	Repost() Post
	Room() Room
	UpvoterIDs() map[string]bool
	DownvoterIDs() map[string]bool
	UpvoteCount() int
	DownvoteCount() int
	ReplyCount() int
	IsUpvoted(accountID string) bool
	IsDownvoted(accountID string) bool

	Upvote(accountID string) bool
	Downvote(accountID string) bool
}

type post struct {
	id           string
	timestamp    int
	name         string
	avatar       string
	body         string
	upvoterIDs   map[string]bool
	downvoterIDs map[string]bool
	replyCount   int
	repost       Post
	room         Room
}

//PostConstructor constructor for Post entity
type PostConstructor struct {
	ID           string
	Timestamp    int
	Name         string
	Avatar       string
	Body         string
	UpvoterIDs   map[string]bool
	DownvoterIDs map[string]bool
	ReplyCount   int
	Repost       Post
	Room         Room
}

//New construtor
func (c PostConstructor) New() Post {
	if c.UpvoterIDs == nil {
		c.UpvoterIDs = map[string]bool{}
	}
	if c.DownvoterIDs == nil {
		c.DownvoterIDs = map[string]bool{}
	}
	return &post{
		id:           c.ID,
		timestamp:    c.Timestamp,
		name:         c.Name,
		avatar:       c.Avatar,
		body:         c.Body,
		upvoterIDs:   c.UpvoterIDs,
		downvoterIDs: c.DownvoterIDs,
		replyCount:   c.ReplyCount,
		repost:       c.Repost,
		room:         c.Room,
	}
}

func (p *post) Upvote(accountID string) bool {
	voted := p.IsUpvoted(accountID)
	if !voted {
		p.upvoterIDs[accountID] = true
	} else {
		delete(p.upvoterIDs, accountID)
	}
	return voted
}
func (p *post) Downvote(accountID string) bool {
	voted := p.IsDownvoted(accountID)
	if !voted {
		p.downvoterIDs[accountID] = true
	} else {
		delete(p.downvoterIDs, accountID)
	}
	return voted
}
func (p *post) IsUpvoted(accountID string) bool {
	_, ok := p.upvoterIDs[accountID]
	return ok
}
func (p *post) IsDownvoted(accountID string) bool {
	_, ok := p.downvoterIDs[accountID]
	return ok
}
func (p *post) Room() Room {
	return p.room
}
func (p *post) Repost() Post {
	return p.repost
}
func (p *post) ReplyCount() int {
	return p.replyCount
}
func (p *post) DownvoteCount() int {
	return len(p.downvoterIDs)
}
func (p *post) UpvoteCount() int {
	return len(p.upvoterIDs)
}
func (p *post) UpvoterIDs() map[string]bool {
	return p.upvoterIDs
}
func (p *post) DownvoterIDs() map[string]bool {
	return p.downvoterIDs
}
func (p *post) Body() string {
	return p.body
}
func (p *post) Avatar() string {
	return p.avatar
}
func (p *post) Name() string {
	return p.name
}
func (p *post) Timestamp() int {
	return p.timestamp
}
func (p *post) ID() string {
	return p.id
}
