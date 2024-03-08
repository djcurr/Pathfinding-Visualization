//go:build js && wasm

package main

import (
	"fmt"
	"pathfinding-algorithms/models"
	"pathfinding-algorithms/pathfinder"
	"reflect"
	"syscall/js"
	"unsafe"
)

var pf *pathfinder.Pathfinder

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
func getGrid(out *uint8, len int) bool {
	grid, err := pf.GetGrid()
	if err != nil {
		log(fmt.Sprintf("Error getting grid: %v", err))
		return false
	}

	if width, height, err := pf.GetDimensions(); err != nil {
		log(fmt.Sprintf("Error getting dimensions: %v", err))
		return false
	} else if width*height != len {
		log(fmt.Sprintf("Grid array is not large enough: %v", len))
		return false
	}

	encodeGrid(grid, out)
	return true
}

func encodeGrid(grid [][]*models.Node, encodedGrid *uint8) {
	width, height, err := pf.GetDimensions()
	if err != nil {
		log(fmt.Sprintf("Error getting dimensions: %v", err))
	}
	// Create a slice header
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&encodedGrid))
	nodes := width * height
	sliceHeader.Len = uintptr(unsafe.Pointer(&nodes))
	sliceHeader.Cap = uintptr(unsafe.Pointer(&nodes))

	// Now, firstElemPtr is a slice that you can index into
	dataSlice := *(*[]uint8)(unsafe.Pointer(sliceHeader))
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			encodedNode := convertNodeToUint(*grid[row][col])
			dataSlice[row*width+col] = encodedNode
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
	grid, err := pf.GetGrid()
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
	grid, err := pf.GetGrid()
	if err != nil {
		log(fmt.Sprintf("Error getting grid: %v", err))
		return -1
	}
	return len(grid[0])
}

//export getSnapshot
func getSnapshot(out *uint8, length int) bool {
	snapshot, err := pf.GetSnapshot()
	if err != nil {
		log(fmt.Sprintf("Error getting snapshot: %v", err))
		return false
	}

	if snapshot == nil {
		return false
	}

	if width, height, err := pf.GetDimensions(); err != nil {
		log(fmt.Sprintf("Error getting dimensions: %v", err))
		return false
	} else if width*height != length {
		log(fmt.Sprintf("Grid array is not large enough: %v", length))
		return false
	}

	encodeGrid(snapshot, out)
	return true
}

//export getPath
func getPath(out *uint32, length int) bool {
	path, err := pf.GetPath()
	if err != nil {
		log(fmt.Sprintf("Error getting snapshot: %v", err))
		return false
	}

	if len(path)*4 > length {
		log("Path array is not large enough")
		return false
	}
	encodePath(path, out)
	return true
}

func encodePath(path map[models.Node]models.Node, out *uint32) {
	nodes, err := pf.GetPath()
	if err != nil {
		log(fmt.Sprintf("Error getting path: %v", err))
	}
	// Create a slice header to manipulate the output array
	length := len(nodes)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&out))
	sliceHeader.Len = uintptr(unsafe.Pointer(&length))
	sliceHeader.Cap = uintptr(unsafe.Pointer(&length))

	// Convert the *uint8 to a []uint8 slice for easier manipulation
	dataSlice := *(*[]uint32)(unsafe.Pointer(sliceHeader))
	i := 0
	for prevNode, curNode := range path {
		// Example encoding: node.X, node.Y -> nextNode.X, nextNode.Y
		// Adjust based on your actual node structure and encoding needs
		if i+1 < len(path) {
			index := i * 4
			dataSlice[index] = uint32(prevNode.X)
			dataSlice[index+1] = uint32(prevNode.Y)
			dataSlice[index+2] = uint32(curNode.X)
			dataSlice[index+3] = uint32(curNode.Y)
			i++
		}
	}
}

func log(a string) {
	js.Global().Get("console").Get("log").Invoke(a)
}
