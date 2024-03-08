package datastructures

type Queue []interface{}

// Enqueue adds an item to the end of the queue
func (q *Queue) Enqueue(item interface{}) {
	*q = append(*q, item)
}

// Dequeue removes an item from the start of the queue and returns it
// Returns nil if the queue is empty
func (q *Queue) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item
}

// IsEmpty checks if the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

// Size returns the number of items in the queue
func (q *Queue) Size() int {
	return len(*q)
}
