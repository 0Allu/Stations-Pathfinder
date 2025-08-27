# Stations Pathfinder

A path-finding algorithm to find the most efficient paths to move trains from one destination to another.

## Overview

This program implements a sophisticated pathfinding system for train networks. It can:
- Parse station network maps from `.map` files
- Find multiple optimal paths between stations
- Simulate train movements with collision avoidance
- Generate random station networks for testing
- Handle multiple trains efficiently using different routes

## Features

- **Multi-path Algorithm**: Finds multiple non-overlapping paths to distribute train traffic
- **Collision Detection**: Prevents trains from occupying the same station or edge simultaneously
- **Smart Train Assignment**: Distributes trains across paths based on path length and train count
- **Network Generation**: Creates random station networks for testing and simulation
- **Comprehensive Validation**: Validates map files for format errors, duplicate stations, and invalid connections

## Usage

### Finding Train Routes

```bash
go run . <map_file> <start_station> <end_station> <number_of_trains>
```

**Parameters:**
- `map_file`: Path to the `.map` file containing the station network
- `start_station`: Name of the starting station
- `end_station`: Name of the destination station  
- `number_of_trains`: Number of trains to route (positive integer)

**Example:**
```bash
go run . testdata/small.map small large 3
```

### Generating Map Files

```bash
go run . <input_txt_file> <output_map_file> <number_of_stations> -g
```

**Parameters:**
- `input_txt_file`: Text file containing station names (one per line)
- `output_map_file`: Output filename for the generated `.map` file
- `number_of_stations`: Number of stations to include in the network
- `-g`: Generator flag

**Example:**
```bash
go run . stations.txt network.map 20 -g
```

### Help

```bash
go run . -h
# or
go run . --help
```

### Running Tests

```bash
go test -v
```

## Map File Format

Map files use a specific format with two sections:

```
stations:
station_name,x_coordinate,y_coordinate
another_station,x_coordinate,y_coordinate

connections:
station_name-another_station
```

### Stations Section
- Format: `name,x,y`
- Station names: lowercase letters, numbers, and underscores only
- Coordinates: positive integers
- No duplicate names or coordinates allowed
- Maximum 10,000 stations per file

### Connections Section
- Format: `station1-station2`
- Both stations must be defined in the stations section
- Connections are bidirectional
- No duplicate connections allowed

### Comments
Lines starting with `#` are treated as comments and ignored.

**Example Map File:**
```
# Simple network example
stations:
start,0,0
middle,1,1
end,2,0

connections:
start-middle
middle-end
```

## Algorithm Details

### Pathfinding Strategy
1. **Multi-path Discovery**: Uses BFS(Breadth First Search) from different neighbors of the start station
2. **Station Removal**: Removes intermediate stations from used paths to ensure non-overlapping routes
3. **Path Sorting**: Prioritizes shorter paths for better efficiency

### Train Movement Simulation
1. **Turn-based Movement**: Trains move one station per turn
2. **Collision Avoidance**: 
   - No two trains can occupy the same station (except start/end)
   - No two trains can use the same edge in the same turn
3. **Priority System**: Trains further along their path get movement priority

### Train Assignment
- Distributes trains across available paths
- Considers path length and existing train count
- Ensures optimal load balancing

## Error Handling

The program provides detailed error messages for:
- Invalid command line arguments
- Malformed map files
- Missing or duplicate stations
- Invalid station names or coordinates
- Unreachable destinations
- File I/O errors

## Example Output

```
Paths found:
Path 1: start -> middle -> end
Path 2: start -> alternative -> end

Train movement:
Turn 1: T1-middle T2-alternative
Turn 2: T1-end T2-end
```

## Testing

The project includes comprehensive test files in the `testdata/` directory:

- `small.map` - Basic network for testing
- `London.map` - Complex real-world example
- `LondonNoPath.map` - Test disconnected networks
- `10000.map` - Large network stress test
- Various error condition tests

## Project Structure

```
stations-pathfinder/
├── main.go              # Main program entry point
├── main_test.go         # Test suite
├── go.mod              # Go module definition
├── stations.txt        # Sample station names for generation
├── pathfinder/         # Core algorithm package
│   ├── parseMapFile.go # Map file parser
│   ├── findPath.go     # Pathfinding algorithms
│   ├── pipeline.go     # Train assignment logic
│   ├── simulate.go     # Movement simulation
│   └── generator.go    # Map file generation
└── testdata/           # Test map files
    ├── small.map
    ├── London.map
    └── ...
```

