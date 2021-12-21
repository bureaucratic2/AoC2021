package main

import (
	"bufio"
	"common"
	"container/list"
	"fmt"
	"strings"
)

func main() {
	s := common.Load()
	l, pairs := listAndPairsFromStr(s)
	part1(l, pairs)
	// reconstruct states
	l, pairs = listAndPairsFromStr(s)
	part2(l, pairs)
}

func part1(l *list.List, pairs map[pair]rune) {
	fmt.Println("Part1:", run(l, pairs, 10))
}

func part2(l *list.List, pairs map[pair]rune) {
	// toooooo slow
	// fmt.Println("Part2:", run(l, pairs, 40))

	pairCounts := make(map[pair]int64)
	e := l.Front()
	for {
		next := e.Next()
		if next == nil {
			break
		}
		p := pair{e.Value.(rune), next.Value.(rune)}
		pairCounts[p]++
		e = next
	}

	for i := 0; i < 40; i++ {
		// temporarily store new added pairs
		tmp := make(map[pair]int64)
		for k, v := range pairCounts {
			if insertion, exist := pairs[k]; exist {
				tmp[pair{k.first, insertion}] += v
				tmp[pair{insertion, k.second}] += v
				// pair k is broke down into pieces, delete it
				pairCounts[k] = 0
			}
		}

		// add tmp result to pairCounts
		for k, v := range tmp {
			pairCounts[k] += v
		}
	}

	count := make(map[rune]int64)
	for k, v := range pairCounts {
		count[k.first] += v
		count[k.second] += v
	}

	most, least := common.MinInt64, common.MaxInt64
	for _, v := range count {
		if v > most {
			most = v
		}
		if v < least {
			least = v
		}
	}

	fmt.Println("Part2:", (most-least+1)/2)
}

func run(l *list.List, pairs map[pair]rune, times int) int {
	for i := 0; i < times; i++ {
		e := l.Front()
		for {
			next := e.Next()
			if next == nil {
				break
			}
			p := pair{e.Value.(rune), next.Value.(rune)}
			if insertion, exist := pairs[p]; exist {
				l.InsertAfter(insertion, e)
			}
			e = next
		}
	}

	count := make(map[rune]int)
	for e := l.Front(); e != nil; e = e.Next() {
		count[e.Value.(rune)]++
	}

	most, least := common.MinInt, common.MaxInt
	for _, v := range count {
		if v > most {
			most = v
		}
		if v < least {
			least = v
		}
	}

	return most - least
}

// must return *list.List
// if return list.List, the memory layout will be broke down!!!
func listAndPairsFromStr(s string) (*list.List, map[pair]rune) {
	scanner := bufio.NewScanner(strings.NewReader(s))

	l := list.New()
	scanner.Scan()
	s = scanner.Text()
	for _, ch := range s {
		l.PushBack(ch)
	}

	scanner.Scan()
	pairs := make(map[pair]rune)
	tuple := make([]rune, 2)
	for scanner.Scan() {
		sub := strings.Split(scanner.Text(), " -> ")
		for i, ch := range sub[0] {
			tuple[i] = ch
		}
		// so awkward
		for _, ch := range sub[1] {
			pairs[pair{tuple[0], tuple[1]}] = ch
		}
	}

	return l, pairs
}

type pair struct {
	first  rune
	second rune
}
