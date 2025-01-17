package utils

import "sync"

// DefaultIndexQueueCap is the default capacity for an index queue.
// It's fairly guaranteed to prevent buffer full conditions during regular use.
// Buffer full wil probably still occur during simulations where there are many spins using the same queue.
// For such use-cases it's probably wise to reset the queue at random intervals.
var DefaultIndexQueueCap = 256

// IndexQueue is a buffered queue for de-duplicating index values.
// It is used for holding duplicate values until they can be used again.
// Note that when the buffer is full, new duplicate values will be discarded!
// Using a buffer prevents memory allocations, and keeps the mechanism as fast as possible.
// IndexQueue is not safe to use across multiple go-routines.
type IndexQueue struct {
	buf    Indexes
	bufMax int
	putPtr int
	getPtr int
	getMax int
}

// NewIndexQueue instantiates a new index queue from the memory pool.
func NewIndexQueue(capacity int) *IndexQueue {
	q := indexQueuePool.Get().(*IndexQueue)
	q.buf = PurgeIndexes(q.buf, capacity)[:capacity]
	q.bufMax = capacity
	return q.Reset()
}

// Release returns the index queue to the memory pool.
func (q *IndexQueue) Release() {
	if q != nil {
		q.buf = q.buf[:0]
		indexQueuePool.Put(q)
	}
}

// Reset discards all values in the queue and returns it to the empty state.
func (q *IndexQueue) Reset() *IndexQueue {
	q.putPtr = 0
	q.getPtr = 0
	q.getMax = 0
	return q
}

// Put pushes a new index value onto the queue.
// If the queue is full the value is discarded.
func (q *IndexQueue) Put(i Index) {
	if q.putPtr >= q.bufMax {
		// wrap around.
		q.putPtr = 0
	}

	if q.putPtr == q.getPtr && q.getPtr < q.getMax {
		// buffer full; too many entries so we drop all!
		q.Reset()
	}

	q.buf[q.putPtr] = i
	q.putPtr++

	if q.putPtr > q.getPtr {
		// can only move if it's not at the end already.
		q.getMax++
	}
}

// Available returns the number of available values in the queue.
func (q *IndexQueue) Available() int {
	if q.putPtr >= q.getMax {
		return q.getMax - q.getPtr
	}
	return q.getMax - q.getPtr + q.putPtr
}

// Get retrieves the next available index value from the queue.
// The returned flag is true if a value was available, otherwise false.
func (q *IndexQueue) Get() (Index, bool) {
	if q.getPtr < q.getMax {
		i := q.buf[q.getPtr]
		q.getPtr++
		return i, true
	}

	if q.getMax >= q.bufMax {
		// wrap around.
		q.getMax = q.putPtr
		q.getPtr = 0
	}

	if q.getPtr == q.getMax {
		return 0, false
	}

	i := q.buf[q.getPtr]
	q.getPtr++
	return i, true
}

// indexQueuePool is the memory pool for index queues.
var indexQueuePool = sync.Pool{New: func() interface{} { return &IndexQueue{} }}
