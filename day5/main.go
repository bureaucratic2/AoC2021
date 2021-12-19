package main

import (
	"bufio"
	"common"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	s := common.Load()
	vs := ventsFromStr(s)
	part1(&vs)
	part2(&vs)
}

func part1(vs *vents) {
	lines := vs.filterLines()
	overlap := count(lines, make([]vent, 0))
	fmt.Println("Part1:", overlap)
}

func part2(vs *vents) {
	lines := vs.filterLines()
	diagonal := vs.filterDiagonalLines()

	overlap := count(lines, diagonal)
	fmt.Println("Part2:", overlap)
}

func count(lines []vent, diagonal []vent) int {
	area := make([][]int, 0)
	for i := 0; i < 1000; i++ {
		area = append(area, make([]int, 1000))
	}

	var overlap int
	for _, v := range lines {
		for y := v.from.y; y <= v.to.y; y++ {
			for x := v.from.x; x <= v.to.x; x++ {
				area[y][x]++
				if area[y][x] == 2 {
					overlap++
				}
			}
		}
	}

	for _, v := range diagonal {
		if v.from.x < v.to.x {
			for i := 0; i+v.from.x <= v.to.x; i++ {
				area[v.from.y+i][v.from.x+i]++
				if area[v.from.y+i][v.from.x+i] == 2 {
					overlap++
				}
			}
		} else {
			for i := 0; i+v.from.x >= v.to.x; i-- {
				area[v.from.y-i][v.from.x+i]++
				if area[v.from.y-i][v.from.x+i] == 2 {
					overlap++
				}
			}
		}
	}

	return overlap
}

type vents struct {
	vents []vent
}

func (vs *vents) filterLines() []vent {
	ret := make([]vent, 0)

	for i := range vs.vents {
		v := &vs.vents[i]
		if v.from.x == v.to.x || v.from.y == v.to.y {
			ret = append(ret, *v)
		}
	}

	return ret
}

func (vs *vents) filterDiagonalLines() []vent {
	ret := make([]vent, 0)

	for i := range vs.vents {
		v := &vs.vents[i]
		if v.to.x-v.from.x == v.to.y-v.from.y || v.to.x-v.from.x == v.from.y-v.to.y {
			ret = append(ret, *v)
		}
	}

	return ret
}

func ventsFromStr(s string) vents {
	v := new(vents)

	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		v.vents = append(v.vents, ventFromStr(scanner.Text()))
	}
	return *v
}

type vent struct {
	from coordinate
	to   coordinate
}

func ventFromStr(s string) vent {
	arr := strings.Split(s, " -> ")
	c1, c2 := coordinateFromStr(arr[0]), coordinateFromStr(arr[1])
	// ensure lines from left to right, from top to down

	if c1.x == c2.x {
		if c1.y >= c2.y {
			c1.y, c2.y = c2.y, c1.y
		}
	} else if c1.y == c2.y {
		if c1.x >= c2.x {
			c1.x, c2.x = c2.x, c1.x
		}
	} else if c1.y >= c2.y {
		c1.x, c2.x = c2.x, c1.x
		c1.y, c2.y = c2.y, c1.y
	}

	return vent{c1, c2}
}

type coordinate struct {
	x int
	y int
}

func coordinateFromStr(s string) coordinate {
	arr := make([]int, 0, 2)
	for _, n := range strings.Split(s, ",") {
		n, err := strconv.Atoi(n)
		common.Check(err)
		arr = append(arr, n)
	}
	if len(arr) != 2 {
		panic("sanity check fail")
	}
	return coordinate{arr[0], arr[1]}
}
