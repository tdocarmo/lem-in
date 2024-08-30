package main

import (
	"fmt"
	"strings"

	"lem-in/algorithm"
)

// Définition de la structure `Ant` pour représenter une fourmi.
type Ant struct {
	ID       int               // Identifiant unique pour la fourmi.
	Position *algorithm.Room   // Position actuelle de la fourmi.
	Path     []*algorithm.Room // Chemin assigné à la fourmi à parcourir.
	Paths    algorithm.PathPriority
}

func simulateAnts(antsCount int, startRoom, endRoom *algorithm.Room, paths []algorithm.PathPriority) {
	if len(paths) == 0 {
		fmt.Println("Aucun chemin trouvé")
		return
	}

	// Initialiser les fourmis avec leurs chemins respectifs
	ants := make([]*Ant, antsCount)
	for i := range ants {
		algorithm.SortPathsPriority(paths, antsCount)
		ants[i] = &Ant{ID: i + 1, Position: startRoom, Path: paths[0].Path}
		paths[0].Worth++
		if ants[i].Position != startRoom && ants[i].Position != endRoom {
			ants[i].Position.Occupied = true
		}
	}
	for {
		var moves []string
		startToEndMoves := 0
		startToEndAnts := make(map[int]bool)
		// var done bool
		// var save int
		for _, ant := range ants {
			// done = false
			if ant.Position != endRoom {
				nextPosIndex := getNextPositionIndex(ant)
				if nextPosIndex < len(ant.Path) {
					nextRoom := ant.Path[nextPosIndex]
					if !nextRoom.Occupied && nextRoom != endRoom || nextRoom == endRoom && ant.Position != startRoom {
						if ant.Position != startRoom {
							ant.Position.Occupied = false
						}

						ant.Position = nextRoom
						ant.Position.Occupied = true
						moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, ant.Position.Name))

						if ant.Position == endRoom {
							if ant.Position == startRoom {
								startToEndMoves++
								startToEndAnts[ant.ID] = true
							}
						}
					} else if nextRoom == endRoom && ant.Position == startRoom {
						if len(startToEndAnts) > 0 {
							continue
						}
						if ant.Position != startRoom {
							ant.Position.Occupied = false
						}
						ant.Position = nextRoom
						ant.Position.Occupied = true
						moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, ant.Position.Name))

						startToEndAnts[ant.ID] = true
					}
				}
			}
		}
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
		if allAntsFinished(ants, endRoom) {
			break
		}
		startToEndMoves = 0
	}
}

// allAntsFinished checks if all ants have reached the end room.
func allAntsFinished(ants []*Ant, endRoom *algorithm.Room) bool {
	for _, ant := range ants {
		if ant.Position != endRoom {
			return false
		}
	}
	return true
}

// getNextPositionIndex finds the index of the next position for the ant on its path.
func getNextPositionIndex(ant *Ant) int {
	for idx, room := range ant.Path {
		if room == ant.Position {
			return idx + 1
		}
	}
	return -1
}
