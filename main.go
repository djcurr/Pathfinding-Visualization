//go:build js && wasm

package main

import (
	"fmt"
	"math"
	"strconv"
	"syscall/js"
	"time"

	"honnef.co/go/js/dom/v2"
)

type Node struct {
	location [2]int
	distance int
	visited  bool
	parent   []*Node
	wall     bool
}

type Nodes struct {
	startNode Node
	endNode   Node
	nodeSlice []Node
}

var nodes Nodes
var start bool
var wallsEnabled = false

func makeTable(rows int, columns int) string {

	tableStr := "<table  draggable=\"false\">"

	for i := 0; i < rows; i++ {

		tableStr += "<tr draggable=\"false\">"

		for j := 0; j < columns; j++ {

			elementLocation := strconv.Itoa(j) + "," + strconv.Itoa(i)
			id := "id=\"" + elementLocation + "\""

			if j == ((columns/8)-1) && i == (rows/2) {
				makeStartNode(j, i, elementLocation, id, &tableStr)
			} else if j == (columns-(columns/8)) && i == (rows/2) {
				makeEndNode(j, i, elementLocation, id, &tableStr)
			} else {
				makeNormalNode(j, i, elementLocation, id, &tableStr)
			}

		}

		tableStr += "</tr>"

	}

	tableStr += "</table>"

	return tableStr

}

func makeStartNode(j int, i int, elementLocation string, id string, tableStr *string) {

	redBg := "style=\"background-color:red\""

	loc := Node{
		location: [2]int{j, i},
		distance: 0,
		visited:  false,
	}
	nodes.startNode = loc
	nodes.nodeSlice = append(nodes.nodeSlice, loc)

	*tableStr += "<td " + redBg + " draggable=\"true\"" + id + ">" + "</td>"
}

func makeEndNode(j int, i int, elementLocation string, id string, tableStr *string) {

	blueBg := "style=\"background-color:blue\""

	loc := Node{
		location: [2]int{j, i},
		distance: math.MaxInt,
		visited:  false,
	}
	nodes.nodeSlice = append(nodes.nodeSlice, loc)

	*tableStr += "<td " + blueBg + " draggable=\"true\"" + id + ">" + "</td>"

	nodes.endNode = loc

}

func makeNormalNode(j int, i int, elementLocation string, id string, tableStr *string) {

	loc := Node{
		location: [2]int{j, i},
		distance: math.MaxInt,
		visited:  false,
	}
	nodes.nodeSlice = append(nodes.nodeSlice, loc)

	*tableStr += "<td onmouseover=\"makeWall(" + elementLocation + ")\"" + id + ">" + "</td>"

}

func surroundingNodes(currentNode *Node) []*Node {

	var out []*Node

	for i := range nodes.nodeSlice {

		x := currentNode.location[0]
		y := currentNode.location[1]

		switch nodes.nodeSlice[i].location {

		case [2]int{x, y + 1}: //2
			out = append(out, &nodes.nodeSlice[i])
		case [2]int{x - 1, y}: //4
			out = append(out, &nodes.nodeSlice[i])
		case [2]int{x + 1, y}: //5
			out = append(out, &nodes.nodeSlice[i])
		case [2]int{x, y - 1}: //7
			out = append(out, &nodes.nodeSlice[i])

		}

	}

	return out
}

func calcDistance(curNode *Node, neighborNode *Node) int {
	dist := math.Sqrt(math.Pow(float64(neighborNode.location[0]-curNode.location[0]), 2)+math.Pow(float64(neighborNode.location[1]-curNode.location[1]), 2)) * 10.0
	return int(dist)
}

func minDistance() *Node {
	min := math.MaxInt
	var minNode *Node
	for i, node := range nodes.nodeSlice {
		if !node.visited && node.distance <= min {
			min = node.distance
			minNode = &nodes.nodeSlice[i]
		}
	}

	return minNode
}

func path(pathSlice []*Node) {

	if len(pathSlice[len(pathSlice)-1].parent) != 0 {
		pathSlice = append(pathSlice, pathSlice[len(pathSlice)-1].parent[0])
		path(pathSlice)
	}

	//fmt.Println(pathSlice)
	for _, node := range pathSlice {
		if node.location != nodes.startNode.location {
			id := fmt.Sprintf("%v,%v", node.location[0], node.location[1])
			el := dom.GetWindow().Document()
			el.GetElementByID(id).SetAttribute("style", "background-color: yellow;")
		}
	}

}

func dijkstra() int {

	time.Sleep(1 * time.Second)

	for range nodes.nodeSlice {

		tempSrcNode := minDistance()
		neighbors := surroundingNodes(tempSrcNode)
		tempSrcNode.visited = true

		for _, neighborNode := range neighbors {

			neighborDist := calcDistance(tempSrcNode, neighborNode)

			if !neighborNode.visited &&
				tempSrcNode.distance != math.MaxInt &&
				tempSrcNode.distance+neighborDist < neighborNode.distance &&
				!tempSrcNode.wall {

				neighborNode.distance = tempSrcNode.distance + neighborDist
				neighborNode.parent = []*Node{tempSrcNode}
				if tempSrcNode.location != nodes.startNode.location {
					id := fmt.Sprintf("%v,%v", tempSrcNode.location[0], tempSrcNode.location[1])
					el := dom.GetWindow().Document()
					el.GetElementByID(id).SetAttribute("style", "background-color: green;")
				}

				if neighborNode.location == nodes.endNode.location {
					parent := neighborNode.parent[0]
					path([]*Node{parent})
					return 0
				}

				time.Sleep(time.Microsecond / 1000)
			}

			if neighborNode.visited && neighborNode.location != nodes.startNode.location && !neighborNode.wall {
				id := fmt.Sprintf("%v,%v", neighborNode.location[0], neighborNode.location[1])
				el := dom.GetWindow().Document()
				el.GetElementByID(id).SetAttribute("style", "background-color: green;")
			}

		}

	}

	return 0

}

func makeWall(this js.Value, i []js.Value) interface{} {
	x := i[0].Int()
	y := i[1].Int()

	for idx := range nodes.nodeSlice {
		if nodes.nodeSlice[idx].location == [2]int{x, y} && !nodes.nodeSlice[idx].wall && wallsEnabled {
			id := fmt.Sprintf("%v,%v", x, y)
			el := dom.GetWindow().Document()
			el.GetElementByID(id).SetAttribute("style", "background-color: black;")
			nodes.nodeSlice[idx].wall = true
		}
	}

	return 0
}

func run(this js.Value, i []js.Value) interface{} {
	start = true
	return 0
}

func toggleWalls(this js.Value, i []js.Value) interface{} {
	wallsEnabled = !wallsEnabled
	return 0
}

func registerCallbacks() {
	js.Global().Set("makeWall", js.FuncOf(makeWall))
	js.Global().Set("run", js.FuncOf(run))
	js.Global().Set("toggleWalls", js.FuncOf(toggleWalls))
}

func main() {
	c := make(chan struct{})
	println("WASM Go Initialized")

	registerCallbacks()

	el := dom.GetWindow().Document()
	div := el.GetElementByID("table")
	div.SetInnerHTML(makeTable(50, 50))

	for {
		if start {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	dijkstra()

	<-c
}
