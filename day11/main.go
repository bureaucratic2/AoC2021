package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"
)

var xLen int
var yLen int
var sub = []int{-1, 0, 1}

func main() {
	s := common.Load()
	part1(octopusesFromStr(s))
	part2(octopusesFromStr(s))
}

func part1(o octopuses) {
	for i := 0; i < 100; i++ {
		o.step()
	}
	fmt.Println("Part1:", o.flashes)
}

func part2(o octopuses) {
	for i := 1; ; i++ {
		if o.step() {
			fmt.Println("Part2:", i)
			break
		}
	}
}

func (o *octopuses) step() bool {
	flag := true
	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			o.energyLevel[y][x]++
			if o.energyLevel[y][x] > 9 && !o.flag[y][x] {
				o.flag[y][x] = true
				o.flash(x, y)
			}
		}
	}

	for y := 0; y < yLen; y++ {
		for x := 0; x < xLen; x++ {
			if o.flag[y][x] {
				o.flag[y][x] = false
				o.energyLevel[y][x] = 0
			} else {
				flag = false
			}
		}
	}
	return flag
}

func (o *octopuses) flash(x int, y int) {
	o.flashes++
	for _, subY := range sub {
		for _, subX := range sub {
			newX := x + subX
			newY := y + subY

			if newX == x && newY == y {
				continue
			}
			if newX < 0 || newY < 0 || newX >= xLen || newY >= yLen {
				continue
			}

			o.energyLevel[newY][newX]++
			if o.energyLevel[newY][newX] > 9 && !o.flag[newY][newX] {
				o.flag[newY][newX] = true
				o.flash(newX, newY)
			}
		}
	}
}

// for debug
// func (o *octopuses) toStr() {
// 	for y := range o.energyLevel {
// 		fmt.Println(o.energyLevel[y])
// 	}
// 	fmt.Println()
// }

type octopuses struct {
	energyLevel [][]int
	flag        [][]bool
	flashes     int
}

func octopusesFromStr(s string) octopuses {
	scanner := bufio.NewScanner(strings.NewReader(s))
	energyLevel := make([][]int, 0)
	flag := make([][]bool, 0)
	for scanner.Scan() {
		tmp := make([]int, 0)
		for _, ch := range scanner.Text() {
			tmp = append(tmp, int(ch-'0'))
		}
		flag = append(flag, make([]bool, len(tmp)))
		energyLevel = append(energyLevel, tmp)
	}
	xLen = len(energyLevel[0])
	yLen = len(energyLevel)
	return octopuses{energyLevel, flag, 0}
}
