package datastructures

import (
	"container/heap"
	"pathfinding-algorithms/models"
)

type Item struct {
	node     *models.Node // The value of the item; arbitrary.
	priority float64      // The priority of the item in the queue.
	index    int          // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority // The less function compares fScores to determine priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Update modifies the priority and reorders the queue.
func (pq *PriorityQueue) Update(node *models.Node, newPriority float64) {
	for i, item := range *pq {
		if item.node == node {
			// Node found, update its priority
			(*pq)[i].priority = newPriority
			// Fix the heap since the priority is updated
			heap.Fix(pq, i)
			return
		}
	}
	// If the node is not found, it might be necessary to add it to the queue
	// depending on your specific requirements. This part is left as an exercise.
}

func (pq PriorityQueue) Contains(node *models.Node) bool {
	for _, item := range pq {
		if item.node == node {
			return true
		}
	}
	return false
}

// NewItem creates a new Item.
func NewItem(node *models.Node, priority float64) *Item {
	return &Item{
		node:     node,
		priority: priority,
	}
}

func (item *Item) GetNode() *models.Node {
	return item.node
}

func (item *Item) GetPriority() float64 {
	return item.priority
}

// Init initializes or clears the priority queue.
func (pq *PriorityQueue) Init() {
	heap.Init(pq)
}
