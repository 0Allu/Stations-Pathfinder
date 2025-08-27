package pathfinder

import (
	"sort"
)

// Pathfinding
func FindMultiplePaths(graph *Graph, start, end string, maxPaths int) [][]string {
	var paths [][]string
	removed := make(map[string]bool)

	neighbors := graph.Connections[start]
	// Sort neighbors by number of connections
	sort.Slice(neighbors, func(i, j int) bool {
		return len(graph.Connections[neighbors[i]]) < len(graph.Connections[neighbors[j]])
	})

	for _, nbr := range neighbors {
		if len(paths) >= maxPaths { // maxpaths == number of trains
			break
		}
		pipe := bfsFromNeighbor(graph, start, nbr, end, removed)
		if len(pipe) == 0 {
			continue
		}
		paths = append(paths, pipe)
		for _, s := range pipe[1 : len(pipe)-1] {
			removed[s] = true
		}
	}
	// Sort by length of the path
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
	return paths
}

func bfsFromNeighbor(graph *Graph, start, nbr, end string, removed map[string]bool) []string {
	type state struct {
		at   string
		path []string
	}
	q := []state{{nbr, []string{start, nbr}}}
	seen := map[string]bool{start: true, nbr: true}

	for len(q) > 0 {
		current := q[0]
		q = q[1:]
	
		if current.at == end {
			return current.path
		}
	
		for _, neighbor := range graph.Connections[current.at] {
			if removed[neighbor] || seen[neighbor] {
				continue
			}
	
			if current.at == start && neighbor != nbr {
				continue
			}
	
			seen[neighbor] = true
	
			newPath := append(append([]string{}, current.path...), neighbor)
	
			q = append(q, state{
				at:   neighbor,
				path: newPath,
			})
		}
	}
	return nil
}