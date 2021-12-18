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
	bits := part1(s)
	part2(bits)
}

func part1(s string) []string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	arr := make([][]int, 0)
	// prepared for part2
	bits := make([]string, 0)

	// initialize
	scanner.Scan()
	lenOfBits := len(scanner.Text())
	arr = append(arr, make([]int, lenOfBits))
	arr = append(arr, make([]int, lenOfBits))
	for idx, ch := range scanner.Text() {
		arr[ch-'0'][idx]++
	}
	bits = append(bits, scanner.Text())

	for scanner.Scan() {
		for idx, ch := range scanner.Text() {
			arr[ch-'0'][idx]++
		}
		bits = append(bits, scanner.Text())
	}

	var gamma, epsilon uint
	for i := 0; i < lenOfBits; i++ {
		if arr[0][i] < arr[1][i] {
			gamma += 1 << (lenOfBits - i - 1)
		} else {
			epsilon += 1 << (lenOfBits - i - 1)
		}
	}

	fmt.Println("Part1:", gamma, epsilon, gamma*epsilon)
	return bits
}

func part2(bits []string) {
	o2 := bits
	co2 := make([]string, len(o2))
	copy(co2, o2)

	limit := len(o2)
	for idx := 0; idx < limit; idx++ {
		var most byte
		if mostCommon(o2, idx) >= 0 {
			most = '1'
		} else {
			most = '0'
		}

		o2 = filter(o2, idx, most)
		if len(o2) == 1 {
			break
		}
	}

	limit = len(co2)
	for idx := 0; idx < limit; idx++ {
		var most byte
		if mostCommon(co2, idx) >= 0 {
			most = '0'
		} else {
			most = '1'
		}

		co2 = filter(co2, idx, most)
		if len(co2) == 1 {
			break
		}
	}

	o2_val, _ := strconv.ParseUint(o2[0], 2, len(o2[0]))
	co2_val, _ := strconv.ParseUint(co2[0], 2, len(co2[0]))
	fmt.Println("Part2:", o2_val*co2_val)
}

func mostCommon(bits []string, idx int) (ret int) {
	for _, s := range bits {
		ret += int(s[idx]-'0')*2 - 1
	}
	return
}

func filter(bits []string, idx int, most byte) []string {
	var n int
	for _, s := range bits {
		if s[idx] == most {
			bits[n] = s
			n++
		}
	}
	return bits[:n]
}
