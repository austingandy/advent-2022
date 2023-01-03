package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	logger := log.Logger{}
	logger.SetOutput(os.Stdout)
	f, err := os.Open("./input.txt")
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	commands := make([]Command, 0)
	for scanner.Scan() {
		input := strings.Split(scanner.Text(), " ")
		switch commandType := toCommandType(input[0]); commandType {
		case AddX:
			v, _ := strconv.Atoi(input[1])
			commands = append(commands, Command{CommandType: commandType, Val: v})
		case Noop:
			commands = append(commands, Command{CommandType: Noop})
		case Unknown:
			logger.Fatalf("unrecognized commandType: %s", string(commandType))
		}
	}
	cycleCount, x, signalStrength, lit, dim, out := 0, 1, 0, "#", " ", ""
	for _, command := range commands {
		var numCycles int
		if command.CommandType == Noop {
			numCycles = 1
		} else {
			numCycles = 2
		}
		for i := 0; i < numCycles; i += 1 {
			if math.Abs(float64((len(out)%40)-x)) < float64(2) {
				out += lit
			} else {
				out += dim
			}
			cycleCount += 1
			if cycleCount == 20 || cycleCount > 20 && (cycleCount-20)%40 == 0 {
				signalStrength += x * cycleCount
			}
		}
		x += command.Val
	}
	logger.Printf("Total (Part 1): %d", signalStrength)
	logger.Printf("Part 2 (squint for letters):")
	for i := 0; i < len(out)/40; i += 1 {
		logger.Printf("%s\n", out[40*i:40*(i+1)])
	}
}

type CommandType string

const Noop CommandType = "noop"
const AddX CommandType = "addx"
const Unknown CommandType = "unknown"

type Command struct {
	CommandType CommandType
	Val         int
}

func toCommandType(s string) CommandType {
	switch s {
	case "noop":
		return Noop
	case "addx":
		return AddX
	default:
		return Unknown
	}
}
