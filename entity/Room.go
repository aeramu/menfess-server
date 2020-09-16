package entity

//Room interface
type Room interface {
	ID() string
	Name() string
}

//RoomConstructor constructor
type RoomConstructor struct {
	ID   string
	Name string
}

//New constructor
func (c RoomConstructor) New() Room {
	return &room{
		id:   c.ID,
		name: c.Name,
	}
}

type room struct {
	id   string
	name string
}

func (r *room) ID() string {
	return r.id
}

func (r *room) Name() string {
	return r.name
}
