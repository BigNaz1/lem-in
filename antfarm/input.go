package antfarm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func ResolveInput(lines []string, field *Field) error {
	for index, line := range lines {
		if line[0] != '#' {
			if isNumber(line) {
				if field.ants != nil {
					return fmt.Errorf("ants already defined")
				}

				numberOfAnts, err := strconv.Atoi(line)
				if err != nil {
					return err
				}

				err = numberToAnts(numberOfAnts, field)
				if err != nil {
					return err
				}
				continue
			}

			if strings.Contains(line, "-") {
				params := strings.Split(line, "-")
				if len(params) == 2 {
					if err := linkRooms(params[0], params[1], field); err != nil {
						return err
					}
					continue
				} else {
					return fmt.Errorf("invalid line %d", index)
				}
			}

			if strings.Contains(line, " ") {
				params := strings.Split(line, " ")
				if len(params) == 3 && index > 0 {
					isStart := false
					isEnd := false

					if lines[index-1] == "##start" {
						isStart = true
					} else if lines[index-1] == "##end" {
						isEnd = true
					}

					if err := addRoom(params[0], isStart, isEnd, field); err != nil {
						return err
					}
					continue
				}
			} else {
				return fmt.Errorf("invalid line %d", index)
			}
		}
	}

	return nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func numberToAnts(numberOfAnts int, field *Field) error {
	if numberOfAnts < 1 {
		return fmt.Errorf("invalid number of ants")
	}

	for i := 1; i <= numberOfAnts; i++ {
		ant := Ant{id: i, isFinished: false}
		field.ants = append(field.ants, &ant)
	}

	return nil
}

func addRoom(roomName string, isStart bool, isEnd bool, field *Field) error {
	if strings.Contains(roomName, " ") {
		return fmt.Errorf("room name cannot contain spaces")
	} else if roomName[0] == 'L' {
		return fmt.Errorf("room name cannot start with 'L'")
	}

	for _, room := range field.rooms {
		if room.name == roomName {
			return fmt.Errorf("Room %s already exists", roomName)
		}
	}

	room := Room{name: roomName, connectedRooms: []string{}}
	field.rooms = append(field.rooms, &room)

	if isStart {
		if field.startRoomName != "" {
			return fmt.Errorf("start room already exists")
		}
		field.startRoomName = roomName
		for _, ant := range field.ants {
			ant.currentRoom = roomName
		}
	} else if isEnd {
		if field.endRoomName != "" {
			return fmt.Errorf("end room already exists")
		}
		field.endRoomName = roomName
	}

	return nil
}

func linkRooms(firstRoom string, secondRoom string, field *Field) error {
	var firstRoomObj, secondRoomObj *Room

	for _, room := range field.rooms {
		if room.name == firstRoom {
			firstRoomObj = room
		}
		if room.name == secondRoom {
			secondRoomObj = room
		}
	}

	if firstRoomObj == nil || secondRoomObj == nil {
		return fmt.Errorf("Room %s or %s does not exist", firstRoom, secondRoom)
	}

	for _, link := range firstRoomObj.connectedRooms {
		if link == secondRoom {
			return fmt.Errorf("rooms %s and %s are already linked", firstRoom, secondRoom)
		}
	}
	for _, link := range secondRoomObj.connectedRooms {
		if link == firstRoom {
			return fmt.Errorf("rooms %s and %s are already linked", firstRoom, secondRoom)
		}
	}

	firstRoomObj.connectedRooms = append(firstRoomObj.connectedRooms, secondRoom)
	secondRoomObj.connectedRooms = append(secondRoomObj.connectedRooms, firstRoom)

	return nil
}
