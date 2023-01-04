package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	hill := NewHill()
	for scanner.Scan() {
		input := scanner.Text()
		line := make([]string, 0, len(input))
		for _, s := range input {
			line = append(line, string(s))
		}
		hill.Topology = append(hill.Topology, line)
	}
	s, err := hill.StartingPos()
	if err != nil {
		log.Fatalf("Error retreiving starting position: %s", err.Error())
	}
	fmt.Printf("Starting position is X: %d, Y: %d\n", s.X, s.Y)
	l, path, err := hill.Climb()
	if err != nil {
		log.Fatalf("fell down the hill: %s", err.Error())
	}
	fmt.Printf("shortest distance was %d\n", l)
	fmt.Println("path was:")
	for _, p := range path {
		fmt.Printf("(%d,%d)\n", p.X, p.Y)
	}
}

type Hill struct {
	Topology   [][]string
	Elevations map[string]int
	visited    map[Pos]bool
}

func NewHill() Hill {
	elevations := map[string]int{"S": 1, "E": 26}
	for i, l := range "abcdefghijklmnopqrstuvwxyz" {
		elevations[string(l)] = i + 1
	}
	return Hill{
		Elevations: elevations,
		visited:    make(map[Pos]bool),
	}
}

func (t Hill) Climb() (int, []Pos, error) {
	t.visited = make(map[Pos]bool)
	p, err := t.StartingPos()
	if err != nil {
		return 0, nil, fmt.Errorf("error retreiving starting position: %s", err.Error())
	}
	var path []Pos
	depth, l := 10, math.MaxInt
	for depth < l {
		fmt.Printf("traversing for depth %d. currrent length is %d\n", depth, l)
		l, path = t.visit(p, []Pos{}, depth)
		depth += 1
	}
	return l, path, nil
}

func (t Hill) PrintGrid() {
	rows, row := make([]string, 0), ""
	for _, pos := range t.Positions() {
		if t.visited[pos] {
			row += "x"
		} else {
			row += t.charAt(pos)
		}
		if len(row) >= len(t.Topology) {
			rows = append(rows, row)
			row = ""
		}
	}
	rows = append(rows, row)
	for _, roww := range rows {
		fmt.Println(roww)
	}
}

func (t Hill) visit(p Pos, path []Pos, maxDepth int) (int, []Pos) {
	if len(path) >= maxDepth {
		//fmt.Printf("[depth %d] exceeded max depth of %d at position (%d, %d)\n", maxDepth, len(path), p.X, p.Y)
		return math.MaxInt, path
	}

	newPath := make([]Pos, len(path))
	for i := range path {
		newPath[i] = path[i]
	}
	t.visited[p] = true
	newPath = append(newPath, p)
	//if len(newPath)%10 == 0 {
	//	fmt.Printf("[depth %d] maxDepth is %d\n", len(newPath), maxDepth)
	//}
	//fmt.Printf("[depth %d] visiting (%d, %d). current path is:\n", len(newPath), p.X, p.Y)
	//t.PrintGrid()
	if t.charAt(p) == "E" {
		fmt.Printf("[depth %d] Made it!!\n", len(newPath))
		return 0, newPath
	}
	candidates := make([]Pos, 0, len(dirs))
	for _, d := range dirs {
		next := d(p)
		if !t.isValid(next) {
			//fmt.Printf("\t[depth %d] position (%d, %d) not valid in grid\n", len(newPath), next.X, next.Y)
			continue
		}
		if t.visited[next] {
			//fmt.Printf("\t[depth %d] already visited (%d, %d)\n", len(newPath), next.X, next.Y)
			continue
		}
		if t.Elevation(next)-t.Elevation(p) > 1 {
			//fmt.Printf("\t[depth %d] elevation of position (%d, %d) too high from position (%d, %d)\n", len(newPath), next.X, next.Y, p.X, p.Y)
			continue
		}
		candidates = append(candidates, next)
	}
	if len(candidates) == 0 {
		return math.MaxInt, newPath
	}
	sort.Slice(candidates, func(i, j int) bool { return t.Elevation(candidates[i]) > t.Elevation(candidates[j]) })
	var minPath []Pos
	min := math.MaxInt
	for _, candidate := range candidates {
		if min == 0 {
			break
		}
		var newMaxDepth int
		if min > maxDepth {
			newMaxDepth = maxDepth
		} else {
			newMaxDepth = min + len(newPath)
		}
		length, cPath := t.visit(candidate, newPath, newMaxDepth)
		if length < min {
			min, minPath = length, cPath
		}
		t.visited[candidate] = false
	}
	if min == math.MaxInt {
		return min, newPath
	}
	return 1 + min, minPath
}

func (t Hill) isValid(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && p.Y < len(t.Topology) && p.X < len(t.Topology[p.Y])
}

func (t Hill) Elevation(p Pos) int {
	return t.Elevations[t.Topology[p.Y][p.X]]
}

func (t Hill) charAt(p Pos) string {
	return t.Topology[p.Y][p.X]
}

func (t Hill) Positions() []Pos {
	positions := make([]Pos, 0, len(t.Topology)*len(t.Topology[0]))
	for y := range t.Topology {
		for x := range t.Topology[y] {
			positions = append(positions, Pos{Y: y, X: x})
		}
	}
	return positions
}

func (t Hill) StartingPos() (Pos, error) {
	for _, p := range t.Positions() {
		if t.charAt(p) == "S" {
			return p, nil
		}
	}
	return Pos{}, fmt.Errorf("no starting position found")
}

type Pos struct {
	X int
	Y int
}

type Direction func(p Pos) Pos

func up(p Pos) Pos {
	return Pos{Y: p.Y - 1, X: p.X}
}

func down(p Pos) Pos {
	return Pos{Y: p.Y + 1, X: p.X}
}

func left(p Pos) Pos {
	return Pos{Y: p.Y, X: p.X - 1}
}

func right(p Pos) Pos {
	return Pos{Y: p.Y, X: p.X + 1}
}

var dirs = []Direction{down, right, up, left}
