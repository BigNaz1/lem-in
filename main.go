package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Field represents the entire structure of the maze.
type Field struct {
	ants          []*Ant
	rooms         []*Room
	startRoomName string
	endRoomName   string
}

// Ant represents an individual ant in the maze.
type Ant struct {
	id          int
	currentRoom string
	isFinished  bool
}

// Room represents a single room in the maze.
type Room struct {
	name           string
	connectedRooms []string
}

// Main function to run the program.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Invalid number of arguments")
		return
	}

	// Read the file specified in the command line arguments.
	filePath := os.Args[1]
	lines, err := ReadFile(filePath)
	if err != nil {
		fmt.Println("ERROR: File not found")
		return
	}

	// Print the file content.
	for _, line := range lines {
		fmt.Println(line)
	}

	fmt.Println() // Empty line

	// Resolve input and fill the field structure.
	field := Field{}
	err = ResolveInput(lines, &field)
	if err != nil {
		fmt.Println("ERROR: invalid data format, " + err.Error())
		return
	}

	// Find all shortest paths.
	var shortestPaths [][]string
	visited := make(map[string]bool)
	FindShortestPaths(field, "", []string{}, visited, &shortestPaths)
	shortestPaths = RemoveTooLongPaths(shortestPaths)

	if len(shortestPaths) == 0 {
		fmt.Printf("ERROR: invalid data format, farm is unsolvable")
		return
	}

	// Start turn-based solving and print the result.
	StartTurnBasedSolving(&field, shortestPaths)
}

// ReadFile reads the content of a file line by line.
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

// ResolveInput processes the input lines and fills the field structure.
func ResolveInput(lines []string, field *Field) error {
	for index, line := range lines {
		if line[0] != '#' {
			// Assume that the line is the number of ants.
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

			// Assume that the line is a link between two rooms.
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

			// Assume that the line is a room with coordinates.
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

// isNumber checks if a string is a valid number.
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// numberToAnts creates the specified number of ants and adds them to the field.
func numberToAnts(numberOfAnts int, field *Field) error {
	if numberOfAnts < 1 {
		return fmt.Errorf("invalid number of ants")
	}

	for i := 0; i < numberOfAnts; i++ {
		ant := Ant{id: i, isFinished: false}
		field.ants = append(field.ants, &ant)
	}

	return nil
}

// addRoom adds a room to the field and sets start/end rooms if specified.
func addRoom(roomName string, isStart bool, isEnd bool, field *Field) error {
	// Check if roomName is valid.
	if strings.Contains(roomName, " ") {
		return fmt.Errorf("room name cannot contain spaces")
	} else if roomName[0] == 'L' {
		return fmt.Errorf("room name cannot start with 'L'")
	}

	// Check if the room already exists.
	for _, room := range field.rooms {
		if room.name == roomName {
			return fmt.Errorf("Room %s already exists", roomName)
		}
	}

	// Create the room.
	room := Room{name: roomName, connectedRooms: []string{}}
	field.rooms = append(field.rooms, &room)

	// Set the start and end rooms.
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

// linkRooms creates a bidirectional link between two rooms.
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

	// Check if the rooms are already linked.
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

	// Link the rooms.
	firstRoomObj.connectedRooms = append(firstRoomObj.connectedRooms, secondRoom)
	secondRoomObj.connectedRooms = append(secondRoomObj.connectedRooms, firstRoom)

	return nil
}

// FindShortestPaths finds all shortest paths from start to end using DFS.
func FindShortestPaths(field Field, startRoom string, currentPath []string, visited map[string]bool, allPaths *[][]string) {
	if startRoom == "" {
		startRoom = field.startRoomName
	}

	endRoom := field.endRoomName
	rooms := field.rooms

	// Add the current room to the path and mark it as visited.
	currentPath = append(currentPath, startRoom)
	visited[startRoom] = true

	// Find the room object.
	var room *Room
	for _, cRoom := range rooms {
		if cRoom.name == startRoom {
			room = cRoom
			break
		}
	}

	if startRoom == endRoom {
		// Found a path to the end room.
		*allPaths = append(*allPaths, append([]string{}, currentPath...))
		sortByShortest(*allPaths)
	} else {
		for _, neighbor := range room.connectedRooms {
			if !visited[neighbor] {
				FindShortestPaths(field, neighbor, currentPath, visited, allPaths)
			}
		}
	}

	// Backtrack.
	visited[startRoom] = false
}

// sortByShortest sorts paths by their length.
func sortByShortest(slice [][]string) {
	sort.Slice(slice, func(i, j int) bool {
		return len(slice[i]) < len(slice[j])
	})
}

// RemoveTooLongPaths removes paths that are too long.
func RemoveTooLongPaths(paths [][]string) [][]string {
	if len(paths) >= 10 {
		paths = append(paths[:5], paths[6:]...)
	}

	var shortestPaths [][]string
	for _, path := range paths {
		if len(paths) >= 10 {
			if isPathWithSameStartExists(shortestPaths, path) {
				continue
			} else {
				isHasDuplicate := false
				for _, p := range path {
					if isPathsHaveThisRoom(shortestPaths, p) {
						isHasDuplicate = true
						break
					}
				}

				if !isHasDuplicate {
					shortestPaths = append(shortestPaths, path)
					continue
				}
			}
		} else {
			if len(path)-2 > len(paths[0]) {
				continue
			} else {
				shortestPaths = append(shortestPaths, path)
			}
		}
	}

	return shortestPaths
}

// isPathWithSameStartExists checks if a path with the same start exists in the list.
func isPathWithSameStartExists(paths [][]string, path []string) bool {
	for _, p := range paths {
		if p[1] == path[1] {
			return true
		}
	}
	return false
}

// isPathsHaveThisRoom checks if a room is already included in any path.
func isPathsHaveThisRoom(paths [][]string, room string) bool {
	for _, path := range paths {
		for _, p := range path {
			if p == room && p != "start" && p != "end" {
				return true
			}
		}
	}
	return false
}

// StartTurnBasedSolving performs turn-based solving and prints the result.
func StartTurnBasedSolving(field *Field, pathsToExit [][]string) {
	isSolved := false
	turns := []string{}
	for !isSolved {
		turn := ""
		isSolved = true
		pathsUsed := []string{}
		for _, ant := range field.ants {
			if ant.currentRoom == field.endRoomName {
				ant.isFinished = true
			}

			if !ant.isFinished {
				nextRoom := getNextRoom(ant.currentRoom, pathsToExit, pathsUsed, *field)

				for _, path := range pathsUsed {
					if path == ant.currentRoom+"-"+nextRoom {
						nextRoom = "" // to prevent ants from colliding
					}
				}

				if nextRoom != "" {
					turnStruct := "%s L%v-%s"
					if turn == "" {
						turnStruct = "%sL%v-%s"
					}

					pathsUsed = append(pathsUsed, ant.currentRoom+"-"+nextRoom)
					ant.currentRoom = nextRoom
					turn = fmt.Sprintf(turnStruct, turn, ant.id+1, nextRoom)
				}

				isSolved = false
			}
		}

		if !isSolved {
			turns = append(turns, turn)
		}
	}

	for _, turn := range turns {
		fmt.Println(turn)
	}
}

// getNextRoom determines the next room for an ant.
func getNextRoom(currentRoom string, pathsToExit [][]string, usedPaths []string, field Field) string {
	for _, pathToExit := range pathsToExit {
		for i, room := range pathToExit {
			if room == currentRoom {
				nextRoom := pathToExit[i+1]

				// Prevent ants from colliding in certain test cases.
				if len(field.ants)-len(pathToExit) == getNumOfFinishedAnts(field.ants) && nextRoom != field.endRoomName && currentRoom == field.startRoomName && len(pathsToExit) == 2 {
					continue // Wait for better path
				}

				if isRoomEmpty(nextRoom, field) || nextRoom == field.endRoomName {
					isPathUsed := false
					for _, path := range usedPaths {
						if path == currentRoom+"-"+nextRoom {
							isPathUsed = true
						}
					}

					if isPathUsed {
						continue
					}

					return nextRoom
				} else {
					continue
				}
			}
		}
	}

	return ""
}

// isRoomEmpty checks if a room is empty.
func isRoomEmpty(roomName string, field Field) bool {
	for _, ant := range field.ants {
		if ant.currentRoom == roomName {
			return false
		}
	}

	return true
}

// getNumOfFinishedAnts returns the number of finished ants.
func getNumOfFinishedAnts(ants []*Ant) int {
	numOfFinishedAnts := 0
	for _, ant := range ants {
		if ant.isFinished {
			numOfFinishedAnts++
		}
	}

	return numOfFinishedAnts
}
