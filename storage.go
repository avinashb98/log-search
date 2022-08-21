package main

import (
	"fmt"
	"sort"
)

type LogsStorage map[LogID]Log

func (l *LogsStorage) Print() {
	for _, log := range *l {
		fmt.Println(log.String())
	}
}

type Storage struct {
	logsStorage LogsStorage
	index       InvertedIndex
	buffer      Buffer
	capacity    int
}

func getNewStore(s int) *Storage {
	return &Storage{
		logsStorage: LogsStorage{},
		index:       getNewIndex(),
		buffer:      getNewBuffer(),
		capacity:    s,
	}
}

func (s *Storage) upsertLog(id LogID, data string) {
	existingLog, err := s.getLogById(id)
	if err != nil {
		newLog := getNewLog(id, data)
		s.addLog(newLog, true)
		return
	}
	updatedLog := existingLog.copy()
	updatedLog.Data = data
	s.updateLog(existingLog, updatedLog)
}

func (s *Storage) updateLog(prevLog, updatedLog Log) {
	s.addLog(updatedLog, false)
	opts := UpdateOpts{previous: &prevLog, current: &updatedLog}
	s.index.update(opts)
}

func (s *Storage) addLog(log Log, updateBuffer bool) {
	s.logsStorage[log.ID] = log
	opts := UpdateOpts{current: &log}
	s.index.update(opts)
	if updateBuffer {
		s.buffer.Enqueue(&log.ID)
	}

	s.cleanup()

}

func (s *Storage) getLogsByWord(word string, limit int) []Log {
	logIds := s.index.getByKey(word)
	if logIds == nil {
		return nil
	}
	var logs []Log
	for i := len(logIds) - 1; i >= 0; i-- {
		id := logIds[i]
		log, err := s.getLogById(id)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.After(logs[j].CreatedAt)
	})
	if len(logs) < limit {
		limit = len(logs)
	}
	return logs[:limit]
}

func (s *Storage) getLogById(id LogID) (Log, error) {
	log, found := s.logsStorage[id]
	if !found {
		return Log{}, fmt.Errorf("document not found")
	}
	return log, nil
}

func (s *Storage) truncate() {
	for s.buffer.Len() > s.capacity {
		lastElem := s.buffer.Dequeue()
		if lastElem == nil {
			break
		}
		s.deleteLogById(*lastElem)
	}
}

func (s *Storage) deleteLogById(id LogID) {
	s.index.deletedByLogId(id)
	delete(s.logsStorage, id)
}

func (s *Storage) cleanup() {
	if s.buffer.Len() > s.capacity {
		s.truncate()
	}
}
