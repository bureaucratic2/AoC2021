package main

import (
	"bufio"
	"common"
	"strings"
)

func main() {
	s := common.Load()
	part1(s)
}

func part1(s string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	arr := make([][]int, 0)

	// initialize
	scanner.Scan()
	lenOfBits := len(scanner.Text())
	arr = append(arr, make([]int, lenOfBits))
	arr = append(arr, make([]int, lenOfBits))
	for idx, ch := range scanner.Text() {
		arr[ch][idx]++
	}

	for scanner.Scan() {
		for idx, ch := range scanner.Text() {
			arr[ch][idx]++

		}
	}

}
