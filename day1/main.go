package main

import (
	"bufio"
	"common"
	"fmt"
	"strconv"
	"strings"

	"github.com/gammazero/deque"
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
	var q deque.Deque

	for i := 0; i < 3; i++ {
		scanner.Scan()
		next, err := strconv.Atoi(scanner.Text())
		common.Check(err)
		q.PushBack(next)
	}

	var count int
	for scanner.Scan() {
		next, err := strconv.Atoi(scanner.Text())
		common.Check(err)
		if q.PopFront().(int) < next {
			count += 1
		}
		q.PushBack(next)
	}

	return count
}
