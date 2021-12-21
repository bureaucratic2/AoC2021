package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"

	"github.com/yourbasic/graph"
)

var bigCaves = make(map[int]bool)

var smallCaves = make(map[int]bool)
var visited = make(map[int]bool)

var paths int
var start int
var end int

func main() {
	s := common.Load()
	g := graphFromStr(s)
	part1(&g)
	part2(&g)
}

func part1(g graph.Iterator) {
	// reset global states
	paths = 0
	visited = make(map[int]bool)

	bfs(g, start, false)
	fmt.Println("Part1:", paths)
}

func part2(g graph.Iterator) {
	// reset global states
	paths = 0
	visited = make(map[int]bool)

	bfs(g, start, true)
	fmt.Println("Part2:", paths)
}

func bfs(g graph.Iterator, v int, revisit bool) {
	if v == start {
		if visited[v] {
			return
		}
	}

	used := false
	if smallCaves[v] {
		if visited[v] && !revisit {
			return
		}
		if !visited[v] {
			visited[v] = true
		} else {
			revisit = false
			used = true
		}
	}

	if v == end {
		paths++
		visited[v] = false
		return
	}

	// bfs neighbor
	g.Visit(v, func(w int, c int64) (skip bool) {
		bfs(g, w, revisit)
		return
	})

	// reset visited and return
	if used {
		return
	}
	if smallCaves[v] {
		visited[v] = false
	}
}

func graphFromStr(s string) graph.Mutable {
	scanner := bufio.NewScanner(strings.NewReader(s))

	// start & end treat as small caves
	idx := 0
	seen := make(map[string]int)
	for scanner.Scan() {
		for _, node := range strings.Split(scanner.Text(), "-") {
			if _, exist := seen[node]; !exist {
				seen[node] = idx
				if node == "start" {
					start = idx
				} else if node == "end" {
					end = idx
				}
				if common.IsUpper(node) {
					bigCaves[idx] = true
				} else {
					smallCaves[idx] = true
				}
				idx++
			}
		}
	}

	g := graph.New(len(seen))
	scanner = bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		nodes := strings.Split(scanner.Text(), "-")
		from, to := nodes[0], nodes[1]
		g.AddBoth(seen[from], seen[to])
	}

	return *g
}
