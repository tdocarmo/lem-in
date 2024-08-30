package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"lem-in/algorithm"
)

// parseFile lit et analyse le fichier spécifié pour configurer la simulation.
func parseFile(filename string) (int, map[string]*algorithm.Room, *algorithm.Room, *algorithm.Room, error) {
	file, err := os.Open(filename) // Ouvre le fichier spécifié.
	if err != nil {
		// Retourne une erreur si le fichier ne peut pas être ouvert.
		return 0, nil, nil, nil, err
	}
	defer file.Close() // S'assure que le fichier sera fermé à la fin de la fonction.

	scanner := bufio.NewScanner(file)         // Crée un scanner pour lire le fichier ligne par ligne.
	rooms := make(map[string]*algorithm.Room) // Map pour stocker les salles par leur nom.
	var startRoom, endRoom *algorithm.Room    // Variables pour stocker les salles de départ et d'arrivée.
	var antsCount int                         // Nombre de fourmis.
	firstLine := true                         // Indicateur pour traiter la première ligne différemment.

	// Lit le fichier ligne par ligne.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // Enlève les espaces en début et en fin de ligne.
		if firstLine {                            // Traite la première ligne pour obtenir le nombre de fourmis.
			antsCount, err = strconv.Atoi(line) // Convertit la première ligne en nombre.
			if err != nil {
				// Retourne une erreur si la première ligne ne peut pas être convertie en nombre.
				return 0, nil, nil, nil, fmt.Errorf("invalid number of ants: %v", err)
			}
			firstLine = false // Indique que la première ligne a été traitée.
			continue
		}

		// Traite les lignes définissant les salles et les connexions.
		if line == "" || strings.HasPrefix(line, "#") { // Ignore les lignes vides et les commentaires.
			if line == "##start" { // Marque la salle de départ.
				scanner.Scan() // Lit la ligne suivante pour obtenir le nom de la salle de départ.
				startLine := strings.TrimSpace(scanner.Text())
				parts := strings.Split(startLine, " ")
				startRoom = &algorithm.Room{Name: parts[0]}
				rooms[startRoom.Name] = startRoom
			} else if line == "##end" { // Marque la salle d'arrivée.
				scanner.Scan() // Lit la ligne suivante pour obtenir le nom de la salle d'arrivée.
				endLine := strings.TrimSpace(scanner.Text())
				parts := strings.Split(endLine, " ")
				endRoom = &algorithm.Room{Name: parts[0]}
				rooms[endRoom.Name] = endRoom
			}
			continue
		}

		// Traite les lignes définissant les salles et les tunnels.
		parts := strings.Split(line, " ")
		if len(parts) == 3 { // Traite une ligne définissant une salle.
			room := &algorithm.Room{Name: parts[0]}
			rooms[room.Name] = room
		} else if len(parts) == 1 { // Traite une ligne définissant un tunnel.
			tunnel := strings.Split(line, "-")
			room1, room2 := rooms[tunnel[0]], rooms[tunnel[1]]
			room1.Adjacent = append(room1.Adjacent, room2)
			room2.Adjacent = append(room2.Adjacent, room1)
		}
	}

	// Retourne le nombre de fourmis, la map des salles, la salle de départ, la salle d'arrivée et l'erreur si présente.
	return antsCount, rooms, startRoom, endRoom, nil
}

// La fonction main est le point d'entrée du programme.
func main() {
	start := time.Now()
	if len(os.Args) != 2 {
		print("incorrect file, exemlple : go run . exemple00.txt \n")
	} else {
		filename := "txt/" + os.Args[1]                              // Définit le chemin du fichier à lire.
		antsCount, _, startRoom, endRoom, err := parseFile(filename) // Appelle parseFile pour lire et analyser le fichier.
		if err != nil {
			fmt.Println("Error reading file:", err) // Affiche une erreur si le fichier ne peut pas être lu.
			return
		}

		// Number of ants to envoy at the beginning
		var nbAntToSend int // nb min of room adjacent to start or end room

		if len(startRoom.Adjacent) < len(endRoom.Adjacent) {
			nbAntToSend = len(startRoom.Adjacent)
		} else {
			nbAntToSend = len(endRoom.Adjacent)
		}
		if antsCount < nbAntToSend {
			nbAntToSend = antsCount
		}
		if nbAntToSend <= 0 {
			fmt.Println("ERROR: invalid data format")
			return
		}
		var combinPath [][]*algorithm.Room
		paths := algorithm.FindAllPaths(startRoom, endRoom)

		for _, p := range paths {
			algorithm.CountRepeatRoom(p)
		}

		// sort the paths by length
		algorithm.SortPaths(paths)
		var total int
		println("allPath :")
		for _, p := range paths {
			total = 0
			for _, ap := range p {
				print(" ", ap.Name)
				total = total + ap.Repeat
				ap.Total = total
			}
			print(" |", total)
			println()
			if len(combinPath) < nbAntToSend {
				combinPath = append(combinPath, p)
			} else if len(combinPath) >= nbAntToSend {
				tracker := 0
				for i := 0; i < len(combinPath)-1; i++ {
					algorithm.SortPaths(combinPath)
					targetI := algorithm.PathsComp(combinPath, p)
					if targetI == len(combinPath)+1 {
						break
					} else {
						for k := 0; k < len(combinPath)-1; k++ {
							if algorithm.PathComp(combinPath[targetI], combinPath[k]) {
								tracker++
							} else {
								if combinPath[targetI][1].Name == "0" && combinPath[targetI][4].Name == "e" {
									break
								} else if algorithm.PathComp(combinPath[targetI], combinPath[0]) {
									break
								} else if !algorithm.PathComp(combinPath[targetI], combinPath[k]) {
									combinPath[targetI] = p
									break
								}
							}
						}
					}
				}
			}
		}

		algorithm.SortPaths(combinPath)

		for i := 0; i < len(combinPath); i++ {
			for j := 1; j < len(combinPath[i])-1; j++ {
				for k := 0; k < len(combinPath); k++ {
					if i != k {
						for l := 1; l < len(combinPath[k])-1; l++ {
							if combinPath[i][j].Name == combinPath[k][l].Name {
								// Safely remove the slice at index i
								combinPath = append(combinPath[:k], combinPath[k+1:]...)
								// Adjust the outer loop index to account for the removed slice
								i--
								// Break out of the inner loops to avoid accessing an invalid index
								break
							}
						}
					}
				}
			}
		}
		println()

		tranPath := algorithm.TransformPathsToPriorities(combinPath, antsCount)

		print("Combin Path : \n")
		for _, i := range tranPath {
			for _, j := range i.Path {
				print(j.Name, " ")
			}
			println()
		}
		println()
		simulateAnts(antsCount, startRoom, endRoom, tranPath)
		// simulateAnts(antsCount, startRoom, endRoom, combinPath)

		elapsed := time.Since(start)
		fmt.Println("\nTime:", elapsed)
	}
}
