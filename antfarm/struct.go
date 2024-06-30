package antfarm

type Field struct {
	ants          []*Ant
	rooms         []*Room
	startRoomName string
	endRoomName   string
}

type Ant struct {
	id          int
	currentRoom string
	isFinished  bool
}

type Room struct {
	name           string
	connectedRooms []string
}
