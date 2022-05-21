//go:build js && wasm

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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
var startMoving = false
var endMoving = false

func makeTable(rows int, columns int) string {

	tableStr := `<table draggable="false" style="border-collapse: collapse">`

	for i := 0; i < rows; i++ {

		tableStr += `<tr draggable="false">`

		for j := 0; j < columns; j++ {

			if j == ((columns/8)-1) && i == (rows/2) {
				makeStartNode(j, i, &tableStr)
			} else if j == (columns-(columns/8)) && i == (rows/2) {
				makeEndNode(j, i, &tableStr)
			} else {
				makeNormalNode(j, i, &tableStr)
			}

		}

		tableStr += "</tr>"

	}

	tableStr += "</table>"

	return tableStr

}

func makeStartNode(j int, i int, tableStr *string) {

	elementLocation := strconv.Itoa(j) + "," + strconv.Itoa(i)
	id := `id="` + elementLocation + `"`
	redBg := `style="background-color:red"`

	loc := Node{
		location: [2]int{j, i},
		distance: 0,
		visited:  false,
	}
	nodes.startNode = loc
	nodes.nodeSlice = append(nodes.nodeSlice, loc)

	*tableStr += `<td ondragstart="enableMoveStart()" ` + redBg + ` draggable="true" ` + id + "></td>"
}

func makeEndNode(j int, i int, tableStr *string) {

	elementLocation := strconv.Itoa(j) + "," + strconv.Itoa(i)
	id := `id="` + elementLocation + `"`

	blueBg := `style="background-color:blue"`

	loc := Node{
		location: [2]int{j, i},
		distance: math.MaxInt,
		visited:  false,
	}
	nodes.nodeSlice = append(nodes.nodeSlice, loc)

	*tableStr += `<td ondragstart="enableMoveEnd()" ` + blueBg + ` draggable="true" ` + id + "></td>"

	nodes.endNode = loc

}

func makeNormalNode(j int, i int, tableStr *string) {

	elementLocation := strconv.Itoa(j) + "," + strconv.Itoa(i)
	id := `id="` + elementLocation + `"`

	loc := Node{
		location: [2]int{j, i},
		distance: math.MaxInt,
		visited:  false,
	}
	nodes.nodeSlice = append(nodes.nodeSlice, loc)

	*tableStr += `<td onpointerdown="toggleWalls()" ondrop="drop_handler(event)" ondragover="dragoverHandler(event)" ondragenter="dragoverHandler(event)" onmouseover="makeWall(` + elementLocation + `)"` + id + "></td>"

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
			el.GetElementByID(id).SetAttribute("style", "background-color: #f29668;")
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
					el.GetElementByID(id).SetAttribute("style", "background-color: #338070;")
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
				el.GetElementByID(id).SetAttribute("style", "background-color: #338070;")
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

func moveNode(this js.Value, i []js.Value) interface{} {
	fmt.Println(startMoving, endMoving)
	if startMoving {
		locStr := i[0].String()
		locStrArr := strings.Split(locStr, ",")
		x, _ := strconv.Atoi(locStrArr[0])
		y, _ := strconv.Atoi(locStrArr[1])
		var newStartNode = ""

		oldEl := dom.GetWindow().Document().GetElementByID(fmt.Sprintf("%v,%v", nodes.startNode.location[0], nodes.startNode.location[1]))
		oldEl.SetAttribute("onpointerdown", "toggleWalls()")
		oldEl.SetAttribute("ondrop", "drop_handler(event)")
		oldEl.SetAttribute("ondragover", "dragoverHandler(event)")
		oldEl.SetAttribute("ondragenter", "dragoverHandler(event)")
		oldEl.SetAttribute("onmouseover", fmt.Sprintf("makeWall(%v,%v)", nodes.startNode.location[0], nodes.startNode.location[1]))
		oldEl.SetAttribute("style", "")
		oldEl.SetAttribute("ondragstart", "")
		oldEl.SetAttribute("draggable", "false")

		for idx := range nodes.nodeSlice {
			if nodes.nodeSlice[idx].location == [2]int{nodes.startNode.location[0], nodes.startNode.location[1]} {
				nodes.nodeSlice[idx].distance = math.MaxInt
			}
		}

		makeStartNode(x, y, &newStartNode)

		el := dom.GetWindow().Document().GetElementByID(fmt.Sprintf("%v,%v", x, y))
		el.SetAttribute("onpointerdown", "")
		el.SetAttribute("ondrop", "")
		el.SetAttribute("ondragover", "")
		el.SetAttribute("ondragenter", "")
		el.SetAttribute("onmouseover", "")
		el.SetAttribute("style", "background-color:red")
		el.SetAttribute("ondragstart", "enableMoveStart()")
		el.SetAttribute("draggable", "true")

		startMoving = false

	}
	if endMoving {
		locStr := i[0].String()
		locStrArr := strings.Split(locStr, ",")
		x, _ := strconv.Atoi(locStrArr[0])
		y, _ := strconv.Atoi(locStrArr[1])
		var newEndNode = ""

		oldEl := dom.GetWindow().Document().GetElementByID(fmt.Sprintf("%v,%v", nodes.endNode.location[0], nodes.endNode.location[1]))
		oldEl.SetAttribute("onpointerdown", "toggleWalls()")
		oldEl.SetAttribute("ondrop", "drop_handler(event)")
		oldEl.SetAttribute("ondragover", "dragoverHandler(event)")
		oldEl.SetAttribute("ondragenter", "dragoverHandler(event)")
		oldEl.SetAttribute("onmouseover", fmt.Sprintf("makeWall(%v,%v)", nodes.endNode.location[0], nodes.endNode.location[1]))
		oldEl.SetAttribute("style", "")
		oldEl.SetAttribute("ondragstart", "")
		oldEl.SetAttribute("draggable", "false")

		for idx := range nodes.nodeSlice {
			if nodes.nodeSlice[idx].location == [2]int{nodes.endNode.location[0], nodes.endNode.location[1]} {
				nodes.nodeSlice[idx].distance = math.MaxInt
			}
		}

		makeEndNode(x, y, &newEndNode)

		el := dom.GetWindow().Document().GetElementByID(fmt.Sprintf("%v,%v", x, y))
		el.SetAttribute("onpointerdown", "")
		el.SetAttribute("ondrop", "")
		el.SetAttribute("ondragover", "")
		el.SetAttribute("ondragenter", "")
		el.SetAttribute("onmouseover", "")
		el.SetAttribute("style", "background-color:blue")
		el.SetAttribute("ondragstart", "enableMoveEnd()")
		el.SetAttribute("draggable", "true")

		endMoving = false

	}

	return 0
}

func enableMoveStart(this js.Value, i []js.Value) interface{} {
	startMoving = true
	return 0
}

func enableMoveEnd(this js.Value, i []js.Value) interface{} {
	endMoving = true
	return 0
}

func registerCallbacks() {
	js.Global().Set("makeWall", js.FuncOf(makeWall))
	js.Global().Set("run", js.FuncOf(run))
	js.Global().Set("toggleWalls", js.FuncOf(toggleWalls))
	js.Global().Set("moveNode", js.FuncOf(moveNode))
	js.Global().Set("enableMoveStart", js.FuncOf(enableMoveStart))
	js.Global().Set("enableMoveEnd", js.FuncOf(enableMoveEnd))
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
