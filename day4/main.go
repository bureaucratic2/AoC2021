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
	game := gameFromStr(s)
	part1(game)
	part2(game)
}

func part1(game game) {
	for _, num := range game.nums[:4] {
		for idx := range game.boards {
			board := &game.boards[idx]
			board.mark(num)
		}
	}

	for _, num := range game.nums[4:] {
		for i := 0; i < len(game.boards); i++ {
			game.boards[i].mark(num)
			if game.boards[i].win() {
				sum := game.boards[i].sum()
				fmt.Println("Part1:", num*sum)
				return
			}
		}
	}
}

func part2(game game) {
	winned := make(map[int]bool)
	record := make([]int, 0)

	for _, num := range game.nums[:4] {
		for idx := range game.boards {
			board := &game.boards[idx]
			board.mark(num)
		}
	}

	for _, num := range game.nums[4:] {
		for i := 0; i < len(game.boards); i++ {
			game.boards[i].mark(num)

			if winned[i] {
				continue
			}

			if game.boards[i].win() {
				sum := game.boards[i].sum()
				winned[i] = true
				record = append(record, num*sum)
			}
		}
	}
	fmt.Println("Part2:", record[len(record)-1])
}

type game struct {
	nums   []int
	boards []board
}

func gameFromStr(s string) game {
	game := new(game)
	scanner := bufio.NewScanner(strings.NewReader(s))

	// get first line of numbers
	scanner.Scan()
	for _, num := range strings.Split(scanner.Text(), ",") {
		num, err := strconv.Atoi(num)
		common.Check(err)
		game.nums = append(game.nums, num)
	}

	// skip blank line
	for scanner.Scan() {
		newBoard := make([][][]int, 0, 5)
		for i := 0; i < 5; i++ {
			newRow := make([][]int, 0, 5)
			newRow = append(newRow, make([]int, 0))
			newRow = append(newRow, make([]int, 0))
			scanner.Scan()
			for _, num := range strings.Split(scanner.Text(), " ") {
				if num != "" {
					num, err := strconv.Atoi(num)
					common.Check(err)
					newRow[0] = append(newRow[0], num)
					newRow[1] = append(newRow[1], 0)
				}
			}

			// san check
			if len(newRow[0]) != 5 {
				panic("parse error")
			}

			newBoard = append(newBoard, newRow)
		}
		game.boards = append(game.boards, newBoard)
	}
	return *game
}

// 5x2x5, 1st value store number, 2nd value store mark
type board [][][]int

// terrible reference experience in Go
func (b *board) mark(num int) {
	for j := 0; j < 5; j++ {
		for i := 0; i < 5; i++ {
			if (*b)[j][0][i] == num {
				(*b)[j][1][i] = 1
			}
		}
	}
}

// terrible reference experience in Go
func (b *board) win() bool {
	colCount := make([]int, 5)

	for j := 0; j < 5; j++ {
		var rowCount int
		for i := 0; i < 5; i++ {
			val := &(*b)[j][1][i]
			rowCount += *val
			colCount[i] += *val
		}
		if rowCount == 5 {
			return true
		}
	}

	for _, count := range colCount {
		if count == 5 {
			return true
		}
	}

	return false
}

// terrible reference experience in Go
func (b *board) sum() int {
	var count int

	for j := 0; j < 5; j++ {
		for i := 0; i < 5; i++ {
			if (*b)[j][1][i] == 0 {
				count += (*b)[j][0][i]
			}
		}
	}

	return count
}
