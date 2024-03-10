package models

import (
	"errors"
)

type Node struct {
	X, Y    int
	Visited bool
	IsWall  bool
	IsStart bool
	IsEnd   bool
}

type Grid struct {
	width, height    int // Dimensions of the grid
	nodes            [][]*Node
	start, end       *Node
	diagonalMovement bool
}

type Point struct{ Dx, Dy int }

var diagonalDirections = []Point{
	{-1, -1}, {0, -1}, {1, -1}, // Above
	{-1, 0}, {1, 0}, // Sides
	{-1, 1}, {0, 1}, {1, 1}, // Below
}

var cartesianDirections = []Point{
	{0, -1},
	{-1, 0}, {1, 0},
	{0, 1},
}

// NewGrid creates a new grid of the given width and height
func NewGrid(width, height int) (*Grid, error) {
	if (width < 10) || (height < 10) {
		return nil, errors.New("width and height must be greater than 10")
	}

	nodes := make([][]*Node, height)
	yMid := height / 2
	xMid1 := width * 1 / 4
	xMid2 := width * 3 / 4
	for i := range nodes {
		nodes[i] = make([]*Node, width)
		for j := range nodes[i] {
			nodes[i][j] = &Node{
				X:       j,
				Y:       i,
				IsWall:  i == 0 || j == 0 || i == height-1 || j == width-1,
				Visited: false,
				IsStart: i == yMid && j == xMid1,
				IsEnd:   i == yMid && j == xMid2+1,
			}
		}
	}

	return &Grid{
		width:            width,
		height:           height,
		nodes:            nodes,
		start:            nodes[yMid][xMid1],
		end:              nodes[yMid][xMid2+1],
		diagonalMovement: false,
	}, nil
}

// SetStart sets the start node
func (g *Grid) SetStart(x, y int) error {
	if g == nil {
		return errors.New("grid is nil")
	}
	if g.nodes[y][x].IsEnd || g.nodes[y][x].IsWall {
		return errors.New("invalid location")
	}
	g.nodes[g.start.Y][g.start.X].IsStart = false
	if y >= 0 && y < g.height && x >= 0 && x < g.width {
		g.nodes[y][x].IsStart = true
		g.start = g.nodes[y][x]
		return nil
	}
	return errors.New("invalid location")
}

// SetEnd sets the end node
func (g *Grid) SetEnd(x, y int) error {
	if g == nil {
		return errors.New("grid is nil")
	}
	if g.nodes[y][x].IsStart || g.nodes[y][x].IsWall {
		return errors.New("invalid location")
	}
	g.nodes[g.end.Y][g.end.X].IsEnd = false
	if y >= 0 && y < g.height && x >= 0 && x < g.width {
		g.nodes[y][x].IsEnd = true
		g.end = g.nodes[y][x]
		return nil
	}
	return errors.New("invalid location")
}

func (g *Grid) GetWidth() int {
	if g == nil {
		return 0
	}
	return g.width
}

func (g *Grid) GetHeight() int {
	if g == nil {
		return 0
	}
	return g.height
}

func (g *Grid) GetNode(x, y int) (*Node, error) {
	if g == nil {
		return nil, errors.New("grid is nil")
	}
	if x >= 0 && x < g.width && y >= 0 && y < g.height {
		return g.nodes[y][x], nil
	}
	return nil, errors.New("invalid location")
}

func (g *Grid) GetStart() (*Node, error) {
	if g == nil {
		return nil, errors.New("grid is nil")
	}
	return g.start, nil
}

func (g *Grid) GetEnd() (*Node, error) {
	if g == nil {
		return nil, errors.New("grid is nil")
	}
	return g.end, nil
}

func (g *Grid) SetWall(x, y int, isWall bool) error {
	if g == nil {
		return errors.New("grid is nil")
	}
	if g.nodes[y][x].IsStart || g.nodes[y][x].IsEnd {
		return errors.New("invalid location")
	}
	if y >= 0 && y < g.height && x >= 0 && x < g.width {
		g.nodes[y][x].IsWall = isWall
		return nil
	}
	return errors.New("invalid location")
}

func (g *Grid) GetNeighbors(node *Node) ([]*Node, error) {
	var neighbors []*Node
	directions, err := g.GetDirections()
	if err != nil {
		return nil, err
	}
	for _, d := range directions {
		nx, ny := node.X+d.Dx, node.Y+d.Dy
		if nx >= 0 && nx < g.width && ny >= 0 && ny < g.height {
			neighbors = append(neighbors, g.nodes[ny][nx])
		}
	}

	return neighbors, nil
}

func (g *Grid) GetNodes() ([][]*Node, error) {
	if g == nil {
		return nil, errors.New("grid is nil")
	}
	return g.nodes, nil
}

func (g *Grid) GetDirections() ([]Point, error) {
	if g == nil {
		return nil, errors.New("grid is nil")
	}
	if g.diagonalMovement {
		return diagonalDirections, nil
	} else {
		return cartesianDirections, nil
	}
}

func (g *Grid) DeepCopy() (*Grid, error) {
	if g == nil {
		return nil, errors.New("grid is nil")
	}
	newGrid, _ := NewGrid(g.width, g.height)
	for y, row := range g.nodes {
		for x, node := range row {
			newGrid.nodes[y][x] = &Node{
				X:       node.X,
				Y:       node.Y,
				IsWall:  node.IsWall,
				Visited: node.Visited,
				IsStart: node.IsStart,
				IsEnd:   node.IsEnd,
			}
		}
	}
	return newGrid, nil
}
