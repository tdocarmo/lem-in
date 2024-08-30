// Déclaration du package 'algorithm'.
package algorithm

import "sort"

// Room définit la structure d'une salle dans la fourmilière.
type Room struct {
	Name          string  // Nom de la salle, utilisé comme identifiant unique.
	Adjacent      []*Room // Liste des salles adjacentes, représentant les tunnels vers les autres salles.
	Occupied      bool    // Indicateur pour savoir si la salle est occupée par une fourmi.
	Repeat, Total int
}

// PathPriority represents the priority of a path based on its length, number of ants, and number of repeating rooms.
type PathPriority struct {
	ID             int
	Path           []*Room
	Length         int
	Ants           int
	RepeatingRooms int
	Worth          int
}

// PathPriorities is a slice of PathPriority.
type PathPriorities []PathPriority

// Len, Less, and Swap are methods required by the sort.Interface interface.
func (p PathPriorities) Len() int {
	return len(p)
}

func (p PathPriorities) Less(i, j int) bool {
	// Prioritize shorter paths.
	if p[i].Length != p[j].Length {
		return p[i].Length < p[j].Length
	}
	// Prioritize paths with fewer repeating rooms.
	if p[i].RepeatingRooms != p[j].RepeatingRooms {
		return p[i].RepeatingRooms < p[j].RepeatingRooms
	}
	// Prioritize paths with fewer ants.
	return p[i].Ants < p[j].Ants
}

func (p PathPriorities) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// SortPathsPriority sorts the paths based on their priorities.
func SortPathsPriority(paths []PathPriority, antsCount int) {
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Worth < paths[j].Worth
	})
}

// countRepeatingRooms counts the number of rooms that appear in more than one path.
func countRepeatingRooms(path []*Room) int {
	roomCounts := make(map[string]int)
	for _, room := range path {
		roomCounts[room.Name]++
	}
	repeatingRooms := 0
	for _, count := range roomCounts {
		if count > 1 {
			repeatingRooms++
		}
	}
	return repeatingRooms
}

// FindShortestPath finds the shortest path from the 'start' room to the 'end' room using BFS.
func FindShortestPath(start, end *Room) []*Room {
	if start == end {
		return []*Room{start} // Return immediately if start is end.
	}

	// A map to track visited rooms and the path used to reach them.
	visited := make(map[*Room][]*Room)
	queue := [][]*Room{{start}}     // Start with a queue containing only the start room.
	visited[start] = []*Room{start} // Start room is considered visited with the path containing itself.

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		currentRoom := path[len(path)-1] // The last room in the current path.

		for _, adjacent := range currentRoom.Adjacent {
			if _, found := visited[adjacent]; !found {
				// Construct new path leading to this adjacent room.
				newPath := make([]*Room, len(path)+1)
				copy(newPath, path)
				newPath[len(newPath)-1] = adjacent

				if adjacent == end {
					return newPath // Return immediately if we reach the end.
				}

				visited[adjacent] = newPath    // Mark this room as visited.
				queue = append(queue, newPath) // Add new path to the queue.
			}
		}
	}

	return nil // Return nil if no path is found.
}

func NewRoom(name string) *Room {
	return &Room{Name: name, Adjacent: []*Room{}}
}

// Ajoute un lien entre deux salles.
func (r *Room) AddAdjacent(room *Room) {
	r.Adjacent = append(r.Adjacent, room)
	room.Adjacent = append(room.Adjacent, r)
}

// TransformPathsToPriorities transforms a slice of slices of *Room into a slice of PathPriority.
func TransformPathsToPriorities(paths [][]*Room, antsCountPath int) []PathPriority {
	var priorities []PathPriority
	for i, path := range paths {
		repeatingRooms := countRepeatingRooms(path)
		// Calculate Worth as the sum of the path's length and the number of ants assigned to this path.
		// Since initially, no ants are assigned, the Worth is equal to the path's length.
		worth := len(path) + antsCountPath
		priority := PathPriority{
			ID:             i + 1,
			Path:           path,
			Length:         len(path),
			Ants:           antsCountPath, // The number of ants assigned to this path.
			RepeatingRooms: repeatingRooms,
			Worth:          worth,
		}
		priorities = append(priorities, priority)
	}
	return priorities
}

// FindAllPaths finds all possible paths from the 'start' room to the 'end' room using DFS.
func FindAllPaths(start, end *Room) [][]*Room {
	var allPaths [][]*Room
	var currentPath []*Room
	findAllPathsDFS(start, end, &currentPath, &allPaths)
	return allPaths
}

// findAllPathsDFS is a helper function that performs DFS to find all paths.
func findAllPathsDFS(current *Room, end *Room, currentPath *[]*Room, allPaths *[][]*Room) {
	*currentPath = append(*currentPath, current)
	if current == end {
		// Make a copy of the current path to avoid modifying paths already added to allPaths.
		pathCopy := make([]*Room, len(*currentPath))
		copy(pathCopy, *currentPath)
		*allPaths = append(*allPaths, pathCopy)
	} else {
		for _, next := range current.Adjacent {
			if !contains(*currentPath, next) {
				findAllPathsDFS(next, end, currentPath, allPaths)
			}
		}
	}
	// Backtrack by removing the current node from the path.
	*currentPath = (*currentPath)[:len(*currentPath)-1]
}

func CountRepeatRoom(rooms []*Room) {
	for i := 0; i < len(rooms); i++ {
		for j := 0; j < len(rooms); j++ {
			if rooms[i] == rooms[j] {
				rooms[i].Repeat++
			}
		}
	}
}

func contains(path []*Room, room *Room) bool {
	for _, r := range path {
		if r == room {
			return true
		}
	}
	return false
}

func TotalRepeat(RoomRepeat []*Room) int {
	var totalRepeat int
	for i := 0; i < len(RoomRepeat); i++ {
		totalRepeat = totalRepeat + RoomRepeat[i].Repeat
	}
	return totalRepeat
}

// return the position of the path to change if there is a duplicate of the debuting room
// if there is none return length of the parameter + 1 as an error
func GetPathToChange(paths [][]*Room) int {
	PathChange := len(paths) + 1
	for i := 0; i < len(paths); i++ {
		totOfI := TotalRepeat(paths[i])
		for j := 0; j < len(paths); j++ {
			totOfJ := TotalRepeat(paths[j])
			if paths[i][1].Name == paths[j][1].Name {
				if totOfJ <= totOfI && i != j {
					PathChange = i
				} else {
					PathChange = j
				}
			}
		}
	}
	return PathChange
}

func PathsComp(paths [][]*Room, path []*Room) int {
	var saveI int
	saveI = len(paths) + 1
	for i := 0; i < len(paths); i++ {
		for j := 1; j < len(paths[i])-2; j++ {
			for k := 1; k < len(path)-2; k++ {
				if path[k].Name == paths[i][j].Name {
					if saveI == len(paths)+1 {
						saveI = i
						break
					} else if len(paths[saveI]) < len(paths[i]) {
						saveI = i
						break
					}
				}
			}
		}
	}
	return saveI
}

func TakeShortestPath(paths [][]*Room) []*Room {
	// Sort the slice of slices by length
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
	return paths[0]
}

func SortPaths(paths [][]*Room) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
}

func PathComp(path1, path2 []*Room) bool {
	var check bool
	check = true
	for i := 1; i < len(path1)-1; i++ {
		for j := 1; j < len(path2)-1; j++ {
			if path1[i].Name == path2[j].Name {
				check = false
			}
		}
	}
	return check
}

func RemoveSliceAtIndex(slices [][]*Room, index int) [][]*Room {
	return append(slices[:index], slices[index+1:]...)
}
