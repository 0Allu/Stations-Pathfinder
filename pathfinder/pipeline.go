package pathfinder

import (
	"fmt"
)

type Train struct {
	Name   string
	Path   []string
	Index  int
	Active bool
}

// Train Assignment
func AssignToPipelines(paths [][]string, numTrains int) []*Train {
	numPipelines := len(paths)

	pipelineLengths := make([]int, numPipelines)
	for i, path := range paths {
		pipelineLengths[i] = len(path) - 1
	}

	trainsPerPipeline := make([]int, numPipelines)
	trainAssignments := make([]int, numTrains)

	for i := 0; i < numTrains && i < numPipelines; i++ {
		trainAssignments[i] = i
		trainsPerPipeline[i] = 1
	}

	for i := numPipelines; i < numTrains; i++ {
		bestPipeline := 0
		bestScore := trainsPerPipeline[0] + pipelineLengths[0]

		for j := 1; j < numPipelines; j++ {
			score := trainsPerPipeline[j] + pipelineLengths[j]

			if score < bestScore || (score == bestScore && pipelineLengths[j] < pipelineLengths[bestPipeline]) {
				bestPipeline = j
				bestScore = score
			}
		}
		trainAssignments[i] = bestPipeline
		trainsPerPipeline[bestPipeline]++
	}

	trains := make([]*Train, numTrains)
	for i, pipelineIndex := range trainAssignments {
		trains[i] = &Train{
			Name:   fmt.Sprintf("T%d", i+1),
			Path:   paths[pipelineIndex],
			Index:  0,
			Active: true,
		}
	}
	return trains
}