package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	i, s := 1, 0
	toCompare := make([]*NestedInteger, 0)
	allPackets := make([]*NestedInteger, 0)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			continue
		}
		packet := deserialize(scanner.Text())
		toCompare = append(toCompare, packet)
		allPackets = append(allPackets, packet)
		if len(toCompare) == 2 {
			if toCompare[0].Less(toCompare[1]) == 1 {
				s += i
			}
			toCompare = nil
			i += 1
		}
	}
	fmt.Printf("Total (Part 1): %d\n", s)
	d2, d6 := deserialize("[[2]]"), deserialize("[[6]]")
	allPackets = append(allPackets, d2, d6)
	sort.Slice(allPackets, func(i, j int) bool {
		return allPackets[i].Less(allPackets[j]) != -1
	})
	i2, i6 := 0, 0
	for i, p := range allPackets {
		if p.Equals(d2) {
			i2 = i + 1
		}
		if p.Equals(d6) {
			i6 = i + 1
		}
		if i2 != 0 && i6 != 0 {
			break
		}
	}
	fmt.Printf("Part 2: %d\n", i2*i6)
}

/**
 * // This is the interface that allows for creating nested lists.
 * // You should not implement it, or speculate about its implementation
 * type NestedInteger struct {
 * }
 *
 * // Return true if this NestedInteger holds a single integer, rather than a nested list.
 * func (n NestedInteger) IsInteger() bool {}
 *
 * // Return the single integer that this NestedInteger holds, if it holds a single integer
 * // The result is undefined if this NestedInteger holds a nested list
 * // So before calling this method, you should have a check
 * func (n NestedInteger) GetInteger() int {}
 *
 * // Set this NestedInteger to hold a single integer.
 * func (n *NestedInteger) SetInteger(value int) {}
 *
 * // Set this NestedInteger to hold a nested list and adds a nested integer to it.
 * func (n *NestedInteger) Add(elem NestedInteger) {}
 *
 * // Return the nested list that this NestedInteger holds, if it holds a nested list
 * // The list length is zero if this NestedInteger holds a single integer
 * // You can access NestedInteger's List element directly if you want to modify it
 *
 */

type NestedInteger struct {
	val  *int
	List []*NestedInteger
}

func (n *NestedInteger) IsInteger() bool {
	return n.val != nil
}

func (n *NestedInteger) Equals(o *NestedInteger) bool {
	return n.Less(o) == 0
}

func (n *NestedInteger) GetInteger() int {
	if n.IsInteger() {
		return *n.val
	}
	return 0
}

func (n *NestedInteger) GetList() []*NestedInteger {
	if n.IsInteger() {
		return nil
	}
	return n.List
}

func (n *NestedInteger) Add(elem NestedInteger) {
	n.List = append(n.List, &elem)
}

func (n *NestedInteger) SetInteger(value int) {
	n.List = nil
	n.val = &value
}

func (n *NestedInteger) Less(o *NestedInteger) int {
	if n.IsInteger() && o.IsInteger() {
		if n.GetInteger() < o.GetInteger() {
			return 1
		}
		if n.GetInteger() > o.GetInteger() {
			return -1
		}
		return 0
	}
	if n.IsInteger() {
		wrapper := &NestedInteger{}
		wrapper.Add(*n)
		return wrapper.Less(o)
	}
	if o.IsInteger() {
		wrapper := &NestedInteger{}
		wrapper.Add(*o)
		return n.Less(wrapper)
	}
	for i := 0; i < len(n.List); i += 1 {
		if i > len(o.List)-1 {
			return -1
		}
		switch res := n.List[i].Less(o.List[i]); res {
		case 1, -1:
			return res
		default:
			continue
		}
	}
	if len(n.List) == len(o.List) {
		return 0
	}
	return 1
}

func deserialize(s string) *NestedInteger {
	ni := &NestedInteger{}
	if s[0] != '[' {
		v, _ := strconv.Atoi(s)
		ni.SetInteger(v)
		return ni
	}
	for i := 1; i < len(s); {
		if s[i] == ',' || s[i] == ']' {
			i += 1
			continue
		}
		if s[i] != '[' {
			elem := &NestedInteger{}
			start := i
			for j := i + 1; j < len(s); j += 1 {
				if !(s[j] >= '0' && s[j] <= '9') {
					i = j
					break
				}
			}
			v, _ := strconv.Atoi(s[start:i])
			elem.SetInteger(v)
			ni.Add(*elem)
			continue
		}
		l := 1
		start := i
		for j := i + 1; j < len(s); j += 1 {
			if s[j] == '[' {
				l += 1
			} else if s[j] == ']' {
				l -= 1
				if l == 0 {
					i = j
					break
				}
			}
		}
		elem := deserialize(s[start : i+1])
		ni.Add(*elem)
	}
	return ni
}
