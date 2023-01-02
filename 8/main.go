package main

import (
	"bufio"
	"log"
	"os"
	"sort"
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
	rows := Grid{}
	for scanner.Scan() {
		row := make([]int, 0)
		for _, t := range scanner.Text() {
			v, _ := strconv.Atoi(string(t))
			row = append(row, v)
		}
		rows = append(rows, row)
	}
	visible, scores := 0, make([]int, 0, len(rows)*len(rows[0]))
	for i := range rows {
		for j := range rows[i] {
			if !rows.isHidden(i, j) {
				visible += 1
			}
			scores = append(scores, rows.scenicScore(i, j))
		}
	}
	logger.Printf("Total visible (Part 1): %d", visible)
	sort.Slice(scores, func(i, j int) bool { return scores[i] > scores[j] })
	logger.Printf("Max scenic score (Part 2): %d", scores[0])
}

type Grid [][]int
type Direction func(i, j int) (int, int)
type ReturnIndicator[T any] func(target, curr int, old T, grid Grid) (T, bool)

func (this Grid) isHidden(i, j int) bool {
	for _, d := range []Direction{north, south, east, west} {
		if !this.isHiddenInDir(i, j, d) {
			return false
		}
	}
	return true
}

func (this Grid) scenicScore(i, j int) int {
	dist := 1
	for _, d := range []Direction{north, south, east, west} {
		dist *= this.scoreInDir(i, j, d)
	}
	return dist
}

func (this Grid) isHiddenInDir(i, j int, d Direction) bool {
	return genericInDir[bool](i, j, d, this, func(target, curr int, _ bool, grid Grid) (bool, bool) {
		if target <= curr {
			return true, true
		}
		return false, false
	})
}

func genericInDir[T any](i, j int, d Direction, grid Grid, ri ReturnIndicator[T]) (v T) {
	for k, l := d(i, j); k >= 0 && l >= 0 && k < len(grid) && l < len(grid[k]); k, l = d(k, l) {
		newVal, shouldReturn := ri(grid[i][j], grid[k][l], v, grid)
		if shouldReturn {
			return newVal
		}
		v = newVal
	}
	return v
}

func (this Grid) scoreInDir(i, j int, d Direction) int {
	return genericInDir[int](i, j, d, this, func(target, curr, old int, grid Grid) (int, bool) { return old + 1, curr >= target })
}

func north(i, j int) (int, int) { return i - 1, j }

func south(i, j int) (int, int) { return i + 1, j }

func east(i, j int) (int, int) { return i, j + 1 }

func west(i, j int) (int, int) { return i, j - 1 }
