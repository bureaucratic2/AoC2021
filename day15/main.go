package main

import (
	"bufio"
	"common"
	"container/heap"
	"fmt"
	"strings"
)

var movement = []int{-1, 1}

func main() {
	s := common.Load()
	c := cavernFromStr(s)
	part1(&c)
	part2(c)
}

func part1(c *cavern) {
	src := coordinate{0, 0}
	dst := coordinate{c.xLen - 1, c.yLen - 1}
	res := c.shortestPath(src, dst)
	fmt.Println("Part1:", res)
}

func part2(c cavern) {
	c = c.extend()
	src := coordinate{0, 0}
	dst := coordinate{c.xLen - 1, c.yLen - 1}
	res := c.shortestPath(src, dst)
	fmt.Println("Part2:", res)
}

type cavern struct {
	risk [][]int
	xLen int
	yLen int
}

func (c *cavern) shortestPath(src coordinate, dst coordinate) int {
	dist := make([]int, c.xLen*c.yLen)
	for i := range dist {
		dist[i] = common.MaxInt
	}
	dist[src.y*c.yLen+src.x] = 0

	visited := make(map[int]bool)

	pq := make(priorityQueue, 0)
	pq = append(pq, &pointWithDistance{src, 0, 0})
	heap.Init(&pq)

	i := 1
	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*pointWithDistance)
		if visited[u.index] {
			continue
		}
		x, y := u.c.x, u.c.y
		s := y*c.yLen + x
		for _, sub := range movement {
			if x+sub >= 0 && x+sub < c.xLen {
				v := y*c.yLen + x + sub
				weight := c.risk[y][x+sub]
				if dist[v] > dist[s]+weight {
					dist[v] = dist[s] + weight
					heap.Push(&pq, &pointWithDistance{coordinate{x + sub, y}, dist[v], i})
					i++
				}
			}
			if y+sub >= 0 && y+sub < c.yLen {
				v := (y+sub)*c.yLen + x
				weight := c.risk[y+sub][x]
				if dist[v] > dist[s]+weight {
					dist[v] = dist[s] + weight
					heap.Push(&pq, &pointWithDistance{coordinate{x, y + sub}, dist[v], i})
					i++
				}
			}
		}
	}

	return dist[dst.y*c.yLen+dst.x]
}

func (c cavern) extend() cavern {
	risk := make([][]int, c.yLen*5)
	for i := range risk {
		risk[i] = make([]int, c.xLen*5)
	}

	ch := make([]chan int, 0)
	ch = append(ch, make(chan int, 5))
	ch = append(ch, make(chan int, 5))

	for i := 0; i < 5; i++ {
		ch[0] <- i
	}

	for j := 0; j < 5; j++ {
		cur := j % 2
		for i := 0; i < 5; i++ {
			adder := <-ch[cur]
			ch[1-cur] <- adder + 1
			for y, row := range c.risk {
				for x, val := range row {
					risk[j*c.yLen+y][i*c.xLen+x] = (val+adder-1)%9 + 1
				}
			}
		}
	}

	return cavern{risk, c.xLen * 5, c.yLen * 5}
}

func cavernFromStr(s string) cavern {
	scanner := bufio.NewScanner(strings.NewReader(s))
	risk := make([][]int, 0)
	for scanner.Scan() {
		row := make([]int, 0)
		for _, digit := range scanner.Text() {
			row = append(row, int(digit-'0'))
		}
		risk = append(risk, row)
	}

	return cavern{risk, len(risk[0]), len(risk)}
}

type coordinate struct {
	x int
	y int
}

type pointWithDistance struct {
	c     coordinate
	dist  int
	index int
}

type priorityQueue []*pointWithDistance

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pointWithDistance)
	item.index = n
	*pq = append(*pq, item)
}

// // update modifies the distance of a point in the queue.
// func (pq *priorityQueue) update(item *pointWithDistance, dist int) {
// 	item.dist = dist
// 	heap.Fix(pq, item.index)
// }
