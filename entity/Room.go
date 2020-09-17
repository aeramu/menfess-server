package entity

//Room interface
type Room interface {
	ID() string
	Name() string
	Avatar() string
}

type room struct {
	id     string
	name   string
	avatar string
}

//RoomConstructor constructor
type RoomConstructor struct {
	ID     string
	Name   string
	Avatar string
}

//New constructor
func (c RoomConstructor) New() Room {
	return &room{
		id:     c.ID,
		name:   c.Name,
		avatar: c.Avatar,
	}
}

func (r *room) ID() string {
	return r.id
}

func (r *room) Name() string {
	return r.name
}

func (r *room) Avatar() string {
	return r.avatar
}
