package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"
)

var digits = []map[rune]bool{
	{'a': true, 'b': true, 'c': true, 'e': true, 'f': true, 'g': true}, // 0
	{'c': true, 'f': true},                                                        // 1
	{'a': true, 'c': true, 'd': true, 'e': true, 'g': true},                       // 2
	{'a': true, 'c': true, 'd': true, 'f': true, 'g': true},                       // 3
	{'b': true, 'c': true, 'd': true, 'f': true},                                  // 4
	{'a': true, 'b': true, 'd': true, 'f': true, 'g': true},                       // 5
	{'a': true, 'b': true, 'd': true, 'e': true, 'f': true, 'g': true},            // 6
	{'a': true, 'c': true, 'f': true},                                             // 7
	{'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true}, // 8
	{'a': true, 'b': true, 'c': true, 'd': true, 'f': true, 'g': true},            // 9
}

func main() {
	s := common.Load()
	part1(s)
	part2(s)
}

func part1(s string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	count := 0
	for scanner.Scan() {
		s = strings.Split(scanner.Text(), " | ")[1]
		for _, sub := range strings.Split(s, " ") {
			switch len(sub) {
			case 2, 3, 4, 7:
				count++
			default:
			}
		}
	}

	fmt.Println("Part1:", count)
}

func part2(s string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	res := 0
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " | ")
		patterns := make([]string, 0)
		output := make([]string, 0)
		patterns = append(patterns, strings.Split(s[0], " ")...)
		output = append(output, strings.Split(s[1], " ")...)
		mapping := parse(patterns)
		res += toDigits(mapping, output)
	}
	fmt.Println("Part2:", res)
}

// 1,4,7,8 is known at first
// 1,7 --> a
// 4,9 --> (9, g)
// 8,9 --> e
// e --> 2
// 2,e --> (3, 5, f)
// 3,5 --> (b, c)
// last is d
func parse(patterns []string) map[rune]rune {
	mappings := make([][]map[rune]bool, 0)

	// for convenience, althouth 0, 1 is never used
	for i := 0; i < 8; i++ {
		mappings = append(mappings, make([]map[rune]bool, 0))
	}

	for _, pattern := range patterns {
		mapping := make(map[rune]bool)
		for _, ch := range pattern {
			mapping[ch] = true
		}
		mappings[len(pattern)] = append(mappings[len(pattern)], mapping)
	}

	mapping := make(map[rune]rune)
	// find 7
	for k := range mappings[3][0] {
		// find 1 and a
		if !mappings[2][0][k] {
			mapping[k] = 'a'
			break
		}
	}

	var nine int
	// traverse 0, 6, 9
outer:
	for idx, lenSixDigit := range mappings[6] {
		// 9 conatins 4
		for k := range mappings[4][0] {
			if !lenSixDigit[k] {
				continue outer
			}
		}
		// find 9
		nine = idx
		for k := range lenSixDigit {
			if mappings[4][0][k] {
				continue
			}
			if _, exist := mapping[k]; exist {
				continue
			}
			mapping[k] = 'g'
			break
		}

		break
	}

	var e rune
	for k := range mappings[7][0] {
		if !mappings[6][nine][k] {
			mapping[k] = 'e'
			e = k
			break
		}
	}

	var two int
	// traverse 2, 3, 5
	for idx, lenFiveDigit := range mappings[5] {
		if lenFiveDigit[e] {
			two = idx
			break
		}
	}

	var three int
	var five int
	// traverse 2, 3, 5
	for idx, lenFiveDigit := range mappings[5] {
		if idx == two {
			continue
		}
		diff := 0
		for k := range mappings[5][two] {
			if !lenFiveDigit[k] {
				diff++
			}
		}
		if diff == 1 {
			three = idx
		} else {
			five = idx
			continue
		}
		for k := range lenFiveDigit {
			if !mappings[5][two][k] {
				mapping[k] = 'f'
				break
			}
		}
	}

	for k := range mappings[5][three] {
		if !mappings[5][five][k] {
			mapping[k] = 'c'
			break
		}
	}

	for k := range mappings[5][five] {
		if !mappings[5][three][k] {
			mapping[k] = 'b'
			break
		}
	}

	// last is d
	segments := map[rune]bool{'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true}
	for k := range mapping {
		delete(segments, k)
	}
	for k := range segments {
		mapping[k] = 'd'
	}

	return mapping
}

func toDigits(mapping map[rune]rune, output []string) int {
	multiplier := 1000
	res := 0
	for _, s := range output {
	outer:
		for num, digit := range digits {
			if len(s) != len(digit) {
				continue
			}
			for _, ch := range s {
				if !digit[mapping[ch]] {
					continue outer
				}
			}
			res += num * multiplier
			multiplier /= 10
			break
		}
	}

	return res
}
