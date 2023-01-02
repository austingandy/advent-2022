package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
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
	defer func() { _ = f.Close() }()
	s := make([]string, 0)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			break
		}
		s = append(s, input)
	}
	logger.Printf("got %d lines\n", len(s))
	i := 1
	stacks := make([][]string, 0)
	for i < len(s[0]) {
		stack := make([]string, 0)
		// don't add the index row at the bottom of the diagram
		for _, row := range s[:len(s)-1] {
			c := string(row[i])
			if c != " " {
				// reverse ordering to make popping easier
				stack = append([]string{c}, stack...)
			}
		}
		stacks = append(stacks, stack)
		i += 4
	}
	isDigit := regexp.MustCompile(`\d`)

	moves := make([]Move, 0)
	for scanner.Scan() {
		input := scanner.Text()
		digits := make([]int, 0, 3)
		i = 0
		d := ""
		for i = 0; i < len(input); i += 1 {
			l := string(input[i])
			if !isDigit.MatchString(l) {
				if len(d) > 0 {
					v, err := strconv.Atoi(d)
					if err != nil {
						logger.Fatalf("Failed to parse digit: %s", err.Error())
					}
					digits = append(digits, v)
					d = ""
				}
				continue
			}
			d += l
		}
		if len(d) > 0 {
			v, err := strconv.Atoi(d)
			if err != nil {
				logger.Fatalf("Failed to parse digit: %s", err.Error())
			}
			digits = append(digits, v)
		}
		if len(digits) != 3 {
			logger.Fatalf("Expected 3 digits in line %s, got %d", input, len(digits))
		}
		moves = append(moves, Move{
			qty:  digits[0],
			from: digits[1] - 1,
			to:   digits[2] - 1,
		})
	}
	stackV2 := make([][]string, 0, len(stacks))
	for _, stack := range stacks {
		newStack := make([]string, 0, len(stack))
		for _, v := range stack {
			newStack = append(newStack, v)
		}
		stackV2 = append(stackV2, newStack)
	}
	for _, move := range moves {
		stacks = makeMove(stacks, move)
	}
	top := ""

	for _, stack := range stacks {
		top += stack[len(stack)-1]
	}

	logger.Printf("Top (Part 1): %s", top)
	for _, move := range moves {
		stackV2 = makeMoveV2(stackV2, move)
	}
	top = ""
	for _, stack := range stackV2 {
		top += stack[len(stack)-1]
	}
	logger.Printf("Top (Part 2): %s", top)
}

type Move struct {
	qty  int
	from int
	to   int
}

func makeMove(stack [][]string, move Move) [][]string {
	for i := 0; i < move.qty; i++ {
		stack[move.from], stack[move.to] = pop(stack[move.from], stack[move.to])
	}
	return stack
}

func makeMoveV2(stack [][]string, move Move) [][]string {
	from := stack[move.from]
	fromIndex := len(stack[move.from]) - move.qty
	stack[move.to] = append(stack[move.to], from[fromIndex:]...)
	stack[move.from] = from[:fromIndex]
	return stack
}

func pop(from []string, to []string) ([]string, []string) {
	to = append(to, from[len(from)-1])
	from = from[:len(from)-1]
	return from, to
}
