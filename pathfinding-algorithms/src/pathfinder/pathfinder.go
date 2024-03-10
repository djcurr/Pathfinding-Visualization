package pathfinder

import (
	"errors"
	"math/rand"
	"pathfinding-algorithms/algorithms"
	"pathfinding-algorithms/models"
)

type Algorithm int

const (
	aStar Algorithm = iota
	dijkstra
)

type Pathfinder struct {
	algorithmsMap   map[Algorithm]func() algorithms.PathfindingAlgorithm
	activeAlgorithm algorithms.PathfindingAlgorithm
}

// NewPathfinder creates a new Pathfinder instance
func NewPathfinder() *Pathfinder {
	return &Pathfinder{
		algorithmsMap: map[Algorithm]func() algorithms.PathfindingAlgorithm{
			aStar: func() algorithms.PathfindingAlgorithm {
				return &algorithms.AStar{}
			},
			dijkstra: func() algorithms.PathfindingAlgorithm {
				return &algorithms.Dijkstra{}
			},
		},
	}
}

// SetActiveAlgorithm sets the active algorithm based on the name, resets the grid everytime
func (p *Pathfinder) SetActiveAlgorithm(algorithm Algorithm, width, height int) error {
	algFunc, exists := p.algorithmsMap[algorithm]
	if !exists {
		return errors.New("algorithm not found")
	}

	p.activeAlgorithm = algFunc()
	return p.activeAlgorithm.Init(width, height)
}

func (p *Pathfinder) SetStart(x, y int) error {
	if p.activeAlgorithm == nil {
		return errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.SetStart(x, y)
}

func (p *Pathfinder) SetEnd(x, y int) error {
	if p.activeAlgorithm == nil {
		return errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.SetEnd(x, y)
}

func (p *Pathfinder) SetWall(x, y int, isWall bool) error {
	if p.activeAlgorithm == nil {
		return errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.SetWall(x, y, isWall)
}

func (p *Pathfinder) FindPath() error {
	if p.activeAlgorithm == nil {
		return errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.FindPath()
}

func (p *Pathfinder) ClearGrid() error {
	if p.activeAlgorithm == nil {
		return errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.Clear()
}

func (p *Pathfinder) ChangeGridSize(width, height int) error {
	if p.activeAlgorithm == nil {
		return errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.Init(width, height)
}

func (p *Pathfinder) GetNodes() ([][]*models.Node, error) {
	if p.activeAlgorithm == nil {
		return nil, errors.New("no active algorithm set")
	}
	grid, err := p.activeAlgorithm.GetGrid()
	if err != nil {
		return nil, err
	}
	return grid.GetNodes()
}

func (p *Pathfinder) GetDimensions() (width int, height int, err error) {
	if p.activeAlgorithm == nil {
		return 0, 0, errors.New("no active algorithm set")
	}
	grid, err := p.activeAlgorithm.GetGrid()
	if err != nil {
		return 0, 0, err
	}
	nodes, err := grid.GetNodes()
	if err != nil {
		return 0, 0, err
	}
	return len(nodes[0]), len(nodes), nil
}

func (p *Pathfinder) GetSnapshot() ([][]*models.Node, error) {
	if p.activeAlgorithm == nil {
		return nil, errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.GetSnapshot()
}

func (p *Pathfinder) GetPath() (map[models.Node]models.Node, error) {
	if p.activeAlgorithm == nil {
		return nil, errors.New("no active algorithm set")
	}
	if grid, err := p.activeAlgorithm.GetPath(); err != nil {
		return nil, err
	} else {
		return grid, nil
	}
}

func (p *Pathfinder) GetStart() (node *models.Node, err error) {
	if p.activeAlgorithm == nil {
		return nil, errors.New("no active algorithm set")
	}
	if node, err := p.activeAlgorithm.GetStart(); err != nil {
		return nil, err
	} else {
		return node, nil
	}
}

func (p *Pathfinder) GetEnd() (node *models.Node, err error) {
	if p.activeAlgorithm == nil {
		return nil, errors.New("no active algorithm set")
	}
	if node, err := p.activeAlgorithm.GetEnd(); err != nil {
		return nil, err
	} else {
		return node, nil
	}
}

func (p *Pathfinder) GenerateMaze() error {
	maze, err := p.GetNodes()
	if err != nil {
		return err
	}
	for _, row := range maze {
		for _, node := range row {
			if !node.IsStart && !node.IsEnd {
				node.IsWall = true
			}
		}
	}
	if err := p.generateMaze(maze, 2, 2); err != nil {
		return err
	}
	for _, row := range maze {
		for _, node := range row {
			node.Visited = false
		}
	}
	return nil
}

func (p *Pathfinder) generateMaze(maze [][]*models.Node, x, y int) error {
	directions := []models.Point{
		{0, -1},
		{-1, 0}, {1, 0},
		{0, 1},
	}
	// Mark the current cell as visited

	maze[y][x].IsWall = false
	maze[y][x].Visited = true

	// Randomly order the directions
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	// Explore the neighbors in a random order
	for _, d := range directions {
		nx, ny := x+2*d.Dx, y+2*d.Dy

		// Check bounds and if the neighbor has been visited
		if nx >= 0 && nx < len(maze[0]) && ny >= 0 && ny < len(maze) && !maze[ny][nx].Visited {
			maze[(y+ny)/2][(x+nx)/2].IsWall = false
			if err := p.generateMaze(maze, nx, ny); err != nil {
				return err
			} // Recursively visit the neighbor
		}
	}
	return nil
}
