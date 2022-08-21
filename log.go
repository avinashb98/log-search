package main

import (
	"fmt"
	"time"
)

type LogID int

type Log struct {
	ID                LogID
	Data              string
	CreatedAt         time.Time
	MarkedForDeletion bool
}

func (l Log) copy() Log {
	return Log{
		ID:                l.ID,
		Data:              l.Data,
		CreatedAt:         l.CreatedAt,
		MarkedForDeletion: l.MarkedForDeletion,
	}
}

func (l Log) String() string {
	return fmt.Sprintf("id: %d createdAt: %v data: %s", l.ID, l.CreatedAt, l.Data)
}

func getNewLog(id LogID, data string) Log {
	return Log{
		ID:        id,
		Data:      data,
		CreatedAt: time.Now(),
	}
}
