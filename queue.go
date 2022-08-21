package main

import "container/list"

type Buffer struct {
	list *list.List
}

func getNewBuffer() Buffer {
	return Buffer{
		list: list.New(),
	}
}

func (q *Buffer) Enqueue(item *LogID) {
	q.list.PushFront(item)
}

func (q *Buffer) Len() int {
	return q.list.Len()
}

func (q *Buffer) Dequeue() *LogID {
	lastElem := q.list.Back()
	if lastElem == nil {
		return nil
	}
	q.list.Remove(lastElem)
	return lastElem.Value.(*LogID)
}
