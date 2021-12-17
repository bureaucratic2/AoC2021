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
	fmt.Println("part1: ", part1(s))
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
