package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	logger := log.Logger{}
	f, err := os.Open("./input.txt")
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	score := 0
	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		theirs, err := throw(input[0])
		if err != nil {
			logger.Fatalf("Failed to parse their throw: %s", err.Error())
		}
		mine, err := myThrow(theirs, input[1])
		if err != nil {
			logger.Fatalf("Failed to parse my throw: %s", err.Error())
		}
		score += scoreGame(mine, theirs)
	}
	fmt.Printf("Score: %d\n", score)
}

type Throw string
type Outcome string

var wins = map[Throw]Throw{
	"R": "P",
	"P": "S",
	"S": "R",
}

var losses = map[Throw]Throw{
	"R": "S",
	"P": "R",
	"S": "P",
}

func myThrow(theirs Throw, c string) (Throw, error) {
	if c == "Y" {
		return theirs, nil
	}
	if c == "Z" {
		win, ok := wins[theirs]
		if !ok {
			return "", fmt.Errorf("invalid throw: %s", theirs)
		}
		return win, nil
	}
	if c == "X" {
		loss, ok := losses[theirs]
		if !ok {
			return "", fmt.Errorf("invalid throw: %s", theirs)
		}
		return loss, nil
	}
	return "", fmt.Errorf("invalid choice: %s", c)
}

func throw(c string) (Throw, error) {
	if c == "X" || c == "A" {
		return "R", nil
	}
	if c == "Y" || c == "B" {
		return "P", nil
	}
	if c == "Z" || c == "C" {
		return "S", nil
	}
	return "", fmt.Errorf("invalid choice: %s", c)
}

func scoreThrow(t Throw) int {
	if t == "R" {
		return 1
	}
	if t == "P" {
		return 2
	}
	if t == "S" {
		return 3
	}
	return 0
}

func getOutcome(a, b Throw) Outcome {
	if a == b {
		return "T"
	}
	if wins[b] == a {
		return "W"
	}
	return "L"
}

func scoreGame(a, b Throw) int {
	throwScore := scoreThrow(a)
	outcome := getOutcome(a, b)
	if outcome == "W" {
		return 6 + throwScore
	}
	if outcome == "T" {
		return 3 + throwScore
	}
	return throwScore
}
