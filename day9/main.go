package main

import (
	"bufio"
	"common"
	"fmt"
	"sort"
	"strings"
)

func main() {
	s := common.Load()
	v := fromStr(s)
	part1(v)
	part2(v)
}

func part1(v heightMap) {
	res := make([]int, 0)
	yLen, xLen := len(v), len(v[0])

	for y, row := range v {
		for x, height := range row {
			if y-1 >= 0 && v[y-1][x] <= height {
				continue
			}
			if x-1 >= 0 && v[y][x-1] <= height {
				continue
			}
			if y+1 < yLen && v[y+1][x] <= height {
				continue
			}
			if x+1 < xLen && v[y][x+1] <= height {
				continue
			}
			res = append(res, height)
		}
	}

	sum := 0
	for _, height := range res {
		sum += height + 1
	}
	fmt.Println("Part1:", sum)
}

// bfs
func part2(v heightMap) {
	res := make([]int, 0)
	yLen, xLen := len(v), len(v[0])
	visited := make(map[coordinate]bool)

	for y, row := range v {
		for x, height := range row {
			cur := coordinate{x, y}
			if visited[cur] || height == 9 {
				continue
			}

			basin := 0
			// if deadlock, increase buffer
			queue := make(chan coordinate, 175)
			queue <- cur

		outer:
			for {
				select {
				case c := <-queue:
					if c.y < 0 || c.y >= yLen || c.x < 0 || c.x >= xLen {
						continue
					}
					if v[c.y][c.x] == 9 || visited[c] {
						continue
					}
					basin++
					visited[c] = true
					for _, neighbor := range c.around() {
						queue <- neighbor
					}
				default:
					break outer
				}
			}

			res = append(res, basin)
			close(queue)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] > res[j]
	})

	fmt.Println("Part2:", res[0]*res[1]*res[2])
}

type heightMap [][]int

func fromStr(s string) heightMap {
	v := make(heightMap, 0)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		tmp := make([]int, 0, len(scanner.Text()))
		for _, ch := range scanner.Text() {
			tmp = append(tmp, int(ch-'0'))
		}
		v = append(v, tmp)
	}
	return v
}

type coordinate struct {
	x int
	y int
}

func (c coordinate) around() []coordinate {
	neighbor := make([]coordinate, 0, 4)
	for _, y := range []int{-1, 1} {
		neighbor = append(neighbor, coordinate{c.x, c.y + y})
	}
	for _, x := range []int{-1, 1} {
		neighbor = append(neighbor, coordinate{c.x + x, c.y})
	}
	return neighbor
}
