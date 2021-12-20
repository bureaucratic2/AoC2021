package main

import (
	"common"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	s := common.Load()
	part1(s)
	part2(s)
}

func part1(s string) {
	pos, idx, max := fromStr(s)
	var left, right, occupied, fuel int

	// initialize
	for k, v := range pos {
		if k == idx {
			occupied = v
		} else if k < idx {
			left += v
			fuel += v * (idx - k)
		} else {
			right += v
			fuel += v * (k - idx)
		}
	}
	idx++

	minFuel := fuel
	for ; idx < max; idx++ {
		if v, exist := pos[idx]; exist {
			left += occupied
			right -= v
			occupied = v
		} else {
			left += occupied
			occupied = 0
		}
		fuel = fuel + left - right - occupied
		if minFuel > fuel {
			minFuel = fuel
		} else {
			break
		}
	}

	fmt.Println("Part1:", minFuel)
}

func part2(s string) {
	pos, idx, max := fromStr(s)
	var occupied, fuel int
	left := make(map[int]int)
	right := make(map[int]int)

	// initialize
	for k, v := range pos {
		if k == idx {
			occupied = v
		} else if k < idx {
			left[k] = v
			fuel += v * dist(idx-k)
		} else {
			right[k] = v
			fuel += v * dist(k-idx)
		}
	}
	idx++

	minFuel := fuel
	for ; idx < max; idx++ {
		if v, exist := pos[idx]; exist {
			if occupied > 0 {
				left[idx-1] = occupied
			}
			delete(right, idx)
			occupied = v
		} else {
			if occupied > 0 {
				left[idx-1] = occupied
			}
			occupied = 0
		}

		fuel = 0
		for k, v := range left {
			fuel += v * dist(idx-k)
		}
		for k, v := range right {
			fuel += v * dist(k-idx)
		}
		if minFuel > fuel {
			minFuel = fuel
		} else {
			break
		}
	}

	fmt.Println("Part2:", minFuel)
}

func dist(n int) int {
	if n < 0 {
		panic("negative distance")
	}
	return (1 + n) * n / 2
}

func fromStr(s string) (map[int]int, int, int) {
	min, max := common.MaxInt, common.MinInt
	pos := make(map[int]int)
	for _, n := range strings.Split(s, ",") {
		n, err := strconv.Atoi(n)
		common.Check(err)
		if min > n {
			min = n
		}
		if max < n {
			max = n
		}
		pos[n]++
	}
	return pos, min, max
}
