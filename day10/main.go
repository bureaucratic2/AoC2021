package main

import (
	"bufio"
	"common"
	"errors"
	"fmt"
	"sort"
	"strings"
)

var syntaxBoard = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var pairs = map[rune]rune{
	')': '(',
	']': '[',
	'}': '{',
	'>': '<',
}

// var revPairs = map[rune]rune{
// 	'(': ')',
// 	'[': ']',
// 	'{': '}',
// 	'<': '>',
// }

var openCharacters = map[rune]bool{
	'(': true,
	'[': true,
	'{': true,
	'<': true,
}

var completeBoard = map[rune]int{
	'(': 1,
	'[': 2,
	'{': 3,
	'<': 4,
}

func main() {
	s := common.Load()
	ss := fromStr(s)
	incomplete := part1(ss)
	part2(incomplete)
}

func part1(ss []string) []stack {
	score := 0
	incomplete := make([]stack, 0)
outer:
	for _, s := range ss {
		stack := make(stack, 0)
		for _, ch := range s {
			if openCharacters[ch] {
				stack.push(ch)
			} else {
				top, err := stack.pop()
				if err != nil || pairs[ch] != top {
					score += syntaxBoard[ch]
					continue outer
				}
			}
		}
		if !stack.isEmpty() {
			incomplete = append(incomplete, stack)
		}
	}

	fmt.Println("Part1:", score)
	return incomplete
}

func part2(incomplete []stack) {
	res := make([]int, 0)
	for _, s := range incomplete {
		score := 0
		for {
			if ch, err := s.pop(); err == nil {
				score *= 5
				score += completeBoard[ch]
			} else {
				break
			}
		}
		res = append(res, score)
	}

	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	fmt.Println("Part2:", res[len(res)/2])
}

type stack []rune

func (s *stack) push(elem rune) {
	*s = append(*s, elem)
}

func (s *stack) pop() (rune, error) {
	len := len(*s)
	if len == 0 {
		return 0, errors.New("stack is empty")
	}
	res := (*s)[len-1]
	*s = (*s)[:len-1]
	return res, nil
}

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}

func fromStr(s string) []string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	ss := make([]string, 0)

	for scanner.Scan() {
		ss = append(ss, scanner.Text())
	}

	return ss
}
