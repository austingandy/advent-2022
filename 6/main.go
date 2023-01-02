package main

import (
	"bufio"
	"log"
	"os"
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
	input := ""
	for scanner.Scan() {
		input += scanner.Text()
	}
	packetMarker := NewLastN(4)
	messageMarker := NewLastN(14)
	for i, l := range input {
		if !packetMarker.foundUnique {
			if packetMarker.IsUnique() {
				logger.Printf("Packets processed (Part 1): %d", i)
			} else {
				packetMarker.Append(string(l))
			}
		}
		if !messageMarker.foundUnique {
			if messageMarker.IsUnique() {
				logger.Printf("Messages processed (Part 2): %d", i)
			} else {
				messageMarker.Append(string(l))
			}
		}
		if packetMarker.foundUnique && messageMarker.foundUnique {
			break
		}
	}
}

type LastN struct {
	vals        []string
	max         int
	foundUnique bool
}

func NewLastN(max int) *LastN {
	return &LastN{
		vals:        make([]string, 0, max),
		max:         max,
		foundUnique: false,
	}
}

func (this *LastN) Append(r string) {
	if len(this.vals) < this.max {
		this.vals = append(this.vals, r)
		return
	}
	this.vals = append(this.vals[1:], r)
}

func (this *LastN) IsUnique() bool {
	if len(this.vals) < this.max {
		return false
	}
	vals := make(map[string]bool)
	for _, v := range this.vals {
		if vals[v] {
			return false
		}
		vals[v] = true
	}
	this.foundUnique = true
	return true
}

func (this *LastN) String() string {
	s := ""
	for _, v := range this.vals {
		s += v
	}
	return s
}
