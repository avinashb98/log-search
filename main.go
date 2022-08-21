package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	dat, err := os.ReadFile("input_sample.txt")
	check(err)
	commands := strings.Split(string(dat), "\n")
	filteredCommands := []string{}
	for _, command := range commands {
		if command != "" {
			filteredCommands = append(filteredCommands, command)
		}
	}

	inputFileDriver(filteredCommands)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func inputFileDriver(commands []string) {
	if len(commands) < 2 {
		panic("the input should have at least 2 commands")
	}
	storeLimit, err := strconv.Atoi(commands[0])
	if err != nil {
		panic("invalid storage limit passed")
	}
	endCommand := commands[len(commands)-1]
	if endCommand != "END" {
		panic("no end received, the last command has to be END")
	}
	store := getNewStore(storeLimit)
	commands = commands[1:]
	for _, command := range commands {
		processCommand(store, command, os.Stdout)
	}
}

func processCommand(store *Storage, command string, output io.Writer) {
	if command == "END" {
		processEnd(output)
		return
	}

	if len(command) >= 3 && command[:3] == "ADD" {
		processAdd(store, command)
		return
	}

	if len(command) >= 6 && command[:6] == "SEARCH" {
		processSearch(store, command, output)
		return
	}

	panic("invalid command")

}

func processEnd(output io.Writer) {
	output.Write([]byte("END\r\n"))
}

func processAdd(store *Storage, command string) {
	arguments := strings.Split(command, " ")
	idLen := len(arguments[1])
	logId, err := strconv.Atoi(arguments[1])
	if err != nil {
		fmt.Print("invalid key id\r\n")
	}
	data := command[5+idLen:]
	store.upsertLog(LogID(logId), data)
}

func processSearch(store *Storage, command string, output io.Writer) {
	arguments := strings.Split(command, " ")
	query := arguments[1]
	limit, err := strconv.Atoi(arguments[2])
	if err != nil {
		fmt.Print("invalid limit")
	}
	logs := store.getLogsByWord(query, limit)
	if logs == nil || len(logs) == 0 {
		output.Write([]byte("NONE\r\n"))
		return
	}
	logIds := []string{}
	for _, log := range logs {
		logIds = append(logIds, strconv.Itoa(int(log.ID)))
	}
	output.Write([]byte(strings.Join(logIds, " ") + "\r\n"))
}
