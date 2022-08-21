package main

import "strings"

type UpdateOpts struct {
	current  *Log
	previous *Log
}

type InvertedIndex struct {
	keyToEntries map[string][]LogID
	entryToKeys  map[LogID][]string
}

func getNewIndex() InvertedIndex {
	return InvertedIndex{
		keyToEntries: map[string][]LogID{},
		entryToKeys:  map[LogID][]string{},
	}
}

func (i *InvertedIndex) update(opts UpdateOpts) {
	i.updateEntries(opts.current)
	if opts.previous != nil {
		i.removeMappings(opts.previous, opts.current)
	}
}

func (i *InvertedIndex) removeMappings(prev, current *Log) {
	prevWords := getWordsFromData(prev.Data)
	currWords := getWordsFromData(current.Data)
	keysDelta := getWordsDelta(prevWords, currWords)
	i.removeKeysFromEntry(keysDelta, prev.ID)
	i.removeEntryFromKeys(keysDelta, prev.ID)
}

func (i *InvertedIndex) updateEntries(log *Log) {
	if log == nil {
		return
	}
	words := getWordsFromData(log.Data)
	for _, word := range words {
		i.updateEntry(word, log.ID)
	}
}

func (i *InvertedIndex) updateEntry(key string, id LogID) {
	_, found := i.keyToEntries[key]
	if !found {
		i.keyToEntries[key] = []LogID{}
	}

	_, found = i.entryToKeys[id]
	if !found {
		i.entryToKeys[id] = []string{}
	}
	entryExists := false
	for _, entry := range i.keyToEntries[key] {
		if entry == id {
			entryExists = true
			break
		}
	}
	if !entryExists {
		i.keyToEntries[key] = append(i.keyToEntries[key], id)
	}

	keyExists := false
	for _, storedKey := range i.entryToKeys[id] {
		if storedKey == key {
			keyExists = true
			break
		}
	}
	if !keyExists {
		i.entryToKeys[id] = append(i.entryToKeys[id], key)
	}
}

func (i *InvertedIndex) getByKey(key string) []LogID {
	entries, found := i.keyToEntries[key]
	if !found {
		return nil
	}
	return entries
}

func (i *InvertedIndex) deletedByLogId(id LogID) {
	keys, found := i.entryToKeys[id]
	if !found {
		return
	}
	delete(i.entryToKeys, id)
	i.removeEntryFromKeys(keys, id)
}

func (i *InvertedIndex) removeEntryFromKeys(keys []string, id LogID) {
	for _, key := range keys {
		if _, found := i.keyToEntries[key]; !found {
			continue
		}
		filteredEntries := []LogID{}
		for _, logId := range i.keyToEntries[key] {
			if logId != id {
				filteredEntries = append(filteredEntries, logId)
			}
		}
		i.keyToEntries[key] = filteredEntries
	}
}

func (i *InvertedIndex) removeKeysFromEntry(keysToBeRemoved []string, id LogID) {
	storedKeys, found := i.entryToKeys[id]
	if !found {
		return
	}
	shouldBeRemoved := map[string]struct{}{}

	for _, key := range keysToBeRemoved {
		shouldBeRemoved[key] = struct{}{}
	}

	filteredKeys := []string{}

	for _, key := range storedKeys {
		if _, yes := shouldBeRemoved[key]; !yes {
			filteredKeys = append(filteredKeys, key)
		}
	}
	i.entryToKeys[id] = filteredKeys
}

func getWordsFromData(data string) []string {
	tokens := strings.Split(data, " ")
	words := []string{}

	for _, token := range tokens {
		if token != "" {
			words = append(words, token)
		}
	}
	return words
}

func getWordsDelta(prev, curr []string) []string {
	delta := []string{}
	existing := map[string]struct{}{}
	for _, word := range curr {
		existing[word] = struct{}{}
	}

	for _, word := range prev {
		if _, found := existing[word]; !found {
			delta = append(delta, word)
		}
	}

	return delta
}
