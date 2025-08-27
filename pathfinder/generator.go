package pathfinder

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Generator creates a map file with random stations and connections
func Generator() {

	txtFile := os.Args[1]
	mapFile := os.Args[2]

	// Check number of stations
	numStations, err := strconv.Atoi(os.Args[3])
	if err != nil || numStations < 2 {
		exitWithError(fmt.Sprintf("Invalid number of stations: %s (must be at least 2)\n", os.Args[3]))
	}

	// Get stations from the input file
	stations, err := parseStations(txtFile, numStations)
	if err != nil {
		exitWithError(fmt.Sprintf("Error processing input file: %v\n", err))
	}

	// Generate coordinates and connections
	coords := generateCoordinates(stations, numStations)
	connections := generateConnections(stations)

	// Save results to file
	err = saveToMapFile(mapFile, coords, connections)
	if err != nil {
		exitWithError(fmt.Sprintf("Error saving .map file: %v\n", err))
	}

	fmt.Printf(".map file successfully created: %s with %d stations.\n", mapFile, len(coords))
}

// Read TXT and extract name of stations, generate new if needed.
func parseStations(filePath string, numStations int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	unique := make(map[string]struct{})
	var stations []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		if name != "" {
			if _, exists := unique[name]; !exists {
				stations = append(stations, name)
				unique[name] = struct{}{}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for len(stations) < numStations {
		newName := fmt.Sprintf("Station%d", len(stations)+1)
		if _, exists := unique[newName]; !exists {
			stations = append(stations, newName)
			unique[newName] = struct{}{}
		}
	}
	if len(stations) > numStations {
		stations = stations[:numStations]
	}
	return stations, nil
}

// Generate random coordinates, avoiding duplicates.
func generateCoordinates(stations []string, numStations int) map[string][2]int {
	coords := make(map[string][2]int)
	used := make(map[string]struct{})

	for _, station := range stations {
		for {
			x, y := rand.Intn(numStations*20), rand.Intn(numStations*20)
			coordKey := fmt.Sprintf("%d,%d", x, y)
			if _, exists := used[coordKey]; !exists {
				coords[station] = [2]int{x, y}
				used[coordKey] = struct{}{}
				break
			}
		}
	}
	return coords
}

// Check for connection in the list (including reverse ones)
func isDuplicateConnection(connections []string, a, b string) bool {
	connectionAB := fmt.Sprintf("%s-%s", a, b)
	connectionBA := fmt.Sprintf("%s-%s", b, a)

	for _, conn := range connections {
		if conn == connectionAB || conn == connectionBA {
			return true
		}
	}
	return false
}

// Generate unique connections between stations
func generateConnections(stations []string) []string {
	// Number of desired connections
	targetConnections := len(stations) * 2
	connections := make([]string, 0)

	// Limit the number of attempts to avoid infinite loop
	maxAttempts := targetConnections * 10
	attempts := 0

	// Use an iterator to ensure minimum connection for each station
	for i := 0; i < len(stations) && len(connections) < targetConnections; i++ {
		// Guarantee that each station has at least one connection
		connected := false
		for j := 0; j < 3 && !connected; j++ { // Try up to 3 times for each station
			b := rand.Intn(len(stations))
			if i != b && !isDuplicateConnection(connections, stations[i], stations[b]) {
				connections = append(connections, fmt.Sprintf("%s-%s", stations[i], stations[b]))
				connected = true
			}
		}
	}

	// Add remaining connections randomly
	for len(connections) < targetConnections && attempts < maxAttempts {
		a := rand.Intn(len(stations))
		b := rand.Intn(len(stations))
		if a != b && !isDuplicateConnection(connections, stations[a], stations[b]) {
			connections = append(connections, fmt.Sprintf("%s-%s", stations[a], stations[b]))
		}
		attempts++
	}

	return connections
}

// Save to `.map` file.
func saveToMapFile(filename string, stations map[string][2]int, connections []string) error {
	if !strings.HasSuffix(filename, ".map") {
		filename += ".map"
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write stations section
	if _, err := file.WriteString("stations:\n"); err != nil {
		return err
	}

	for name, coord := range stations {
		line := fmt.Sprintf("%s,%d,%d\n", name, coord[0], coord[1])
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}

	// Write connections section
	if _, err := file.WriteString("\nconnections:\n"); err != nil {
		return err
	}

	for _, connection := range connections {
		if _, err := file.WriteString(connection + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, red+"Error: "+reset, yellow+msg+reset)
	os.Exit(1)
}