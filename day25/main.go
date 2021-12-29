package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"
)

func main() {
	s := common.Load()
	sf := seafloorFromStr(s)
	part1(&sf)
}

func part1(sf *seafloor) {
	for i := 0; ; i++ {
		if !sf.step() {
			fmt.Println("Part1:", i+1)
			return
		}
	}
}

type seafloor struct {
	floor    [][]occupied
	eastLen  int
	southLen int
	canMove  [][]bool
}

func (sf *seafloor) step() bool {
	moved := false

	// first, consider east-facing herd
	for y := range sf.floor {
		for x := range sf.floor[y] {
			if sf.floor[y][x] == east {
				ny, nx := y, (x+1)%sf.eastLen
				if sf.floor[ny][nx] == empty {
					sf.canMove[y][x] = true
					moved = true
				}
			}
		}
	}

	for y := range sf.canMove {
		for x := range sf.canMove[y] {
			if sf.canMove[y][x] {
				ny, nx := y, (x+1)%sf.eastLen
				sf.canMove[y][x] = false
				sf.floor[ny][nx] = east
				sf.floor[y][x] = empty
			}
		}
	}

	// then, south-facing herd move
	for y := range sf.floor {
		for x := range sf.floor[y] {
			if sf.floor[y][x] == south {
				ny, nx := (y+1)%sf.southLen, x
				if sf.floor[ny][nx] == empty {
					sf.canMove[y][x] = true
					moved = true
				}
			}
		}
	}

	for y := range sf.canMove {
		for x := range sf.canMove[y] {
			if sf.canMove[y][x] {
				ny, nx := (y+1)%sf.southLen, x
				sf.canMove[y][x] = false
				sf.floor[ny][nx] = south
				sf.floor[y][x] = empty
			}
		}
	}
	return moved
}

func seafloorFromStr(s string) seafloor {
	scanner := bufio.NewScanner(strings.NewReader(s))
	floor := make([][]occupied, 0)

	for scanner.Scan() {
		tmp := make([]occupied, 0)
		for _, ch := range scanner.Text() {
			switch ch {
			case '.':
				tmp = append(tmp, empty)
			case '>':
				tmp = append(tmp, east)
			case 'v':
				tmp = append(tmp, south)
			}
		}
		floor = append(floor, tmp)
	}
	eastLen := len(floor[0])
	southLen := len(floor)
	canMove := make([][]bool, 0)
	for i := 0; i < southLen; i++ {
		canMove = append(canMove, make([]bool, eastLen))
	}

	return seafloor{floor, eastLen, southLen, canMove}
}

type occupied int

const (
	empty occupied = iota
	east
	south
)
