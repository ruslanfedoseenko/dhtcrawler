package Services

import (
	"container/heap"
	"sync"
)

type Priority int32

const (
	High   Priority = 2
	Normal Priority = 1
	Low    Priority = 0
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    work     // The value of the item; arbitrary.
	priority Priority // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue struct {
	sync.RWMutex
	items        []*Item
	HasMoreItems chan bool
}

func newPriorityQueue(initialCapacity int) *PriorityQueue {
	return &PriorityQueue{
		items:        make([]*Item, initialCapacity),
		HasMoreItems: make(chan bool),
	}
}

func (pq *PriorityQueue) SetItem(index int, item *Item) {
	pq.Lock()
	defer pq.Unlock()
	pq.items[index] = item
}
func (pq PriorityQueue) Len() int {
	pq.RLock()
	defer pq.RUnlock()
	return len(pq.items)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	pq.RLock()
	defer pq.RUnlock()
	return pq.items[i].priority > pq.items[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.Lock()
	defer pq.Unlock()
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	pq.Lock()
	defer pq.Unlock()
	n := len(pq.items)
	item := x.(*Item)
	item.index = n
	pq.items = append(pq.items, item)
	if n == 0 {
		select {
		case pq.HasMoreItems <- true:
			{

			}
		default:
			{

			}
		}
	}
}

func (pq *PriorityQueue) Pop() interface{} {
	pq.Lock()
	defer pq.Unlock()
	old := *pq
	n := len(old.items)
	item := old.items[n-1]
	item.index = -1 // for safety
	pq.items = old.items[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value work, priority Priority) {
	pq.Lock()
	defer pq.Unlock()
	item.value = value
	item.priority = priority
	heap.Fix(pq, int(item.index))
}
