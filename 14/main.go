package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func getPos(s string, i int) (Pos, int) {
	x, i := getInt(s, i)
	// skip the comma -> i+1
	y, i := getInt(s, i+1)
	return Pos{X: x, Y: y}, i
}

func getInt(s string, i int) (int, int) {
	var j int
	for j = i; j < len(s); j += 1 {
		if s[j] >= '0' && s[j] <= '9' {
			continue
		}
		break
	}
	val, _ := strconv.Atoi(s[i:j])
	return val, j
}

func (p Pos) LeftOf(o Pos) bool {
	return p.X < o.X
}

func (p Pos) Above(o Pos) bool {
	return p.Y < o.Y
}

func main() {
	logger := log.Logger{}
	f, err := os.Open("./input.txt")
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	rocks := make([][]Line, 0)
	for scanner.Scan() {
		s := scanner.Text()
		points := make([]Pos, 0)
		i := 0
		for i < len(s) {
			if !(s[i] >= '0' && s[i] <= '9') {
				i += 1
				continue
			}
			var p Pos
			p, i = getPos(s, i)
			points = append(points, p)
		}
		rock := make([]Line, 0)
		for i = 1; i < len(points); i += 1 {
			rock = append(rock, NewLine(points[i-1], points[i]))
		}
		rocks = append(rocks, rock)
	}
	maxX, maxY := 0, 0
	for _, rock := range rocks {
		for _, line := range rock {
			if line.end.X > maxX {
				maxX = line.end.X
			}
			if line.end.Y > maxY {
				maxY = line.end.Y
			}
			if line.start.X > maxX {
				maxX = line.start.X
			}
			if line.start.Y > maxY {
				maxY = line.start.Y
			}
		}
	}
	grid := NewGrid(maxX, maxY)
	for _, r := range rocks {
		grid.PlaceRock(r)
	}
	fmt.Printf("Part 1: %d\n", grid.DropAllSand())
	grid = NewGrid(maxX, maxY)
	for _, r := range rocks {
		grid.PlaceRock(r)
	}
	grid.SetIsBounded(false)
	fmt.Printf("Part 2: %d\n", grid.DropAllSand())
	fmt.Printf(grid.String())
}

func NewLine(p, q Pos) Line {
	if p.LeftOf(q) || p.Above(q) {
		return Line{p, q}
	}
	return Line{q, p}
}

func (g *Grid) String() string {
	s := ""
	maxY := g.maxPos.Y
	if !g.isBounded {
		maxY += 3
	}
	for y := 0; y < maxY; y += 1 {
		for x := 200; x < g.maxPos.X; x += 1 {
			switch g.Get(Pos{X: x, Y: y}).(type) {
			case Rock:
				s += "#"
			case Air:
				s += "."
			case Sand:
				s += "o"
			}
		}
		s += "\n"
	}
	return s
}

type Grid struct {
	g                      [][]Element
	newG                   map[Pos]Element
	maxPos, sandEntryPoint Pos
	isBounded              bool
}

func (g *Grid) Set(p Pos, e Element) {
	g.newG[p] = e
}

func (g *Grid) Get(p Pos) Element {
	if e, ok := g.newG[p]; ok {
		return e
	}
	if !g.isBounded && p.Y == g.maxPos.Y+2 {
		return Rock{}
	}
	return Air{}
}

func (g *Grid) IsInBounds(p Pos) bool {
	return p.X >= 0 && p.X <= g.maxPos.X && p.Y >= 0 && p.Y <= g.maxPos.Y
}

func (g *Grid) IsAir(p Pos) bool {
	return g.Get(p) == Air{}
}

func (g *Grid) PlaceRock(r []Line) {
	for _, line := range r {
		if line.start.X != line.end.X {
			for p := line.start; !line.end.LeftOf(p); p = right(p) {
				g.Set(p, Rock{})
			}
		}
		if line.start.Y != line.end.Y {
			for p := line.start; !line.end.Above(p); p = down(p) {
				if !g.IsInBounds(p) {
					fmt.Printf("Out of bounds: %d, %d\n", p.X, p.Y)
					break
				}
				g.Set(p, Rock{})
			}
		}
	}
}

func (g *Grid) MaxPos() Pos {
	maxPos := Pos{}
	for p := range g.newG {
		if maxPos.Above(p) {
			maxPos = Pos{X: maxPos.X, Y: p.Y}
		}
		if maxPos.LeftOf(p) {
			maxPos = Pos{X: p.X, Y: maxPos.Y}
		}
	}
	return maxPos
}

// DropAllSand drops grains of sand until one overflows or the outlet is blocked
func (g *Grid) DropAllSand() int {
	fmt.Printf("len(rocks) is %d\n", len(g.newG))
	t := 0
	for g.IsAir(g.sandEntryPoint) {
		// if !g.isBounded && t%1000 == 0 {
		// 	fmt.Printf(g.String())
		// }
		if !g.DropSand() {
			break
		}
		t += 1
	}
	return t
}

// DropSand drops a single grain of sand until it is at rest
func (g *Grid) DropSand() bool {
	p := g.sandEntryPoint
	directions := []func(p Pos) Pos{
		down,
		func(p Pos) Pos { return down(left(p)) },
		func(p Pos) Pos { return down(right(p)) },
	}
	for {
		newP, stillInBounds := func() (Pos, bool) {
			for _, d := range directions {
				q := d(p)
				if g.isBounded && !g.IsInBounds(q) {
					return Pos{}, false
				}
				if g.IsAir(q) {
					return q, true
				}
			}
			return p, true
		}()
		if !stillInBounds {
			return false
		}
		if newP == p {
			break
		}
		p = newP
	}
	g.Set(p, Sand{atRest: true})
	return true
}

func (g *Grid) SetIsBounded(b bool) {
	g.isBounded = b
}

func NewGrid(maxX, maxY int) Grid {
	return Grid{
		sandEntryPoint: Pos{500, 0},
		newG:           make(map[Pos]Element),
		maxPos:         Pos{X: maxX, Y: maxY},
		isBounded:      true,
	}
}

type Element interface {
	isAir() bool
	isAtRest() bool
}

type Sand struct {
	atRest bool
}

func (s Sand) isAir() bool { return false }

func (s Sand) isAtRest() bool { return s.atRest }

type Rock struct{}

func (r Rock) isAir() bool { return false }

func (r Rock) isAtRest() bool { return true }

type Air struct{}

func (a Air) isAir() bool { return true }

func (a Air) isAtRest() bool { return false }

type Pos struct {
	X int
	Y int
}

func down(p Pos) Pos { return Pos{Y: p.Y + 1, X: p.X} }

func up(p Pos) Pos { return Pos{Y: p.Y - 1, X: p.X} }

func left(p Pos) Pos { return Pos{Y: p.Y, X: p.X - 1} }

func right(p Pos) Pos { return Pos{Y: p.Y, X: p.X + 1} }

type Line struct {
	start Pos
	end   Pos
}
