package algorithms

import (
	"container/heap"
	"errors"
	"pathfinding-algorithms/datastructures"
	"pathfinding-algorithms/models"
	"sync"
)

// Dijkstra struct holds the necessary components for the pathfinding algorithm
type Dijkstra struct {
	grid      *models.Grid
	solved    bool
	openSet   *datastructures.PriorityQueue
	snapshots *datastructures.Queue
	path      map[models.Node]models.Node
	distances map[*models.Node]float64
	mu        sync.Mutex
}

// Init initializes the Dijkstra with a grid and necessary data structures
func (d *Dijkstra) Init(width, height int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	var err error
	d.grid, err = models.NewGrid(width, height)
	if err != nil {
		return err
	}
	d.resetDataStructures()
	return nil
}

// Clear resets the Dijkstra for a new pathfinding operation
func (d *Dijkstra) Clear() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	var err error
	d.grid, err = models.NewGrid(d.grid.GetWidth(), d.grid.GetHeight())
	if err != nil {
		return err
	}
	d.resetDataStructures()
	return nil
}

func (d *Dijkstra) resetDataStructures() {
	d.solved = false
	d.openSet = &datastructures.PriorityQueue{}
	d.openSet.Init()
	d.distances = make(map[*models.Node]float64)
	d.snapshots = &datastructures.Queue{}
	d.path = make(map[models.Node]models.Node)
}

// FindPath implements Dijkstra's algorithm for finding the shortest path
func (d *Dijkstra) FindPath() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.grid == nil {
		return errors.New("grid is nil")
	}
	if d.solved {
		return errors.New("grid is solved")
	}

	startNode, err := d.grid.GetStart()
	if err != nil {
		return err
	}
	endNode, err := d.grid.GetEnd()
	if err != nil {
		return err
	}
	d.distances[startNode] = 0
	heap.Push(d.openSet, datastructures.NewItem(startNode, 0))

	for d.openSet.Len() > 0 {
		current := heap.Pop(d.openSet).(*datastructures.Item).GetNode()

		if current == endNode {
			d.solved = true
			return nil
		}

		if !current.IsStart && !current.IsEnd {
			current.Visited = true
		}

		neighbors, err := d.grid.GetNeighbors(current)
		if err != nil {
			return err
		}
		for _, neighbor := range neighbors {
			if neighbor.Visited || neighbor.IsWall {
				continue
			}

			tentativeDistance := d.distances[current] + distBetween(current, neighbor)
			if dist, exists := d.distances[neighbor]; !exists || tentativeDistance < dist {
				d.distances[neighbor] = tentativeDistance
				d.path[*neighbor] = *current
				if !d.openSet.Contains(neighbor) {
					heap.Push(d.openSet, datastructures.NewItem(neighbor, tentativeDistance))
				} else {
					d.openSet.Update(neighbor, tentativeDistance)
				}
			}
		}

		snapshot, err := d.grid.DeepCopy()
		if err != nil {
			return err
		}
		if nodes, err := snapshot.GetNodes(); err != nil {
			return err
		} else {
			d.snapshots.Enqueue(nodes)
		}
	}
	d.solved = true
	return errors.New("path to destination not found")
}

func (d *Dijkstra) GetNodes() ([][]*models.Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return nil, errors.New("grid is nil")
	}
	return d.grid.GetNodes()
}

func (d *Dijkstra) GetSnapshot() ([][]*models.Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return nil, errors.New("grid is nil")
	}
	if !d.solved {
		return nil, errors.New("astar is not solved")
	}
	if d.snapshots.IsEmpty() {
		return nil, nil
	} else {
		return d.snapshots.Dequeue().([][]*models.Node), nil
	}
}

func (d *Dijkstra) GetPath() (map[models.Node]models.Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return nil, errors.New("grid is nil")
	}
	if !d.solved {
		return nil, errors.New("grid is not solved")
	}
	return d.path, nil
}

func (d *Dijkstra) SetStart(x, y int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return errors.New("grid is nil")
	}
	if d.solved {
		return errors.New("grid is solved")
	}
	return d.grid.SetStart(x, y)
}

func (d *Dijkstra) SetEnd(x, y int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return errors.New("grid is nil")
	}
	if d.solved {
		return errors.New("grid is solved")
	}
	return d.grid.SetEnd(x, y)
}

func (d *Dijkstra) GetStart() (*models.Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return nil, errors.New("grid is nil")
	}
	return d.grid.GetStart()
}

func (d *Dijkstra) GetEnd() (*models.Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return nil, errors.New("grid is nil")
	}
	return d.grid.GetEnd()
}

func (d *Dijkstra) SetWall(x, y int, isWall bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.grid == nil {
		return errors.New("grid is nil")
	}
	if d.solved {
		return errors.New("grid is solved")
	}
	return d.grid.SetWall(x, y, isWall)
}
