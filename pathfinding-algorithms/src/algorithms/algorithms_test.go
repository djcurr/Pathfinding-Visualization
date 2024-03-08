package algorithms_test

import (
	"fmt"
	"pathfinding-algorithms/algorithms"

	"testing"
)

func TestFindPathSimple(t *testing.T) {
	aStar := algorithms.AStar{}
	aStar.Init(20, 20)
	aStar.SetStart(1, 1)
	aStar.SetEnd(18, 18)

	path, _, err := aStar.FindPath()
	if err != nil {
		t.Fatalf("FindPath failed: %v", err)
	}

	if len(path) == 0 {
		t.Fatalf("No path found when one was expected")
	}
	fmt.Printf("Path found: %v\n", path)
}
