package main

import (
	"bufio"
	"log"
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
	rootDir := NewDir("/", nil)
	currDir := rootDir
	for scanner.Scan() {
		tokens, isCommand := ParseLine(scanner.Text())
		if (tokens[0] == "ls" && len(tokens) != 1) || (tokens[0] != "ls" && len(tokens) != 2) {
			logger.Fatalf("malformed input: %s", tokens)
		}
		if isCommand {
			switch tokens[0] {
			case "ls":
				continue
			case "cd":
				dest := tokens[1]
				if dest == "/" {
					currDir = rootDir
					continue
				}
				if dest == ".." {
					parent := currDir.parent
					if parent == nil {
						logger.Fatalf("no parent directory exists for dir: %v", currDir)
					}
					currDir = parent
					continue
				}
				currDir = currDir.GetOrAddSubdir(dest)
				continue
			default:
				logger.Fatalf("unrecognized command: %s", tokens[0])
			}
		}
		if tokens[0] == "dir" {
			currDir.GetOrAddSubdir(tokens[1])
			continue
		}
		size, err := strconv.ParseInt(tokens[0], 10, 64)
		if err != nil {
			logger.Fatalf("unrecognized token for line: %s", tokens)
		}
		currDir.AddFile(tokens[1], size)
	}
	dirs := rootDir.GetDirsUnderSize(int64(100000))
	sizeOfSmallDirs := int64(0)
	for _, d := range dirs {
		sizeOfSmallDirs += d.GetSize()
	}
	logger.Printf("Total size: %d", sizeOfSmallDirs)
	totalSpace := int64(70000000)
	neededSpace := int64(30000000)
	currentSpace := totalSpace - rootDir.GetSize()
	toFree := neededSpace - currentSpace
	toDelete := rootDir.FindSmallestDirGreaterThan(toFree)
	logger.Printf("Size of directory to delete: %d", toDelete.GetSize())
}

func ParseLine(l string) ([]string, bool) {
	tokens := strings.Split(l, " ")
	if tokens[0] == "$" {
		return tokens[1:], true
	}
	return tokens, false
}

type Dir struct {
	name    string
	subdirs []*Dir
	files   map[string]int64
	parent  *Dir
	size    *int64
}

func NewDir(name string, parent *Dir) *Dir {
	return &Dir{name: name, parent: parent, files: make(map[string]int64)}
}

func (this *Dir) AddFile(name string, size int64) {
	this.files[name] = size
}

func (this *Dir) AddDir(d *Dir) {
	d.parent = this
	this.subdirs = append(this.subdirs, d)
}

func (this *Dir) GetSubdir(name string) *Dir {
	for _, subdir := range this.subdirs {
		if subdir.name == name {
			return subdir
		}
	}
	return nil
}

func (this *Dir) GetOrAddSubdir(name string) *Dir {
	subdir := this.GetSubdir(name)
	if subdir == nil {
		subdir = NewDir(name, this)
		this.AddDir(subdir)
	}
	return subdir
}

func (this *Dir) GetSize() int64 {
	if this.size != nil {
		return *this.size
	}
	fs := int64(0)
	for _, size := range this.files {
		fs += size
	}
	sds := int64(0)
	for _, sd := range this.subdirs {
		sds += sd.GetSize()
	}
	size := fs + sds
	this.size = &size
	return size
}

func (this *Dir) GetDirsUnderSize(threshold int64) []*Dir {
	return this.GetDirsCompSize(func(a int64) bool { return a < threshold })
}

func (this *Dir) GetDirsCompSize(comp func(a int64) bool) []*Dir {
	dirs := make([]*Dir, 0)
	size := this.GetSize()
	if comp(size) {
		dirs = append(dirs, this)
	}
	for _, sd := range this.subdirs {
		dirs = append(dirs, sd.GetDirsCompSize(comp)...)
	}
	return dirs
}

func (this *Dir) FindSmallestDirGreaterThan(threshold int64) *Dir {
	dirs := this.GetDirsCompSize(func(a int64) bool { return a > threshold })
	v, d := int64(1<<63-1), (*Dir)(nil)
	for _, dir := range dirs {
		size := dir.GetSize()
		if size < v {
			v, d = size, dir
		}
	}
	return d
}
