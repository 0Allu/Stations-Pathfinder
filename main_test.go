package main

import (
	"os/exec"
	"strings"
	"testing"
)

// It's best to create each function separately, without using the structure.
// ### go test -v ###
var testData = []struct {
	command  []string
	message  string
	errorMsg string
}{
	// ### TestTooFewArgs ###
	{[]string{"go", "run", "main.go", "testdata/London.map", "waterloo"},
		"Too few command line arguments are used\ngo run main.go testdata/London.map waterloo",
		"expected error for invalid number of trains, got"},

	// ### TestTooManyArgs ###
	{[]string{"go", "run", "main.go", "testdata/London.map", "waterloo", "st_pancras", "1", "extra"},
		"Too many command line arguments are used\ngo run main.go testdata/London.map waterloo st_pancras 1 extra",
		"expected error for invalid number of trains, got"},

	// ### TestStartStationNotExist ###
	{[]string{"go", "run", "main.go", "testdata/London.map", "not_exist", "st_pancras", "1"},
		"The start station does not exist\ngo run main.go testdata/London.map not_exist st_pancras 1",
		"expected error for non-existent start station, got"},

	// ### TestEndStationNotExist ###
	{[]string{"go", "run", "main.go", "testdata/London.map", "waterloo", "not_exist", "1"},
		"The end station does not exist\ngo run main.go testdata/London.map waterloo not_exist 1",
		"expected error for non-existent end station, got"},

	// ### TestSameStartEnd ###
	{[]string{"go", "run", "main.go", "testdata/London.map", "waterloo", "waterloo", "1"},
		"The start and end station are the same\ngo run main.go testdata/London.map waterloo waterloo 1",
		"expected error for same start and end station, got"},

	// ### TestInvalidTrains ###
	{[]string{"go", "run", "main.go", "testdata/London.map", "waterloo", "st_pancras", "-1"},
		"The number of trains is not a valid positive integer\ngo run main.go testdata/London.map waterloo st_pancras -1",
		"expected error for invalid number of trains, got"},

	// ### TestZeroTrains ###
	{[]string{"go", "run", "main.go", "testdata/London.map", "waterloo", "st_pancras", "0"},
		"The number of trains is 0\ngo run main.go testdata/London.map waterloo st_pancras 0",
		"expected error for invalid number of trains, got"},

	// ### TestNoPath ###
	{[]string{"go", "run", "main.go", "testdata/LondonNoPath.map", "waterloo", "st_pancras", "2"},
		"No path exists between the start and end stations\ngo run main.go testdata/LondonNoPath.map waterloo st_pancras 2",
		"expected error for no path between stations, got"},

	// ### TestDuplicateRoutes ###
	{[]string{"go", "run", "main.go", "testdata/LondonDuplicateRoutes.map", "waterloo", "st_pancras", "2"},
		"Duplicate routes exist between two stations\ngo run main.go testdata/LondonDuplicateRoutes.map waterloo st_pancras 2",
		"expected error for duplicate routes, got"},

	// ### TestNotPositiveIntengerCoordinate ###
	{[]string{"go", "run", "main.go", "testdata/LondonNotPositiveIntengerCoordinate.map", "waterloo", "st_pancras", "2"},
		"The coordinates are not valid positive integers\ngo run main.go testdata/LondonNotPositiveIntengerCoordinate.map waterloo st_pancras 2",
		"expected error for non-positive integer coordinates, got"},

	// ### TestSameCoordinates ###
	{[]string{"go", "run", "main.go", "testdata/LondonSameCoordinates.map", "waterloo", "st_pancras", "2"},
		"Two stations exist at the same coordinates\ngo run main.go testdata/LondonSameCoordinates.map waterloo st_pancras 2",
		"expected error for same coordinates, got"},

	// ### TestConnectionWithNotExistStation ###
	{[]string{"go", "run", "main.go", "testdata/LondonConnectionWithNotExistStation.map", "waterloo", "st_pancras", "2"},
		"A connection is made with a station which does not exist\ngo run main.go testdata/LondonConnectionWithNotExistStation.map waterloo st_pancras 2",
		"expected error for connection with non-existent station, got"},

	// ### TestNameDuplicated ###
	{[]string{"go", "run", "main.go", "testdata/LondonNameDuplicated.map", "waterloo", "st_pancras", "2"},
		"Station names are duplicated\ngo run main.go testdata/LondonNameDuplicated.map waterloo st_pancras 2",
		"expected error for duplicated station names, got"},

	// ### TestStationNamesInvalid ###
	{[]string{"go", "run", "main.go", "testdata/LondonStationNamesInvalid.map", "waterloo", "st_pancras", "2"},
		"Station names are invalid\ngo run main.go testdata/LondonStationNamesInvalid.map waterloo st_pancras 2",
		"expected error for invalid station names, got"},

	// ### TestWithoutStations ###
	{[]string{"go", "run", "main.go", "testdata/LondonWithoutStations.map", "waterloo", "st_pancras", "2"},
		"the map does not contain a \"stations:\" section\ngo run main.go testdata/LondonWithoutStations.map waterloo st_pancras 2",
		"expected error for no stations in the map, got"},

	// ### TestWithoutConnections ###
	{[]string{"go", "run", "main.go", "testdata/LondonWithoutConnections.map", "waterloo", "st_pancras", "2"},
		"The map does not contain a \"connections:\" section\ngo run main.go testdata/LondonWithoutConnections.map waterloo st_pancras 2",
		"expected error for no connections in the map, got"},

	// ### TestMapWith10001Stations ###
	{[]string{"go", "run", "main.go", "testdata/10001.map", "waterloo", "st_pancras", "2"},
		"The map file contains more than 10000 stations\ngo run main.go testdata/10001.map waterloo st_pancras 2",
		"expected error for more than 10000 stations, got"},
}

var count int // Bad idea if you want to run tests in parallel t.Parallel()

func errorTest(t *testing.T, command []string, message string, errorMsg string) {
	cmd := exec.Command(command[0], command[1:]...)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "Error") {
		t.Errorf(red+"%s: %s"+reset, errorMsg, out)
	}
	count++
	t.Logf(green+"Test %d: %s"+reset, count, message)
}

func TestTooFewArgs(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestTooManyArgs(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestStartStationNotExist(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestEndStationNotExist(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestSameStartEnd(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestInvalidTrains(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestZeroTrains(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestNoPath(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestDuplicateRoutes(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestNotPositiveIntengerCoordinate(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestSameCoordinates(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestConnectionWithNotExistStation(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestNameDuplicated(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestStationNamesInvalid(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestWithoutStations(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestWithoutConnections(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}

func TestMapWith10001Stations(t *testing.T) {
	errorTest(t, testData[count].command, testData[count].message, testData[count].errorMsg)
}
