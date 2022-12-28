package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	logger := log.Logger{}
	f, err := os.Open("./input.txt")
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	foods := make([][]int, 0)
	currFood := make([]int, 0)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			foods = append(foods, currFood)
			currFood = make([]int, 0)
			continue
		}
		v, err := strconv.Atoi(input)
		if err != nil {
			logger.Fatalf("Failed to parse input: %s", err.Error())
		}
		currFood = append(currFood, v)
	}
	fmt.Printf("got foods for %d elves\n", len(foods))
	maxes := make([]int, 0, 3)
	for _, foodSet := range foods {
		curr := 0
		for _, food := range foodSet {
			curr += food
		}
		if len(maxes) < 3 {
			maxes = append(maxes, curr)
			continue
		}
		lowest, lowestIdx := int(^uint(0)>>1), 0
		for i, max := range maxes {
			if max < lowest {
				lowest, lowestIdx = max, i
			}
		}
		if curr > lowest {
			fmt.Printf("new max: %d\n", curr)
			maxes[lowestIdx] = curr
		}
	}
	sum := 0
	for i := range maxes {
		sum += maxes[i]
	}
	fmt.Printf("Max: %v\n", sum)
}
