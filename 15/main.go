package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	sensors := createSensors("./input.txt")
	fmt.Printf(sensors.String())
	fmt.Printf("Example result: %d\n", sensors.nonBeaconPositionsForRow(10))
	sensors = createSensors("./input2.txt")
	fmt.Printf("Part 1: %d\n", sensors.nonBeaconPositionsForRow(200000))
}

func createSensors(fileName string) *Sensors {
	logger := log.Logger{}
	f, err := os.Open(fileName)
	if err != nil {
		logger.Fatalf("Failed to open file: %s", err.Error())
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	sensors := NewSensors()
	for scanner.Scan() {
		s := scanner.Text()
		ints := getInts(s)
		if len(ints) != 4 {
			logger.Fatalf("Expected 4 ints, got %d", len(ints))
		}
		sensorPos := Pos{
			X: ints[0],
			Y: ints[1],
		}
		beaconPos := Pos{
			X: ints[2],
			Y: ints[3],
		}
		sensors.AddSensor(NewSensor(sensorPos, beaconPos))
	}
	return sensors
}

func (p Pos) String() string {
	return fmt.Sprintf("(X: %d, Y: %d)", p.X, p.Y)
}

func getInts(s string) []int {
	ints := make([]int, 0)
	for i := 0; i < len(s); {
		j := i
		foundInt := false
		for ; j < len(s); j += 1 {
			if isIntChar(s[j]) {
				foundInt = true
				continue
			}
			break
		}
		if foundInt {
			val, _ := strconv.Atoi(s[i:j])
			ints = append(ints, val)
		}
		i = j + 1
	}
	return ints
}

func isIntChar(ch byte) bool {
	return ch >= '0' && ch <= '9' || ch == '-'
}

type Sensor struct {
	p, closestBeacon Pos
	d                int
}

func NewSensor(p, closestBeacon Pos) Sensor {
	d := dist(p, closestBeacon)
	return Sensor{p, closestBeacon, d}
}

func (s Sensor) NoBeaconAt(p Pos) bool {
	return dist(p, s.p) <= s.d
}

type Sensors struct {
	s        []Sensor
	topology map[Pos]Point
	// upper left and bottom right most Pos in the topology
	topLeft, bottomRight Pos
}

func NewSensors() *Sensors {
	return &Sensors{
		topology: make(map[Pos]Point),
	}
}

type Point string

const (
	BeaconPoint   Point = "B"
	NoBeaconPoint Point = "#"
	SensorPoint   Point = "S"
	UnknownPoint  Point = "."
)

func (s *Sensors) stringForPos(p Pos) string {
	return string(s.Get(p))
}

func (s *Sensors) String() string {
	str := ""
	for y := s.topLeft.Y; y <= s.bottomRight.Y; y += 1 {
		start := strconv.Itoa(y)
		for i := len(start); len(start) < 4; i += 1 {
			start = start + " "
		}
		str += start
		for x := s.topLeft.X; x <= s.bottomRight.X; x += 1 {
			str += s.stringForPos(Pos{X: x, Y: y})
		}
		str += "\n"
	}
	return str
}

func (s Sensor) LeftmostX() int {
	return s.p.X - s.d
}

func (s Sensor) RightmostX() int {
	return s.p.X + s.d
}

func (s Sensor) TopY() int {
	return s.p.Y - s.d
}

func (s Sensor) BottomY() int {
	return s.p.Y + s.d
}

func (s *Sensors) nonBeaconPositionsForRow(y int) int {
	count := 0
	var leftmost, rightmost int
	for _, sensor := range s.s {
		if sensor.LeftmostX() < leftmost {
			leftmost = sensor.LeftmostX()
		}
		if sensor.RightmostX() > rightmost {
			rightmost = sensor.RightmostX()
		}
	}
	p := Pos{X: leftmost, Y: y}
	for ; p.X <= s.bottomRight.X; p = right(p) {
		if s.Get(p) == NoBeaconPoint || s.Get(p) == SensorPoint {
			count += 1
		}
	}
	return count
}

func (s *Sensors) Get(p Pos) Point {
	for _, sensor := range s.s {
		if sensor.p == p {
			return SensorPoint
		}
		if sensor.closestBeacon == p {
			return BeaconPoint
		}
		if sensor.NoBeaconAt(p) {
			return NoBeaconPoint
		}
	}
	return UnknownPoint
}

func (s *Sensors) AddSensor(sensor Sensor) {
	s.s = append(s.s, sensor)
	if s.topLeft.X > sensor.LeftmostX() {
		s.topLeft = s.topLeft.SetX(sensor.LeftmostX())
	}
	if s.bottomRight.X < sensor.RightmostX() {
		s.bottomRight = s.bottomRight.SetX(sensor.RightmostX())
	}
	if s.topLeft.Y > sensor.TopY() {
		s.topLeft = s.topLeft.SetY(sensor.TopY())
	}
	if s.bottomRight.Y < sensor.BottomY() {
		s.bottomRight = s.bottomRight.SetY(sensor.BottomY())
	}
}

func dist(p, q Pos) int {
	return abs(p.X-q.X) + abs(p.Y-q.Y)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// MarkPositions marks all positions within dst distance of p with v, unless they've already been marked
func (s *Sensors) MarkPositions(p Pos, v Point, dst int) {
	for yOffset := 0; yOffset < dst; yOffset += 1 {
		for xOffset := 0; xOffset < dst; xOffset += 1 {
			for _, q := range []Pos{
				{X: p.X + xOffset, Y: p.Y + yOffset},
				{X: p.X + xOffset, Y: p.Y - yOffset},
				{X: p.X - xOffset, Y: p.Y + yOffset},
				{X: p.X - xOffset, Y: p.Y - yOffset},
			} {
				if dist(p, q) > dst {
					continue
				}
				if _, ok := s.topology[q]; !ok {
					s.topology[q] = v
				}
			}
		}
	}
}

type Pos struct{ X, Y int }

func down(p Pos) Pos { return Pos{Y: p.Y + 1, X: p.X} }

func up(p Pos) Pos { return Pos{Y: p.Y - 1, X: p.X} }

func left(p Pos) Pos { return Pos{Y: p.Y, X: p.X - 1} }

func right(p Pos) Pos { return Pos{Y: p.Y, X: p.X + 1} }

func (p Pos) LeftOf(o Pos) bool {
	return p.X < o.X
}

func (p Pos) Above(o Pos) bool {
	return p.Y < o.Y
}

func (p Pos) SetX(x int) Pos {
	p.X = x
	return p
}

func (p Pos) SetY(y int) Pos {
	p.Y = y
	return p
}
