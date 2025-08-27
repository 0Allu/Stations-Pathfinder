package main

import (
	"fmt"
	"os"
	"pathfinder/pathfinder"
	"strconv"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	yellow = "\033[33m"
	reset = "\033[0m"
)

func main() {
	for i := range os.Args {
		if os.Args[i] == "-h" || os.Args[i] == "--help" {
			help()
			return
		}
	}

	if len(os.Args) != 5 {
		exitWithError("Incorrect number of arguments.", true)
	}
	if os.Args[4] == "-g" {
		pathfinder.Generator() // Generate a map file
		return
	}

	mapFile := os.Args[1]
	start := os.Args[2]
	end := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains < 0 {
		exitWithError("Number of trains must be a positive integer", false)
	}
	if numTrains == 0 {
		exitWithError("Number of trains must be greater than 0", false)
	}

	graph, err := pathfinder.ParseMapFile(mapFile)
	if err != nil {
		exitWithError(fmt.Sprintf("Error parsing map: %s", err), false)
	}
	if _, ok := graph.Stations[start]; !ok {
		exitWithError(fmt.Sprintf("Start station, %q does not exist", start), false)
	}
	if _, ok := graph.Stations[end]; !ok {
		exitWithError(fmt.Sprintf("End station, %q does not exist", end), false)
	}
	if start == end {
		exitWithError(fmt.Sprintf("Start and end stations, %q and %q are the same", start, end), false)
	}

	paths := pathfinder.FindMultiplePaths(graph, start, end, numTrains)

	if len(paths) == 0 {
		exitWithError(fmt.Sprintf("No path between %q and %q stations.", start, end), false)
	}

	fmt.Println(green+"Paths found:"+reset)
	for i, path := range paths {
		fmt.Printf(green+"Path %d:"+reset+" %s\n", i+1, strings.Join(path, " -> "))
	}
	fmt.Println()

	trains := pathfinder.AssignToPipelines(paths, numTrains)
	pathfinder.SimulateMovements(trains)
}

func help() {
	fmt.Println("To find train routes, use: go run . [path to file containing network map] [start station] [end station] [number of trains]")
	fmt.Println("To generate a map file, use: go run main.go [txt file] [map file] [number of stations] -g")
	fmt.Println("Run error tests: go test -v")
}

func exitWithError(msg string, showHelp bool) {
	fmt.Fprintln(os.Stderr, red+"Error: "+reset, yellow+msg+reset)
	if showHelp {
		help()
	}
	os.Exit(1)
}
