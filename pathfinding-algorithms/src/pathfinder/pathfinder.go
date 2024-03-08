package pathfinder

import (
	"errors"
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

func (p *Pathfinder) GetGrid() ([][]*models.Node, error) {
	if p.activeAlgorithm == nil {
		return nil, errors.New("no active algorithm set")
	}
	return p.activeAlgorithm.GetNodes()
}

func (p *Pathfinder) GetDimensions() (width int, height int, err error) {
	if p.activeAlgorithm == nil {
		return 0, 0, errors.New("no active algorithm set")
	}
	nodes, err := p.activeAlgorithm.GetNodes()
	if err != nil {
		return 0, 0, err
	}
	return len(nodes[0]), len(nodes), nil
}

func (p *Pathfinder) GetSnapshot() ([][]*models.Node, error) {
	if p.activeAlgorithm == nil {
		return nil, errors.New("no active algorithm set")
	}
	if grid, err := p.activeAlgorithm.GetSnapshot(); err != nil {
		return nil, err
	} else {
		return grid, nil
	}
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
