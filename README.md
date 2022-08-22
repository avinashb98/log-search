# log-search
in memory logs search

## Installation
```shell
git clone https://github.com/avinashb98/log-search
cd log-search
```

## Usage

### commands
#### ADD
* O(n) worst-case, n is the total number of logs 
* O(1) best-case 
```shell
ADD [key] [text] 
```
#### SEARCH
Time Complexities
* O(n) worst-case, n is the total number of logs 
* O(1) best-case

```shell
SEARCH [word] [limit]
```
### input file format
```shell
# input.txt

n # number which is the maximum logs stored
.
.
command
.
.
END
```

### Build and Run
```shell
cp input_sample.txt input.txt
# edit the sample commands file

go build .
./log-search
```

## Test
```shell
go test -v ./...
```

## Design
### Inverted Index
The inverted index consists of 2 data structures.
#### KeyToEntries
It is a map of a word to a list of entryIds.
Used to optimally query entries corresponding to a word.
#### EntryToKeys
It is a map of an entryId to a list of keys.
Used to optimally unmap a deleted or updated entry.
