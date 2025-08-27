package pathfinder

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"slices"
)

// ---- Data Structures ----
type Station struct {
	Name string
	X, Y int
}

type Graph struct {
	Stations    map[string]*Station
	Connections map[string][]string
}

// ---- Parsing ----
func ParseMapFile(path string) (*Graph, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("cannot open map file")
	}

	// Remove comments and spaces
	text := regexp.MustCompile(`#.*`).ReplaceAllString(string(file), "")
	text = regexp.MustCompile(` +`).ReplaceAllString(text, "")

	// Check "stations:" and "connections:"
	if !strings.Contains(text, "stations:") {
		return nil, fmt.Errorf("missing stations section in %q", path)
	}
	if !strings.Contains(text, "connections:") {
		return nil, fmt.Errorf("missing connections section in %q", path)
	}

	g := &Graph{
		Stations:    make(map[string]*Station),
		Connections: make(map[string][]string),
	}
	scanner := bufio.NewScanner(strings.NewReader(text))

	coords := make(map[[2]int]string)
	stationCount := 0
	countStrings := 0
	section := ""
	re := regexp.MustCompile(`^[a-z0-9_]+$`)

	for scanner.Scan() {
		countStrings++
		line := scanner.Text()
		if line == "" {
			continue
		}
		switch line {
		case "stations:":
			section = "stations"
			continue
		case "connections:":
			section = "connections"
			continue
		}

		switch section {
		case "stations":
			stationCount++
			if stationCount > 10000 {
				return nil, errors.New("the map file contains more than 10000 stations")
			}
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid station format: %q\nString number in the map file: %d", line, countStrings)
			}

			name := parts[0]
			if !re.MatchString(name) {
				return nil, fmt.Errorf("invalid station name: %q\nString number in the map file: %d", name, countStrings)
			}

			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])
			if err1 != nil || err2 != nil || x < 0 || y < 0 {
				return nil, fmt.Errorf("invalid coordinates %q. Station coordinates must be positive integers\nString number in the map file: %d", line, countStrings)
			}
			if _, exists := g.Stations[name]; exists {
				return nil, fmt.Errorf("duplicate station %q\nString number in the map file: %d", name, countStrings)
			}
			coord := [2]int{x, y}
			if _, dup := coords[coord]; dup {
				return nil, fmt.Errorf("the stations %q and \"%s,%d,%d\" have the same coordinates\nString number in the map file: %d", line, coords[coord], x, y, countStrings)
			}
			coords[coord] = name
			g.Stations[name] = &Station{name, x, y}

		case "connections":
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid connection format: %q\nString number in the map file: %d", line, countStrings)
			}
			u, v := parts[0], parts[1]
			if _, ok := g.Stations[u]; !ok {
				return nil, fmt.Errorf("unknown station %q in connection %q\nString number in the map file: %d", u, line, countStrings)
			}
			if _, ok := g.Stations[v]; !ok {
				return nil, fmt.Errorf("unknown station %q in connection %q\nString number in the map file: %d", v, line, countStrings)
			}
			// prevent duplicates
			if slices.Contains(g.Connections[u], v) {
				return nil, fmt.Errorf("duplicate connection between %q and %q\nString number in the map file: %d", u, v, countStrings)
			}
			g.Connections[u] = append(g.Connections[u], v)
			g.Connections[v] = append(g.Connections[v], u)
		}
	}
	return g, nil
}