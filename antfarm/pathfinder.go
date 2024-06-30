package antfarm

import "sort"

func FindShortestPaths(field Field, startRoom string, currentPath []string, visited map[string]bool, allPaths *[][]string) {
	if startRoom == "" {
		startRoom = field.startRoomName
	}

	endRoom := field.endRoomName
	rooms := field.rooms

	currentPath = append(currentPath, startRoom)
	visited[startRoom] = true

	var room *Room
	for _, cRoom := range rooms {
		if cRoom.name == startRoom {
			room = cRoom
			break
		}
	}

	if startRoom == endRoom {
		*allPaths = append(*allPaths, append([]string{}, currentPath...))
		sortByShortest(*allPaths)
	} else {
		for _, neighbor := range room.connectedRooms {
			if !visited[neighbor] {
				FindShortestPaths(field, neighbor, currentPath, visited, allPaths)
			}
		}
	}

	visited[startRoom] = false
}

func sortByShortest(slice [][]string) {
	sort.Slice(slice, func(i, j int) bool {
		return len(slice[i]) < len(slice[j])
	})
}

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

func isPathWithSameStartExists(paths [][]string, path []string) bool {
	for _, p := range paths {
		if p[1] == path[1] {
			return true
		}
	}
	return false
}

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
