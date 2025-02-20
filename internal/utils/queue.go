package utils

import (
	"sync"
)

type TCPBuffData struct {
	Counter uint64
	Data    []byte
}

type OrderedQueue interface {
	Add(*TCPBuffData)
	Pop() *TCPBuffData
	Length() int
}

func NewQueue() OrderedQueue {
	return &orderedQueue{
		data: make([]*TCPBuffData, 0),
	}
}

type orderedQueue struct {
	data []*TCPBuffData
	mu   sync.Mutex
}

func (oq *orderedQueue) Add(data *TCPBuffData) {
	oq.mu.Lock()
	defer oq.mu.Unlock()

	index := oq.findInsertIndex(data.Counter)
	oq.data = append(oq.data, nil)
	copy(oq.data[index+1:], oq.data[index:])
	oq.data[index] = data
}

func (oq *orderedQueue) findInsertIndex(counter uint64) int {
	low, high := 0, len(oq.data)-1
	for low <= high {
		mid := low + (high-low)/2
		if oq.data[mid].Counter < counter {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return low
}

func (oq *orderedQueue) Pop() *TCPBuffData {
	oq.mu.Lock()
	defer oq.mu.Unlock()

	if len(oq.data) == 0 {
		return nil
	}
	d := oq.data[0]
	oq.data = oq.data[1:]

	return d
}

func (oq *orderedQueue) Length() int {
	oq.mu.Lock()
	defer oq.mu.Unlock()
	return len(oq.data)
}
