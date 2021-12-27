package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"

	"github.com/gammazero/deque"
)

func main() {
	s := common.Load()
	steps, nodes := stepsAndNodesFromStr(s)
	part1(steps)
	part2(nodes)
}

func part1(steps []step) {
	limit := 50
	cubes := make(map[coordinate]bool)
	for _, s := range steps {
		for z := common.IntMax(-limit, s.min.z); z <= common.IntMin(limit, s.max.z); z++ {
			for y := common.IntMax(-limit, s.min.y); y <= common.IntMin(limit, s.max.y); y++ {
				for x := common.IntMax(-limit, s.min.x); x <= common.IntMin(limit, s.max.x); x++ {
					cubes[coordinate{x, y, z}] = s.on
				}
			}
		}
	}

	for k, v := range cubes {
		if !v {
			delete(cubes, k)
		}
	}

	fmt.Println("Part1:", len(cubes))
}

// Inspired by this clean but extremely fast approach
// https://github.com/mattHawthorn/advent_of_code_2021/blob/main/solutions/day22.py#L32-#L53
func part2(nodes []node) {
	graph := make(graph)
	res := make(map[bool]int)

	for _, n1 := range nodes {
		neighbor := make(map[node]coboid)
		for n2 := range graph {
			res := intersect(n1.c[:], n2.c[:])
			if res != nil {
				neighbor[n2] = res
			}
		}

		graph[n1] = neighbor

		ch := make(chan element, 10)
		go allPaths(graph, n1, ch)
		for {
			if elem, ok := <-ch; ok {
				if elem.positive {
					res[elem.n.on] += size(elem.intersection)
				} else {
					res[elem.n.on] -= size(elem.intersection)
				}
			} else {
				break
			}
		}
	}

	fmt.Println("Part2:", res[true])

}

func size(c coboid) int {
	min := c[0]
	max := c[1]
	return (max.x - min.x + 1) * (max.y - min.y + 1) * (max.z - min.z + 1)
}

func intersect(c1, c2 coboid) coboid {
	if c1 == nil || c2 == nil {
		return nil
	}

	res := make([]coordinate, 2)
	res[0].x = common.IntMax(c1[0].x, c2[0].x)
	res[1].x = common.IntMin(c1[1].x, c2[1].x)

	res[0].y = common.IntMax(c1[0].y, c2[0].y)
	res[1].y = common.IntMin(c1[1].y, c2[1].y)

	res[0].z = common.IntMax(c1[0].z, c2[0].z)
	res[1].z = common.IntMin(c1[1].z, c2[1].z)

	if res[0].x > res[1].x || res[0].y > res[1].y || res[0].z > res[1].z {
		return nil
	}

	return res
}

func allPaths(g graph, n node, ch chan element) {
	elem := element{true, n, n.c[:]}
	q := deque.New()
	q.PushBack(elem)
	for {
		if q.Len() == 0 {
			break
		}
		p := q.PopFront()
		elem := p.(element)
		ch <- elem

		if g[elem.n] != nil {
			for neighbor, intersection := range g[elem.n] {
				intersection = intersect(elem.intersection, intersection)
				if intersection != nil {
					q.PushBack(element{!elem.positive, neighbor, intersection})
				}

			}
		}
	}
	close(ch)
}

type step struct {
	on       bool
	min, max coordinate
}

func stepsAndNodesFromStr(s string) ([]step, []node) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	steps := make([]step, 0)
	nodes := make([]node, 0)

	for scanner.Scan() {
		min := coordinate{}
		max := coordinate{}
		s := ""
		fmt.Sscanf(scanner.Text(), "%s x=%d..%d,y=%d..%d,z=%d..%d",
			&s, &min.x, &max.x, &min.y, &max.y, &min.z, &max.z)

		if s == "on" {
			steps = append(steps, step{true, min, max})
			nodes = append(nodes, node{true, hashCoboid{min, max}})
		} else {
			steps = append(steps, step{false, min, max})
			nodes = append(nodes, node{false, hashCoboid{min, max}})
		}

	}

	return steps, nodes
}

type coordinate struct {
	x, y, z int
}

type coboid []coordinate

type hashCoboid [2]coordinate

type node struct {
	on bool
	c  hashCoboid
}

type graph map[node](map[node]coboid)

type element struct {
	positive     bool
	n            node
	intersection coboid
}
