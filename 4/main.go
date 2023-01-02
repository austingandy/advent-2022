package main

import (
	"bufio"
	"log"
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
	t1 := 0
	t2 := 0
	for scanner.Scan() {
		input := scanner.Text()
		rs := strings.Split(input, ",")
		if len(rs) != 2 {
			logger.Fatalf("Expected 2 ranges, got %d", len(rs))
		}
		sn1, sn2 := strings.Split(rs[0], "-"), strings.Split(rs[1], "-")
		if len(sn1) != 2 || len(sn2) != 2 {
			logger.Fatalf("Expected 2 numbers in each range, got %d and %d", len(sn1), len(sn2))
		}
		start1, err := strconv.Atoi(sn1[0])
		if err != nil {
			logger.Fatalf("Failed to parse number: %s", err.Error())
		}
		e1, err := strconv.Atoi(sn1[1])
		if err != nil {
			logger.Fatalf("Failed to parse number: %s", err.Error())
		}
		start2, err := strconv.Atoi(sn2[0])
		if err != nil {
			logger.Fatalf("Failed to parse number: %s", err.Error())
		}
		e2, err := strconv.Atoi(sn2[1])
		if err != nil {
			logger.Fatalf("Failed to parse number: %s", err.Error())
		}
		if start1 <= start2 && e1 >= e2 || start2 <= start1 && e2 >= e1 {
			t1 += 1
		}
		var gs, le int
		if start1 <= start2 {
			le = e1
			gs = start2
		} else {
			le = e2
			gs = start1
		}
		if le >= gs {
			t2 += 1
		}
	}
	logger.Printf("Total (part 1): %d\n", t1)
	logger.Printf("Total (part 2): %d\n", t2)
}
