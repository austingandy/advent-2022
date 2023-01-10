package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		t := tokenize(scanner.Text())
		fmt.Println(t)

	}
}

func processLine(l []string) *Container {
	if len(l) == 0 {
		return &Container{Type: Empty}
	}
	switch l[0] {
	case "[":
		bc, i := 1, 0
		for bc != 0 && i < len(l) {
			i += 1
			switch l[i] {
			case "[":
				bc += 1
			case "]":
				bc -= 1
			default:
				continue
			}
		}
		if i-0 <= 1 {
			// TODO: should this actually be nil? idk
			return nil
		} else {
			return processLine(l[1 : i-1]).BuildArray(processLine(l[i+1:]))
		}
	case "]":
		fmt.Printf("hi")
	// okay we probably have numbers now
	default:
		vals := []
	}
	return nil
}


func processValues(vs []int) []*Container {

}

func tokenize(line string) []string {
	digits := make(map[string]bool)
	for i := range [10]int{} {
		digits[fmt.Sprintf("%d", i)] = true
	}
	l := strings.Split(line, "")
	filteredL := make([]string, 0, len(l))
	for i := 0; i < len(l); i += 1 {
		c := l[i]
		if c == "," {
			continue
		}
		if digits[c] {
			i += 1
			for i < len(l) && digits[l[i]] {
				c += l[i]
				i += 1
			}
			i -= 1
		}
		filteredL = append(filteredL, c)
	}
	return filteredL
}

func createPacket() *Container {
	return &Container{
		Type: Array,
		ArrayData: []*Container{
			{
				Type: Array,
				ArrayData: []*Container{
					{
						Type:    Int,
						IntData: 4,
					},
					{
						Type:    Int,
						IntData: 4,
					},
				},
			},
			{
				Type:    Int,
				IntData: 4,
			},
			{
				Type:    Int,
				IntData: 4,
			},
		},
	}
}

type Container struct {
	ArrayData []*Container
	IntData   int
	IntsData []*Container
	Type      Type
}

type Type int

const (
	Array Type = iota
	Int
	Ints
	Empty
)

func (this *Container) BuildArray(o *Container) *Container {
	if this.Type == Ints {
		switch o.Type {
		case Ints:
			vals := make([]*Container, 0)
			for _, v := range this.IntsData {
				vals = append(vals, v)
			}
			for _, v := range this.IntsData {
				vals = append(vals, v)
			}
			return &Container{Type: Array, ArrayData: vals}
		}
	}
	return &Container{
		Type: Array,
		ArrayData: []*Container{this, o},
	}
}
