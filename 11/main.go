package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Panicf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	inputs, currInput := make([][]string, 0), make([]string, 0)
	for scanner.Scan() {
		in := strings.TrimSpace(scanner.Text())
		if len(in) == 0 {
			inputs = append(inputs, currInput)
			currInput = make([]string, 0, len(currInput))
			continue
		}
		currInput = append(currInput, in)
	}
	inputs = append(inputs, currInput)
	monkeys := make([]*Monkey, 0, len(inputs))
	for range inputs {
		monkeys = append(monkeys, &Monkey{})
	}
	for i, input := range inputs {
		// already dealt with first line, so just ignore it for now
		for _, line := range input[1:] {
			switch tokens := strings.Split(line, ": "); tokens[0] {
			case "Starting items":
				for _, item := range strings.Split(tokens[1], ", ") {
					v, _ := strconv.Atoi(item)
					monkeys[i].items = append(monkeys[i].items, v)
				}
			case "Operation":
				equation := strings.Split(strings.Trim(tokens[1], "new = "), " ")
				update := func(old int) int {
					a, b := getVal(equation[0], old), getVal(equation[2], old)
					if equation[1] == "+" {
						return a + b
					}
					return a * b
				}
				monkeys[i].update = update
			case "Test":
				divisibleBy, _ := strconv.Atoi(strings.Trim(tokens[1], "divisible by "))
				monkeys[i].test = divisibleBy
			case "If true":
				index, _ := strconv.Atoi(strings.Trim(tokens[1], "throw to monkey "))
				monkeys[i].sourceTrue = monkeys[index]
			case "If false":
				index, _ := strconv.Atoi(strings.Trim(tokens[1], "throw to monkey "))
				monkeys[i].sourceFalse = monkeys[index]
			default:
				log.Panicf("unrecognized starting token: %s", tokens[0])
			}
		}
	}
	passed := make([]bool, len(monkeys))
	threshold := 1
	for _, m := range monkeys {
		threshold *= m.test
	}
	for i := 0; i < 10000; i += 1 {
		for j := range monkeys {
			passed[j] = monkeys[j].takeTurn(threshold) || passed[j]
		}
		k := i + 1
		if k == 1 || k == 20 || k > 20 && k%1000 == 0 {
			fmt.Printf("== After round %d ==\n", k)
			for l := range monkeys {
				fmt.Printf("Monkey %d inspected items %d times\n", l, monkeys[l].inspectionCount)
			}
			fmt.Println()
			for l := range passed {
				if !passed[l] {
					fmt.Printf("monkey %d still hasn't passed its tes\n", l)
				}
			}
			fmt.Println()
		}

	}
	sort.Slice(monkeys, func(i, j int) bool { return monkeys[i].inspectionCount > monkeys[j].inspectionCount })
	fmt.Printf("Total inspections of top 2 (Part 1): %d\n", monkeys[0].inspectionCount*monkeys[1].inspectionCount)
	for _, m := range monkeys {
		fmt.Printf("%d\n", m.inspectionCount)
	}
	fmt.Printf("maxint64: %d\n", math.MaxInt64)
}

func getVal(operand string, old int) int {
	if operand == "old" {
		return old
	}
	n, _ := strconv.Atoi(operand)
	return n
}

type Monkey struct {
	items           []int
	update          func(old int) int
	test            int
	sourceTrue      *Monkey
	sourceFalse     *Monkey
	inspectionCount int64
}

func (this *Monkey) takeTurn(threshold int) bool {
	passedTest := false
	for len(this.items) > 0 {
		item := this.update(this.items[0])
		if item > threshold {
			item = item % threshold
		}
		this.items = this.items[1:]
		if item%this.test == 0 {
			passedTest = true
			this.sourceTrue.items = append(this.sourceTrue.items, item)
		} else {
			this.sourceFalse.items = append(this.sourceFalse.items, item)
		}
		this.inspectionCount += 1
	}
	return passedTest
}
