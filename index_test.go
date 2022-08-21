package main

import (
	"reflect"
	"testing"
)

func TestInvertedIndex_update(t *testing.T) {
	type fields struct {
		keyToEntries map[string][]LogID
		entryToKeys  map[LogID][]string
	}
	type args struct {
		opts UpdateOpts
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"prev is nil",
			fields{
				keyToEntries: map[string][]LogID{},
				entryToKeys:  map[LogID][]string{},
			},
			args{opts: UpdateOpts{
				previous: nil,
				current:  &Log{ID: 123, Data: "hello world"},
			}},
		},
		{
			"prev is not nil",
			fields{
				keyToEntries: map[string][]LogID{},
				entryToKeys:  map[LogID][]string{},
			},
			args{opts: UpdateOpts{
				previous: &Log{ID: 123, Data: "hello world"},
				current:  &Log{ID: 123, Data: "hello again"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				keyToEntries: tt.fields.keyToEntries,
				entryToKeys:  tt.fields.entryToKeys,
			}
			i.update(tt.args.opts)
		})
	}
}

func TestInvertedIndex_updateEntries(t *testing.T) {
	type fields struct {
		keyToEntries map[string][]LogID
		entryToKeys  map[LogID][]string
	}
	type args struct {
		log *Log
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"log in nil",
			fields{},
			args{log: nil},
		},
		{
			"log is not nil",
			fields{
				keyToEntries: map[string][]LogID{},
				entryToKeys:  map[LogID][]string{},
			},
			args{log: &Log{
				ID:   123,
				Data: "hello world",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				keyToEntries: tt.fields.keyToEntries,
				entryToKeys:  tt.fields.entryToKeys,
			}
			i.updateEntries(tt.args.log)
		})
	}
}

func TestInvertedIndex_removeMappings(t *testing.T) {
	type fields struct {
		keyToEntries map[string][]LogID
		entryToKeys  map[LogID][]string
	}
	type args struct {
		prev    *Log
		current *Log
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			"no update",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
			args{
				prev: &Log{
					ID:   123,
					Data: "hello world",
				},
				current: &Log{
					ID:   123,
					Data: "hello world",
				},
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
		},
		{
			"complete update",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
			args{
				prev: &Log{
					ID:   123,
					Data: "hello world",
				},
				current: &Log{
					ID:   123,
					Data: "something else",
				},
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {}, "world": {}},
				entryToKeys:  map[LogID][]string{123: {}},
			},
		},
		{
			"partial update",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
			args{
				prev: &Log{
					ID:   123,
					Data: "hello world",
				},
				current: &Log{
					ID:   123,
					Data: "hello again",
				},
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {}},
				entryToKeys:  map[LogID][]string{123: {"hello"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				keyToEntries: tt.fields.keyToEntries,
				entryToKeys:  tt.fields.entryToKeys,
			}
			i.removeMappings(tt.args.prev, tt.args.current)
			if !reflect.DeepEqual(i.entryToKeys, tt.want.entryToKeys) {
				t.Errorf("removeMappings() got = %v, want %v", i.entryToKeys, tt.want.entryToKeys)
			}
			if !reflect.DeepEqual(i.keyToEntries, tt.want.keyToEntries) {
				t.Errorf("removeMappings() got = %v, want %v", i.keyToEntries, tt.want.keyToEntries)
			}
		})
	}
}

func TestInvertedIndex_deletedByLogId(t *testing.T) {
	type fields struct {
		keyToEntries map[string][]LogID
		entryToKeys  map[LogID][]string
	}
	type args struct {
		id LogID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			"log exists",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
			args{id: 123},
			fields{
				keyToEntries: map[string][]LogID{"hello": {}, "world": {}},
				entryToKeys:  map[LogID][]string{},
			},
		},
		{
			"log doesn't exist",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
			args{id: 456},
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				keyToEntries: tt.fields.keyToEntries,
				entryToKeys:  tt.fields.entryToKeys,
			}
			i.deletedByLogId(tt.args.id)
			if !reflect.DeepEqual(i.entryToKeys, tt.want.entryToKeys) {
				t.Errorf("deletedByLogId() = %v, want %v", i.entryToKeys, tt.want.entryToKeys)
			}
			if !reflect.DeepEqual(i.keyToEntries, tt.want.keyToEntries) {
				t.Errorf("deletedByLogId() = %v, want %v", i.keyToEntries, tt.want.keyToEntries)
			}

		})
	}
}

func TestInvertedIndex_getByKey(t *testing.T) {
	type fields struct {
		keyToEntries map[string][]LogID
		entryToKeys  map[LogID][]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []LogID
	}{
		{
			"no entry exists",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
			args{
				key: "something",
			},
			nil,
		},
		{
			"single entry exists",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
			args{
				key: "hello",
			},
			[]LogID{123},
		},
		{
			"multiple entries exist",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123, 345, 567}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}, 345: {"hello"}, 567: {"hello"}},
			},
			args{
				key: "hello",
			},
			[]LogID{123, 345, 567},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				keyToEntries: tt.fields.keyToEntries,
				entryToKeys:  tt.fields.entryToKeys,
			}
			if got := i.getByKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvertedIndex_removeKeysFromEntry(t *testing.T) {
	type fields struct {
		entryToKeys map[LogID][]string
	}
	type args struct {
		keysToBeRemoved []string
		id              LogID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			"entry doesn't exist",
			fields{
				entryToKeys: map[LogID][]string{123: {"hello", "world"}, 345: {"hello"}, 567: {"hello"}},
			},
			args{
				keysToBeRemoved: []string{"hello"},
				id:              789,
			},
			fields{
				entryToKeys: map[LogID][]string{123: {"hello", "world"}, 345: {"hello"}, 567: {"hello"}},
			},
		},
		{
			"entry exists",
			fields{
				entryToKeys: map[LogID][]string{123: {"hello", "world"}, 345: {"hello"}, 567: {"hello"}},
			},
			args{
				keysToBeRemoved: []string{"hello"},
				id:              123,
			},
			fields{
				entryToKeys: map[LogID][]string{123: {"world"}, 345: {"hello"}, 567: {"hello"}},
			},
		},
		{
			"entry exists but key in it doesn't",
			fields{
				entryToKeys: map[LogID][]string{123: {"world"}, 345: {"hello"}, 567: {"hello"}},
			},
			args{
				keysToBeRemoved: []string{"hello"},
				id:              123,
			},
			fields{
				entryToKeys: map[LogID][]string{123: {"world"}, 345: {"hello"}, 567: {"hello"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				entryToKeys: tt.fields.entryToKeys,
			}
			i.removeKeysFromEntry(tt.args.keysToBeRemoved, tt.args.id)
			if !reflect.DeepEqual(i.entryToKeys, tt.want.entryToKeys) {
				t.Errorf("removeKeysFromEntry() got = %v, want %v", i.entryToKeys, tt.want.entryToKeys)
			}
		})
	}
}

func TestInvertedIndex_removeEntryFromKeys(t *testing.T) {
	type fields struct {
		keyToEntries map[string][]LogID
	}
	type args struct {
		keys []string
		id   LogID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			"key doesn't exist",
			fields{
				keyToEntries: map[string][]LogID{"world": {123}},
			},
			args{
				keys: []string{"hello"},
				id:   789,
			},
			fields{
				keyToEntries: map[string][]LogID{"world": {123}},
			},
		},
		{
			"key exists",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123, 345, 567}, "world": {123}},
			},
			args{
				keys: []string{"hello"},
				id:   123,
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {345, 567}, "world": {123}},
			},
		},
		{
			"keys exists but entry in it doesn't",
			fields{
				keyToEntries: map[string][]LogID{"hello": {345, 567}, "world": {123}},
			},
			args{
				keys: []string{"hello"},
				id:   123,
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {345, 567}, "world": {123}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				keyToEntries: tt.fields.keyToEntries,
			}
			i.removeEntryFromKeys(tt.args.keys, tt.args.id)
			if !reflect.DeepEqual(i.keyToEntries, tt.want.keyToEntries) {
				t.Errorf("removeEntryFromKeys() got = %v, want %v", i.keyToEntries, tt.want.keyToEntries)
			}
		})
	}
}

func TestInvertedIndex_updateEntry(t *testing.T) {
	type fields struct {
		keyToEntries map[string][]LogID
		entryToKeys  map[LogID][]string
	}
	type args struct {
		key string
		id  LogID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			"no prior key or entry",
			fields{
				keyToEntries: map[string][]LogID{},
				entryToKeys:  map[LogID][]string{},
			},
			args{
				key: "hello",
				id:  123,
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello"}},
			},
		},
		{
			"existing entry, no prior key",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello"}},
			},
			args{
				key: "world",
				id:  123,
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}, "world": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello", "world"}},
			},
		},
		{
			"existing key, no prior entry",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello"}},
			},
			args{
				key: "hello",
				id:  456,
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {123, 456}},
				entryToKeys:  map[LogID][]string{123: {"hello"}, 456: {"hello"}},
			},
		},
		{
			"existing key and entry",
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello"}},
			},
			args{
				key: "hello",
				id:  123,
			},
			fields{
				keyToEntries: map[string][]LogID{"hello": {123}},
				entryToKeys:  map[LogID][]string{123: {"hello"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InvertedIndex{
				keyToEntries: tt.fields.keyToEntries,
				entryToKeys:  tt.fields.entryToKeys,
			}
			i.updateEntry(tt.args.key, tt.args.id)
			if !reflect.DeepEqual(i.entryToKeys, tt.want.entryToKeys) {
				t.Errorf("updateEntry() = %v, want %v", i.entryToKeys, tt.want.entryToKeys)
			}
			if !reflect.DeepEqual(i.keyToEntries, tt.want.keyToEntries) {
				t.Errorf("updateEntry() = %v, want %v", i.keyToEntries, tt.want.keyToEntries)
			}
		})
	}
}

func Test_getNewIndex(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"properties shouldn't be null",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNewIndex()
			if got.keyToEntries == nil {
				t.Errorf("getNewIndex() returned nil keyToEntries")
			}
			if got.entryToKeys == nil {
				t.Errorf("getNewIndex() returned nil entryToKeys")
			}
		})
	}
}

func Test_getWordsDelta(t *testing.T) {
	type args struct {
		prev []string
		curr []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"no delta",
			args{
				prev: []string{"hello", "world"},
				curr: []string{"hello", "world"},
			},
			[]string{},
		},
		{
			"delta = 1",
			args{
				prev: []string{"hello", "world"},
				curr: []string{"world"},
			},
			[]string{"hello"},
		},
		{
			"first argument nil",
			args{
				prev: nil,
				curr: []string{"hello", "world"},
			},
			[]string{},
		},
		{
			"second argument nil",
			args{
				prev: []string{"hello", "world"},
				curr: nil,
			},
			[]string{"hello", "world"},
		},
		{
			"both arguments nil",
			args{
				prev: nil,
				curr: nil,
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWordsDelta(tt.args.prev, tt.args.curr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWordsDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getWordsFromData(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"add 2 words", args{data: "hello world"}, []string{"hello", "world"},
		},
		{
			"add 9 words", args{data: "a quick brown fox jumps over the lazy dog"}, []string{"a", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"},
		},
		{
			"multi space", args{data: "hello   world"}, []string{"hello", "world"},
		},
		{
			"empty", args{data: ""}, []string{},
		},
		{
			"no space", args{data: "helloworld"}, []string{"helloworld"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWordsFromData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWordsFromData() = %v, want %v", got, tt.want)
			}
		})
	}
}
