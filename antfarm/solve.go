package antfarm

import "fmt"

func StartTurnBasedSolving(field *Field, pathsToExit [][]string) []string {
	movements := []string{}
	isSolved := false
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
						nextRoom = ""
					}
				}

				if nextRoom != "" {
					turnStruct := "%s L%v-%s"
					if turn == "" {
						turnStruct = "%sL%v-%s"
					}

					pathsUsed = append(pathsUsed, ant.currentRoom+"-"+nextRoom)
					ant.currentRoom = nextRoom
					turn = fmt.Sprintf(turnStruct, turn, ant.id, nextRoom)
				}

				isSolved = false
			}
		}

		if !isSolved {
			movements = append(movements, turn)
		}
	}

	return movements
}

func getNextRoom(currentRoom string, pathsToExit [][]string, usedPaths []string, field Field) string {
	for _, pathToExit := range pathsToExit {
		for i, room := range pathToExit {
			if room == currentRoom {
				nextRoom := pathToExit[i+1]

				if len(field.ants)-len(pathToExit) == getNumOfFinishedAnts(field.ants) && nextRoom != field.endRoomName && currentRoom == field.startRoomName && len(pathsToExit) == 2 {
					continue
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

func isRoomEmpty(roomName string, field Field) bool {
	for _, ant := range field.ants {
		if ant.currentRoom == roomName {
			return false
		}
	}

	return true
}

func getNumOfFinishedAnts(ants []*Ant) int {
	numOfFinishedAnts := 0
	for _, ant := range ants {
		if ant.isFinished {
			numOfFinishedAnts++
		}
	}

	return numOfFinishedAnts
}
