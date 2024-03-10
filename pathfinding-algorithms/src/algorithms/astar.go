package algorithms

import (
	"container/heap"
	"errors"
	"math"
	"pathfinding-algorithms/datastructures"
	"pathfinding-algorithms/models"
	"sync"
)

type AStar struct {
	grid      *models.Grid
	solved    bool
	openSet   *datastructures.PriorityQueue
	snapshots *datastructures.Queue
	path      map[models.Node]models.Node
	closedSet map[*models.Node]bool
	fScore    map[*models.Node]float64
	gScore    map[*models.Node]float64
	mu        sync.Mutex
}

func (a *AStar) Init(width, height int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	var err error
	a.grid, err = models.NewGrid(width, height)
	if err != nil {
		return err
	}
	a.resetDataStructures()
	return nil
}

func (a *AStar) Clear() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return errors.New("grid is nil")
	}
	a.grid, _ = models.NewGrid(a.grid.GetWidth(), a.grid.GetHeight())
	a.resetDataStructures()
	return nil
}

func (a *AStar) resetDataStructures() {
	a.openSet = &datastructures.PriorityQueue{}
	a.openSet.Init()
	a.closedSet = make(map[*models.Node]bool)
	a.fScore = make(map[*models.Node]float64)
	a.gScore = make(map[*models.Node]float64)
	a.snapshots = &datastructures.Queue{}
	a.solved = false
	a.path = make(map[models.Node]models.Node)
}

// FindPath implements the pathfinding algorithm
func (a *AStar) FindPath() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return errors.New("grid is nil")
	}
	if a.solved {
		return errors.New("grid is solved")
	}
	// Pseudocode:
	// 1. Initialize open and closed lists
	// 2. Add the start node to the open list
	// 3. Loop until the open list is empty or the end node is reached
	//    a. Find the node with the lowest f score in the open list (current node)
	//    b. Remove current node from the open list and add it to the closed list
	//    c. For each neighbor of the current node:
	//       i. If neighbor is not traversable or in the closed list, skip to the next neighbor
	//       ii. If new path to neighbor is shorter OR neighbor is not in open list:
	//           - Set f, g, and h scores of the neighbor
	//           - Set parent of neighbor to current
	//           - If neighbor not in open list, add it
	// 4. Once the end node is reached, backtrack from the end node to start node to get the path
	startNode, err := a.grid.GetStart()
	if err != nil {
		return err
	}
	endNode, err := a.grid.GetEnd()
	if err != nil {
		return err
	}
	heap.Push(a.openSet, datastructures.NewItem(startNode, 0))
	a.gScore[startNode] = 0
	a.fScore[startNode] = heuristic(startNode, endNode)

	for a.openSet.Len() > 0 {
		current := heap.Pop(a.openSet).(*datastructures.Item).GetNode()
		if current == nil {
			return errors.New("current is nil")
		}
		if current.IsEnd {
			a.solved = true
			return nil
		}

		a.closedSet[current] = true
		if !current.IsStart && !current.IsEnd {
			current.Visited = true
		}

		neighbors, err := a.grid.GetNeighbors(current)
		if err != nil {
			return err
		}
		for _, neighbor := range neighbors {
			if a.closedSet[neighbor] || neighbor.IsWall {
				continue
			}

			tentativeGScore := a.gScore[current] + distBetween(current, neighbor)
			//if tentativeGScore >= a.gScore[neighbor] && a.openSet.Contains(neighbor) {
			//	continue
			//}
			if tentativeGScore >= a.gScore[neighbor] && a.gScore[neighbor] != 0 {
				continue
			}

			// This path is the best until now. Record it!
			a.path[*neighbor] = *current
			a.gScore[neighbor] = tentativeGScore
			a.fScore[neighbor] = a.gScore[neighbor] + heuristic(neighbor, endNode)
			if !a.openSet.Contains(neighbor) {
				heap.Push(a.openSet, datastructures.NewItem(neighbor, a.fScore[neighbor]))
			} else {
				a.openSet.Update(neighbor, a.fScore[neighbor])
			}
		}
		snapshot, err := a.grid.DeepCopy()
		if err != nil {
			return err
		}
		if nodes, err := snapshot.GetNodes(); err != nil {
			return err
		} else {
			a.snapshots.Enqueue(nodes)
		}
	}
	a.solved = true
	return errors.New("solution not found")

}

func (a *AStar) GetGrid() (*models.Grid, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return nil, errors.New("grid is nil")
	}
	return a.grid, nil
}

func (a *AStar) GetSnapshot() ([][]*models.Node, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return nil, errors.New("grid is nil")
	}
	if !a.solved {
		return nil, errors.New("astar is not solved")
	}
	if a.snapshots.IsEmpty() {
		return nil, nil
	} else {
		return a.snapshots.Dequeue().([][]*models.Node), nil
	}
}

func (a *AStar) GetPath() (map[models.Node]models.Node, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return nil, errors.New("grid is nil")
	}
	if !a.solved {
		return nil, errors.New("grid is not solved")
	}
	return a.path, nil
}

func (a *AStar) SetStart(x, y int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return errors.New("grid is nil")
	}
	if a.solved {
		return errors.New("grid is solved")
	}
	return a.grid.SetStart(x, y)
}

func (a *AStar) SetEnd(x, y int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return errors.New("grid is nil")
	}
	if a.solved {
		return errors.New("grid is solved")
	}
	return a.grid.SetEnd(x, y)
}

func (a *AStar) GetStart() (*models.Node, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return nil, errors.New("grid is nil")
	}
	return a.grid.GetStart()
}

func (a *AStar) GetEnd() (*models.Node, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return nil, errors.New("grid is nil")
	}
	return a.grid.GetEnd()
}

func (a *AStar) SetWall(x, y int, isWall bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.grid == nil {
		return errors.New("grid is nil")
	}
	if a.solved {
		return errors.New("grid is solved")
	}
	return a.grid.SetWall(x, y, isWall)
}

func heuristic(a, b *models.Node) float64 {
	dx := math.Abs(float64(a.X - b.X))
	dy := math.Abs(float64(a.Y - b.Y))
	return math.Sqrt(dx*dx + dy*dy)
}

func distBetween(current, neighbor *models.Node) float64 {
	dx := math.Abs(float64(current.X - neighbor.X))
	dy := math.Abs(float64(current.Y - neighbor.Y))
	if dx == 1 && dy == 1 {
		// Diagonal movement
		return 1.414 // sqrt(2), assuming diagonal cost is sqrt(2) times an orthogonal step
	}
	// Orthogonal movement
	return 1
}
