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
	fmt.Println("part1:", part1(s))
	fmt.Println("part2:", part2(s))
}

func part1(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))

	scanner.Scan()
	prev, err := strconv.Atoi(scanner.Text())
	common.Check(err)

	var count int
	for scanner.Scan() {
		next, err := strconv.Atoi(scanner.Text())
		common.Check(err)
		if next > prev {
			count++
		}
		prev = next
	}

	return count
}

func part2(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	vec := make([]int, 0)

	for scanner.Scan() {
		next, err := strconv.Atoi(scanner.Text())
		common.Check(err)
		vec = append(vec, next)
	}

	if len(vec) <= 3 {
		return 0
	}

	var prev int
	for _, n := range vec[:3] {
		prev += n
	}

	var next int
	var count int
	for idx, n := range vec[3:] {
		next = prev - vec[idx] + n
		if next > prev {
			count++
		}
		prev = next
	}

	return count
}
