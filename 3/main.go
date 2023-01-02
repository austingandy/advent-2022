package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	logger := log.Logger{}
	f, err := os.Open("./input.txt")
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	prio := make(map[string]int)
	t := 0
	for i, l := range "abcdefghijklmnopqrstuvwxyz" {
		s := string(l)
		prio[s] = i + 1
		prio[strings.ToUpper(s)] = prio[s] + 26
	}
	group := make([]string, 0, 3)
	groups := make([][]string, 0)
	for scanner.Scan() {
		input := scanner.Text()
		if len(group) == 3 {
			groups = append(groups, group)
			group = make([]string, 0, 3)
		}
		group = append(group, input)
		r := []rune(input)
		intersects, err := findIntersect(string(r[0:len(r)/2]), string(r[len(r)/2:]))
		if err != nil {
			logger.Fatalf("Failed to find intersection: %s", err.Error())
		}
		if len(intersects) != 1 {
			logger.Fatalf("Found %d intersections, expected 1", len(intersects))
		}
		intersect := intersects[0]
		iPrio, ok := prio[intersect]
		if !ok {
			logger.Fatalf("Failed to find priority for %s", intersect)
		}
		t += iPrio
	}
	groups = append(groups, group)
	fmt.Printf("Total (part 1): %d\n", t)
	t = 0
	for _, g := range groups {
		i, err := findIntersect(g...)
		if err != nil {
			logger.Fatalf("Failed to find intersection: %s", err.Error())
		}
		if len(i) != 1 {
			logger.Fatalf("Found %d intersections, expected 1", len(i))
		}
		intersect := i[0]
		p, ok := prio[intersect]
		if !ok {
			logger.Fatalf("Failed to find priority for %s", intersect)
		}
		t += p
	}
	fmt.Printf("Total (part 2): %d\n", t)

}

func splitItems(l string) (string, string) {
	r := []rune(l)
	return string(r[0 : len(r)/2]), string(r[len(r)/2:])
}

func findIntersect(i ...string) ([]string, error) {
	sets := make([]map[rune]bool, 0)
	if len(i) == 0 {
		return nil, fmt.Errorf("no items to intersect")
	}
	for _, str := range i {
		if len(str) == 0 {
			return nil, fmt.Errorf("empty string")
		}
		s := make(map[rune]bool)
		for _, l := range str {
			s[l] = true
		}
		sets = append(sets, s)
	}
	intersect := sets[0]
	for _, s := range sets[1:] {
		for k := range intersect {
			if !s[k] {
				delete(intersect, k)
			}
		}
	}
	if len(intersect) == 0 {
		return nil, fmt.Errorf("no intersection")
	}
	rtn := make([]string, 0, len(intersect))
	for k := range intersect {
		rtn = append(rtn, string(k))
	}
	return rtn, nil
}
