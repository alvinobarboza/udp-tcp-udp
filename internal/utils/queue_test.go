package utils

import (
	"testing"
)

func TestQueue(t *testing.T) {
	t.Run("Add item", func(tt *testing.T) {
		queue := NewQueue()
		queue.Add(&TCPBuffData{})
		queue.Add(&TCPBuffData{})
		queue.Add(&TCPBuffData{})
		queue.Add(&TCPBuffData{})
	})
	t.Run("Pop more them have", func(tt *testing.T) {
		queue := NewQueue()
		queue.Pop()

		queue.Add(&TCPBuffData{})
		queue.Add(&TCPBuffData{})

		queue.Pop()
		queue.Pop()
		queue.Pop()
	})
}
