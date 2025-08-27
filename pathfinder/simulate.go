package pathfinder

import (
	"fmt"
	"sort"
	"strings"
)

const (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)
// movement simulation
func SimulateMovements(trains []*Train) {
	if len(trains) == 0 {
		return
	}

	startStation := trains[0].Path[0]
	endStation := trains[0].Path[len(trains[0].Path)-1]

	occupied := make(map[string]string)
	occupied[startStation] = ""

	turn := 0
	fmt.Println(green + "Train movement:" + reset)

	for {
		moved := false
		var turnMoves []string
		usedEdges := make(map[string]bool)

		// Process trains furthest along their path first
		sort.Slice(trains, func(i, j int) bool {
			return trains[i].Index > trains[j].Index
		})

		for _, train := range trains {
			if !train.Active || train.Index+1 >= len(train.Path) {
				continue
			}

			current := train.Path[train.Index]
			next := train.Path[train.Index+1]

			// Free up current station if not start
			if current != startStation {
				delete(occupied, current)
			}

			edgeKey := normalizeEdgeKey(current, next)

			// Check for conflicts: edge already used or next station occupied
			if usedEdges[edgeKey] || (next != endStation && occupied[next] != "") {
				// Re-occupy current if we vacated it
				if current != startStation {
					occupied[current] = train.Name
				}
				continue
			}

			// Move train to next station
			train.Index++
			if train.Index == len(train.Path)-1 {
				train.Active = false
			}

			usedEdges[edgeKey] = true
			if next != endStation {
				occupied[next] = train.Name
			}

			turnMoves = append(turnMoves, fmt.Sprintf("%s-%s", train.Name, next))
			moved = true
		}

		if !moved {
			break
		}

		turn++
		fmt.Printf(yellow+"Turn %d:"+reset+" %s\n", turn, strings.Join(turnMoves, " "))
	}
}

// normalizeEdgeKey produces a consistent key for an undirected edge
func normalizeEdgeKey(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}