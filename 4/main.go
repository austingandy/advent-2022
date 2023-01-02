package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	logger := log.Logger{}
	f, err := os.Open("./input.txt")
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
}
