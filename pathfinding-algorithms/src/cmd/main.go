//go:build js && wasm

package main

import (
	"fmt"
	"pathfinding-algorithms/models"
	"pathfinding-algorithms/pathfinder"
	"syscall/js"
)

var pf *pathfinder.Pathfinder
var snapshotPointer *[]uint8

func main() {
	c := make(chan struct{}, 0)
	pf = pathfinder.NewPathfinder()
	<-c
}

//export setActiveAlgorithm
func setActiveAlgorithm(algorithm pathfinder.Algorithm, width, height int) bool {
	err := pf.SetActiveAlgorithm(algorithm, width, height)
	if err != nil {
		log(fmt.Sprintf("Error setting active algorithm %v with dimensions (%v, %v): %v", algorithm, width, height, err))
		return false
	}
	return true
}

//export setStart
func setStart(x, y int) bool {
	err := pf.SetStart(x, y)
	if err != nil {
		log(fmt.Sprintf("Error setting start (%v, %v): %v", x, y, err))
		return false
	}
	return true
}

//export setEnd
func setEnd(x, y int) bool {
	err := pf.SetEnd(x, y)
	if err != nil {
		log(fmt.Sprintf("Error setting end (%v,%v): %v", x, y, err))
		return false
	}
	return true
}

//export getStart
func getStart() uint64 {
	node, err := pf.GetStart()
	if err != nil {
		log(fmt.Sprintf("Error getting start: %v", err))
		return 0
	}
	return convertNodeLocationToUint64(node.X, node.Y)
}

//export getEnd
func getEnd() uint64 {
	node, err := pf.GetEnd()
	if err != nil {
		log(fmt.Sprintf("Error getting end: %v", err))
		return 0
	}
	return convertNodeLocationToUint64(node.X, node.Y)
}

//export setWall
func setWall(x, y int, isWall bool) bool {
	err := pf.SetWall(x, y, isWall)
	if err != nil {
		log(fmt.Sprintf("Error setting wall (%v,%v) to %v: %v", x, y, isWall, err))
		return false
	}
	return true
}

//export clearGrid
func clearGrid() bool {
	err := pf.ClearGrid()
	if err != nil {
		log(fmt.Sprintf("Error clearing grid: %v", err))
		return false
	}
	return true
}

//export changeGridSize
func changeGridSize(width, height int) bool {
	err := pf.ChangeGridSize(width, height)
	if err != nil {
		log(fmt.Sprintf("Error changing grid size to (%d,%d): %v", width, height, err))
		return false
	}
	return true
}

//export getGrid
func getGrid() *[]uint8 {
	grid, err := pf.GetNodes()
	if err != nil {
		log(fmt.Sprintf("Error getting grid: %v", err))
		return nil
	}
	out := make([]uint8, len(grid)*len(grid[0]))

	encodeGrid(grid, out)
	return &out
}

func encodeGrid(grid [][]*models.Node, encodedGrid []uint8) {
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[0]); col++ {
			encodedNode := convertNodeToUint(*grid[row][col])
			encodedGrid[row*len(grid[0])+col] = encodedNode
		}
	}
}

func convertNodeToUint(node models.Node) uint8 {

	// Packing boolean values into a single uint32
	var boolPack uint8 = 0
	if node.Visited {
		boolPack |= 1 << 0 // Use bit 0 for Visited
	}
	if node.IsWall {
		boolPack |= 1 << 1 // Use bit 1 for IsWall
	}
	if node.IsStart {
		boolPack |= 1 << 2 // Use bit 2 for IsStart
	}
	if node.IsEnd {
		boolPack |= 1 << 3 // Use bit 3 for IsEnd
	}

	return boolPack
}

func convertNodeLocationToUint64(x, y int) uint64 {
	return uint64(x)<<32 | uint64(y)
}

//export findPath
func findPath() bool {
	err := pf.FindPath()
	if err != nil {
		log(fmt.Sprintf("Error finding path: %v", err))
		return false
	}
	return true
}

//export getNumNodes
func getNumNodes() int {
	grid, err := pf.GetNodes()
	if err != nil {
		log(fmt.Sprintf("Error getting grid: %v", err))
		return -1
	}
	return len(grid) * len(grid[0])
}

//export getNumPathNodes
func getNumPathNodes() int {
	path, err := pf.GetPath()
	if err != nil {
		log(fmt.Sprintf("Error getting path: %v", err))
		return -1
	}
	return len(path)
}

//export getWidth
func getWidth() int {
	grid, err := pf.GetNodes()
	if err != nil {
		log(fmt.Sprintf("Error getting grid: %v", err))
		return -1
	}
	return len(grid[0])
}

//export getSnapshot
func getSnapshot() *[]uint8 {
	snapshot, err := pf.GetSnapshot()
	if err != nil || snapshot == nil {
		return nil
	}
	if snapshotPointer == nil || len(*snapshotPointer) != len(snapshot)*len(snapshot[0]) {
		ptr := make([]uint8, len(snapshot)*len(snapshot[0]))
		snapshotPointer = &ptr
	}
	encodeGrid(snapshot, *snapshotPointer)
	return snapshotPointer
}

//export getPath
func getPath() *[]uint32 {
	path, err := pf.GetPath()
	if err != nil {
		log(fmt.Sprintf("Error getting path: %v", err))
		return nil
	}
	out := make([]uint32, len(path)*4)
	encodePath(path, out)
	return &out
}

func encodePath(path map[models.Node]models.Node, out []uint32) {
	i := 0
	for prevNode, curNode := range path {
		// Example encoding: node.X, node.Y -> nextNode.X, nextNode.Y
		// Adjust based on your actual node structure and encoding needs
		if i < len(path) {
			index := i * 4
			out[index] = uint32(prevNode.X)
			out[index+1] = uint32(prevNode.Y)
			out[index+2] = uint32(curNode.X)
			out[index+3] = uint32(curNode.Y)
			i++
		}
	}
}

//export generateMaze
func generateMaze() bool {
	if err := pf.GenerateMaze(); err != nil {
		log(fmt.Sprintf("Error generating Maze: %v", err))
		return false
	}
	return true
}

func log(a string) {
	js.Global().Get("console").Get("log").Invoke(a)
}
