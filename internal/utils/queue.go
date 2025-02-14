package utils

import (
	"sync"
)

type TCPBuffData struct {
	Counter uint64
	MS      uint32
	Data    []byte
}

type OrderedQueue interface {
	Add(*TCPBuffData)
	Pop() *TCPBuffData
	Length() int
}

func NewQueue() OrderedQueue {
	return &orderedQueue{}
}

type orderedQueue struct {
	data []*TCPBuffData
	mu   sync.Mutex
}

func (ld *orderedQueue) Add(data *TCPBuffData) {
	ld.mu.Lock()
	ld.data = append(ld.data, data)
	for i, d := range ld.data {
		for j, d2 := range ld.data {
			if d.Counter > d2.Counter {
				temp := ld.data[i]
				ld.data[i] = ld.data[j]
				ld.data[j] = temp
			}
		}
	}
	ld.mu.Unlock()
}

func (ld *orderedQueue) Pop() *TCPBuffData {
	ld.mu.Lock()
	defer ld.mu.Unlock()

	i := len(ld.data) - 1
	if i < 0 {
		return nil
	}
	d := ld.data[i]
	ld.data = ld.data[:i]

	return d
}

func (ls *orderedQueue) Length() int {
	return len(ls.data)
}
