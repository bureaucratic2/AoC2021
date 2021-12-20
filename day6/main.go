package main

import (
	"common"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	s := common.Load()

	state := stateFromStr(s)
	part1(state, 80)
	// reconstruct state, avoid using the same memory
	state = stateFromStr(s)
	part2(state, 256)
}

func part1(s state, days int) {
	fmt.Println("Part1:", compute(s, days))
}

func part2(s state, days int) {
	fmt.Println("Part2:", compute(s, days))
}

func compute(s state, days int) int {
	lookUpTable := make([][]int, 0)
	tmp := make([]int, 9)
	for i := 0; i < 9; i++ {
		tmp[i] = 1
		lookUpTable = append(lookUpTable, cycle7(&tmp))
		tmp[i] = 0
	}

	// res computation start from 14 days
	res := make([][][]int, 0)
	res = append(res, make([][]int, 0))
	res = append(res, make([][]int, 0))
	for i := 0; i < 9; i++ {
		res[0] = append(res[0], make([]int, 9))
		res[1] = append(res[1], make([]int, 9))
		copy(res[0][i], lookUpTable[i])
	}
	curr := 0

	for i := 1; i < days/7; i++ {
		for idx := 0; idx < 9; idx++ {
			for day, multiplier := range res[curr][idx] {
				for k, v := range lookUpTable[day] {
					res[1-curr][idx][k] += multiplier * v
					res[curr][idx][day] = 0
				}
			}
		}
		curr = 1 - curr
	}

	for day, multiplier := range s.fishes[s.curr] {
		for k, v := range res[curr][day] {
			s.fishes[1-s.curr][k] += multiplier * v
			s.fishes[s.curr][day] = 0
		}
	}
	s.curr = 1 - s.curr
	c := sum(s.fishes[s.curr])
	c += sum(s.fishes[s.curr][:(days % 7)])
	return c
}

// state: 012345678
// for each 7 days
// [01]+=[78]' (' means last 7 day's value)
// [78]=[56]'
// [23456]+=[01234]'
// Diff = [012345678]-[012345678]'
//      = [78]' + [01234]' + [56]' - [78]'
//      = [0123456]'
func cycle7(s *[]int) []int {
	tmp := make([]int, 9)
	copy(tmp, *s)
	tmp[0] += (*s)[7]
	tmp[1] += (*s)[8]
	tmp[2] += (*s)[0]
	tmp[3] += (*s)[1]
	tmp[4] += (*s)[2]
	tmp[5] += (*s)[3]
	tmp[6] += (*s)[4]
	tmp[7] = (*s)[5]
	tmp[8] = (*s)[6]
	return tmp
}

func sum(arr []int) int {
	sum := 0
	for i := range arr {
		sum += arr[i]
	}
	return sum
}

type state struct {
	fishes [][]int
	curr   int
}

func stateFromStr(s string) state {
	fishes := make([][]int, 0)
	fishes = append(fishes, make([]int, 9))
	fishes = append(fishes, make([]int, 9))

	for _, n := range strings.Split(s, ",") {
		n, err := strconv.Atoi(n)
		common.Check(err)
		fishes[0][n]++
	}

	return state{fishes, 0}
}
