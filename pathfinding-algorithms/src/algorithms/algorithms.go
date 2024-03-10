package algorithms

import "pathfinding-algorithms/models"

type PathfindingAlgorithm interface {
	Init(width, height int) error
	Clear() error
	FindPath() error
	GetGrid() (*models.Grid, error)
	GetSnapshot() ([][]*models.Node, error)
	GetPath() (map[models.Node]models.Node, error)
	SetStart(x, y int) error
	SetEnd(x, y int) error
	SetWall(x, y int, visited bool) error
	GetStart() (node *models.Node, err error)
	GetEnd() (node *models.Node, err error)
}
