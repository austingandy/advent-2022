package main

import (
	"bufio"
	"log"
	"os"
	"sort"
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
	root := NewDir(nil)
	curr := root
	for scanner.Scan() {
		tokens, isCommand := ParseLine(scanner.Text())
		if (tokens[0] == "ls" && len(tokens) != 1) || (tokens[0] != "ls" && len(tokens) != 2) {
			logger.Fatalf("malformed input: %s", tokens)
		}
		if isCommand {
			if tokens[0] != "ls" && tokens[0] != "cd" {
				logger.Fatalf("unrecognized command: %s", tokens[0])
			} else if tokens[0] == "cd" {
				if tokens[1] == "/" {
					curr = root
				} else if tokens[1] == ".." {
					curr = curr.parent
					continue
				} else {
					curr = curr.AddOrGetDir(tokens[1])
				}
			}
			continue
		}
		if tokens[0] == "dir" {
			curr.AddOrGetDir(tokens[1])
			continue
		}
		size, err := strconv.ParseInt(tokens[0], 10, 64)
		if err != nil {
			logger.Fatalf("malformed line: %s", tokens)
		}
		curr.files[tokens[1]] = size
	}
	smallDirSize := int64(0)
	for _, d := range root.GetDirsCompSize(func(a int64) bool { return a < int64(100000) }) {
		smallDirSize += d.GetSize()
	}
	logger.Printf("Size (Part 1): %d", smallDirSize)
	candidates := root.GetDirsCompSize(func(a int64) bool { return a > int64(30000000)-(int64(70000000)-root.GetSize()) })
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].GetSize() < candidates[j].GetSize() })
	logger.Printf("Size (Part 2): %d", candidates[0].GetSize())
}

func ParseLine(l string) ([]string, bool) {
	tokens := strings.Split(l, " ")
	if tokens[0] == "$" {
		return tokens[1:], true
	}
	return tokens, false
}

type Dir struct {
	subdirs map[string]*Dir
	files   map[string]int64
	parent  *Dir
	size    *int64
}

func NewDir(parent *Dir) *Dir {
	return &Dir{parent: parent, files: make(map[string]int64), subdirs: make(map[string]*Dir)}
}

func (this *Dir) AddOrGetDir(name string) *Dir {
	if _, ok := this.subdirs[name]; !ok {
		this.subdirs[name] = NewDir(this)
	}
	return this.subdirs[name]
}

func (this *Dir) GetSize() int64 {
	if this.size != nil {
		return *this.size
	}
	totalSize := int64(0)
	for _, size := range this.files {
		totalSize += size
	}
	for _, sd := range this.subdirs {
		totalSize += sd.GetSize()
	}
	this.size = &totalSize
	return totalSize
}

func (this *Dir) GetDirsCompSize(comp func(a int64) bool) (dirs []*Dir) {
	if comp(this.GetSize()) {
		dirs = append(dirs, this)
	}
	for _, sd := range this.subdirs {
		dirs = append(dirs, sd.GetDirsCompSize(comp)...)
	}
	return
}
