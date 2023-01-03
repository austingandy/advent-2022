package main

import (
	"bufio"
	"log"
	"math"
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
	v1, v2, total1, total2 := map[Pos]bool{Pos{}: true}, map[Pos]bool{Pos{}: true}, 1, 1
	h1, t1 := NewRope(2, Pos{})
	h2, t2 := NewRope(10, Pos{})
	for scanner.Scan() {
		input := strings.Split(scanner.Text(), " ")
		dir := commandToDirections[input[0]]
		dist, _ := strconv.Atoi(input[1])
		for i := 0; i < dist; i += 1 {
			h1.MoveHead(dir)
			if !v1[t1.Pos] {
				total1 += 1
				v1[t1.Pos] = true
			}
			h2.MoveHead(dir)
			if !v2[t2.Pos] {
				total2 += 1
				v2[t2.Pos] = true
			}
		}
	}
	logger.Printf("Total visited (Part 1): %d", total1)
	logger.Printf("Total visited (Part 2): %d", total2)
}

type Direction func(p Pos) Pos
type Pos struct{ X, Y int }
type Rope struct {
	Pos  Pos
	Tail *Rope
}

func NewRope(n int, p Pos) (*Rope, *Rope) {
	head := &Rope{Pos: p}
	tail := head
	for i := 0; i < n-1; i += 1 {
		tail.Tail = &Rope{Pos: p}
		tail = tail.Tail
	}
	return head, tail
}

func (r *Rope) MoveHead(d Direction) {
	r.Pos = d(r.Pos)
	if r.Tail != nil && !isAdjacent(r.Pos, r.Tail.Pos) {
		r.Tail.MoveHead(getTailDir(r.Pos, r.Tail.Pos))
	}
}

func north(p Pos) Pos     { return Pos{p.X, p.Y + 1} }
func south(p Pos) Pos     { return Pos{p.X, p.Y - 1} }
func east(p Pos) Pos      { return Pos{p.X + 1, p.Y} }
func west(p Pos) Pos      { return Pos{p.X - 1, p.Y} }
func southwest(p Pos) Pos { return south(west(p)) }
func southeast(p Pos) Pos { return south(east(p)) }
func northwest(p Pos) Pos { return north(west(p)) }
func northeast(p Pos) Pos { return north(east(p)) }

var directions = []Direction{north, south, east, west, southwest, southeast, northeast, northwest}
var commandToDirections = map[string]Direction{"R": east, "D": south, "U": north, "L": west}

func isAdjacent(p1, p2 Pos) bool {
	return math.Abs(float64(p1.X)-float64(p2.X)) <= float64(1) && math.Abs(float64(p1.Y)-float64(p2.Y)) <= float64(1)
}

func isDiag(p1, p2 Pos) bool { return p1.X != p2.X && p1.Y != p2.Y }

func getTailDir(head, tail Pos) Direction {
	if isDiag(head, tail) {
		if head.X > tail.X {
			if head.Y > tail.Y {
				return northeast
			}
			return southeast
		}
		if head.Y > tail.Y {
			return northwest
		}
		return southwest
	}
	if head.X > tail.X {
		return east
	}
	if head.Y > tail.Y {
		return north
	}
	if head.X < tail.X {
		return west
	}
	return south
}
