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
	topology := make([][]string, 0)
	for scanner.Scan() {
		input := scanner.Text()
		line := make([]string, 0, len(input))
		for _, s := range input {
			line = append(line, string(s))
		}
		topology = append(topology, line)
	}
	hill := NewHill(topology)
	start := hill.Find("S")
	path := hill.aStar(start, "E")
	fmt.Printf("shortest distance (Part 1): %d\n", len(path)-1)
	hill.PrintPath(path)
	c := hill.FindAll("a")
	lengths := make([]int, 0)
	best := len(path)
	for _, v := range c {
		hill.Reset()
		currPath := hill.aStar(v, "E")
		currLen := len(currPath)
		if currPath == nil {
			continue
		}
		want := Pos{Y: 4, X: 0}
		if v.Pos == want {
			hill.PrintPath(currPath)
		}
		fmt.Printf("path length was: %d\n", currLen)
		lengths = append(lengths, currLen)
		if currLen < best {
			best = currLen
		}
	}
	fmt.Printf("shortest distance (Part 2): %d\n", best-1)
}

type Vertex struct {
	value     string
	gScore    float64
	cameFrom  *Vertex
	Pos       Pos
	Elevation int
}

func NewVertex(value string, p Pos, elevation int) *Vertex {
	return &Vertex{
		value:     value,
		gScore:    math.MaxFloat64,
		Pos:       p,
		Elevation: elevation,
	}
}

type Hill [][]*Vertex

func NewHill(top [][]string) Hill {
	elevations := map[string]int{"S": 1, "E": 26}
	for i, l := range "abcdefghijklmnopqrstuvwxyz" {
		elevations[string(l)] = i + 1
	}
	hill := make([][]*Vertex, len(top))
	for i := range top {
		hill[i] = make([]*Vertex, len(top[i]))
		for j := range top[i] {
			val := top[i][j]
			hill[i][j] = NewVertex(val, Pos{Y: i, X: j}, elevations[val])
		}
	}
	return hill
}

func (this Hill) aStar(start *Vertex, goal string) []Pos {
	discovered := []*Vertex{start}
	start.gScore = 0
	curr := start
	for len(discovered) != 0 {
		sort.Slice(discovered, func(i, j int) bool {
			return discovered[i].fScore(curr.Pos) < discovered[j].fScore(curr.Pos)
		})
		curr, discovered = discovered[0], discovered[1:]
		if curr.value == goal {
			p := buildPath(curr)
			this.Reset()
			return p
		}
		for _, n := range this.Neighbors(curr.Pos) {
			score := curr.gScore + curr.Pos.Euclid(n.Pos)
			if score < n.gScore {
				n.cameFrom = curr
				n.gScore = score
				isDiscovered := false
				for _, d := range discovered {
					if d == n {
						isDiscovered = true
						break
					}
				}
				if !isDiscovered {
					discovered = append(discovered, n)
				}
			}
		}
	}
	return nil
}

func (this *Vertex) fScore(curr Pos) float64 {
	if this.gScore == math.MaxFloat64 {
		return this.gScore
	}
	return this.gScore + this.Pos.Euclid(curr)
}

func (this Hill) Neighbors(p Pos) []*Vertex {
	n := make([]*Vertex, 0, len(dirs))
	for _, d := range dirs {
		next := d(p)
		if !this.isValid(next) || this.Get(next).Elevation-this.Get(p).Elevation > 1 {
			continue
		}
		n = append(n, this.Get(next))
	}
	return n
}

func (this Hill) Find(val string) *Vertex {
	for _, v := range this.GetVertices() {
		if v.value == val {
			return v
		}
	}
	return nil
}

func (this Hill) FindAll(val string) []*Vertex {
	vs := make([]*Vertex, 0)
	for _, v := range this.GetVertices() {
		if v.value == val {
			vs = append(vs, v)
		}
	}
	return vs
}

func (this Hill) GetVertices() []*Vertex {
	vs := make([]*Vertex, 0)
	for i := range this {
		for j := range this[i] {
			vs = append(vs, this[i][j])
		}
	}
	return vs
}

func (this Hill) Get(p Pos) *Vertex {
	return this[p.Y][p.X]
}

func (this Hill) PrintPath(path []Pos) {
	posMap := make(map[Pos]string)
	for i := 1; i < len(path); i++ {
		last, curr := path[i-1], path[i]
		switch curr {
		case down(last):
			posMap[last] = "v"
		case up(last):
			posMap[last] = "^"
		case left(last):
			posMap[last] = "<"
		case right(last):
			posMap[last] = ">"
		default:
			posMap[last] = "?"
		}
	}
	for i, v := range this.GetVertices() {
		if c, ok := posMap[v.Pos]; ok {
			fmt.Printf(c)
		} else {
			fmt.Printf(v.value)
		}
		if i%len(this[0]) == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Println()
}

func (this Hill) Reset() {
	for _, v := range this.GetVertices() {
		v.cameFrom = nil
		v.gScore = math.MaxFloat64
	}
}

func (this Hill) isValid(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && p.Y < len(this) && p.X < len(this[p.Y])
}

func buildPath(to *Vertex) []Pos {
	path := []Pos{to.Pos}
	for v := to.cameFrom; v != nil; v = v.cameFrom {
		path = append([]Pos{v.Pos}, path...)
	}
	return path
}

type Pos struct {
	X int
	Y int
}

func (this Pos) Euclid(that Pos) float64 {
	xDist := this.X - that.X
	yDist := this.Y - that.Y
	return math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))
}

type Direction func(p Pos) Pos

func up(p Pos) Pos { return Pos{Y: p.Y - 1, X: p.X} }

func down(p Pos) Pos { return Pos{Y: p.Y + 1, X: p.X} }

func left(p Pos) Pos { return Pos{Y: p.Y, X: p.X - 1} }

func right(p Pos) Pos { return Pos{Y: p.Y, X: p.X + 1} }

var dirs = []Direction{down, right, up, left}
