package main

import (
	"bufio"
	"common"
	"container/heap"
	"fmt"
	"strings"
)

func main() {
	s := common.Load()
	node := initialNodeFromStr(s)
	part1(node)
	node = initialNodeFromStr(s)
	part2(node)
}

func part1(n node) {
	for i := range n.s.rooms {
		n.s.rooms[i][2] = i*2 + 2
		n.s.rooms[i][3] = i*2 + 2
	}

	g, dst := traverse(n)
	res := dijkstra(n, g, dst)

	fmt.Println("Part1:", res)
}

func part2(n node) {
	for i := range n.s.rooms {
		n.s.rooms[i][3] = n.s.rooms[i][1]
	}
	n.s.rooms[0][1] = D
	n.s.rooms[0][2] = D
	n.s.rooms[1][1] = C
	n.s.rooms[1][2] = B
	n.s.rooms[2][1] = B
	n.s.rooms[2][2] = A
	n.s.rooms[3][1] = A
	n.s.rooms[3][2] = C

	g, dst := traverse(n)
	res := dijkstra(n, g, dst)

	fmt.Println("Part2:", res)
}

func dijkstra(n node, g graph, dst int) int {
	energy := make([]int, len(g))

	for i := range energy {
		energy[i] = common.MaxInt
	}
	energy[0] = 0

	visited := make(map[int]bool)

	pq := make(priorityQueue, 0)
	pq = append(pq, &stateWithEnergy{0, 0, 0})
	heap.Init(&pq)

	i := 1
	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*stateWithEnergy)
		if visited[u.id] {
			continue
		}
		visited[u.id] = true
		for k, v := range g[u.id].next {
			if energy[k] > energy[u.id]+v {
				energy[k] = energy[u.id] + v
				heap.Push(&pq, &stateWithEnergy{k, energy[k], i})
				i++
			}
		}
	}

	return energy[dst]
}

/// Graph&Node implementation

const (
	A int = (iota + 1) * 2
	B
	C
	D
)

type graph map[int]*node

var cost = map[int]int{A: 1, B: 10, C: 100, D: 1000}

func traverse(n node) (graph, int) {
	g := make(graph)
	g[n.id] = &n

	visited := make(map[int]bool)
	existed := make(map[state]int)
	existed[*(n.s)] = n.id
	dst := 0
	for {
		if len(g) == len(visited) {
			break
		}

		tmp := make(graph)
		for id, stableNode := range g {
			if visited[id] {
				continue
			}
			visited[id] = true
			move(stableNode, g, tmp, existed)
		}
		for k, v := range tmp {
			if v.end() {
				dst = k
			}
			g[k] = v
		}
	}

	return g, dst
}

// todo id dispaction
func move(n *node, g, tmp graph, existed map[state]int) {
	// first, traverse hallway to move something
	n.hallwayToRoom(g, tmp, existed)

	// then, check room move to hallway or room
	for amphipod, room := range n.s.rooms {
		// this room is empty, skip
		amphipod = amphipod*2 + 2

		if !n.roomCanMove(amphipod) {
			// this room is finished or wait other amphipod, don't move
			continue
		}

		for i := range room {
			if room[i] == 0 {
				continue
			}
			if room[i] == amphipod {
				// this amphipod can't move to other room
				n.roomToHallway(amphipod, i, g, tmp, existed)
			} else {
				// move to other room or hallway
				n.roomToHallway(amphipod, i, g, tmp, existed)
				n.roomToRoom(amphipod, i, room[i], g, tmp, existed)
			}
			break
		}
	}
}

func (n *node) hallwayToRoom(g, tmp graph, existed map[state]int) bool {
	flag := false
hallway:
	for i, amphipod := range n.s.hallway {
		if amphipod != 0 {
			// check room available?
			deep, available := n.roomAvailable(amphipod)
			if !available {
				continue
			}

			// check road connection
			var road []int
			if amphipod > i {
				road = n.s.hallway[i+1 : amphipod]
			} else {
				road = n.s.hallway[amphipod+1 : i]
			}

			for _, block := range road {
				// the road to room is blocked, can't move
				if block != 0 {
					continue hallway
				}
			}
			// new state
			cloned := n.clone()
			flag = true

			cloned.s.rooms[amphipod/2-1][deep] = amphipod
			cloned.s.hallway[i] = 0
			c := (len(road) + deep + 2) * cost[amphipod]

			if n.removeDuplicate(&cloned, existed, c) {
				continue hallway
			}

			cloned.id = len(g) + len(tmp)
			n.next[cloned.id] = c
			tmp[cloned.id] = &cloned
			existed[*cloned.s] = cloned.id
		}
	}

	return flag
}

func (n *node) roomToHallway(src, deep int, g, tmp graph, existed map[state]int) bool {
	flag := false
	for i := src - 1; i >= 0; i-- {
		// skip the space immediately outside any room
		if _, ok := cost[i]; ok {
			continue
		}

		if n.s.hallway[i] == 0 {
			cloned := n.clone()
			flag = true

			amphipod := n.s.rooms[src/2-1][deep]
			cloned.s.hallway[i] = amphipod
			cloned.s.rooms[src/2-1][deep] = 0
			c := (src - i + 1 + deep) * cost[amphipod]

			if n.removeDuplicate(&cloned, existed, c) {
				continue
			}

			cloned.id = len(g) + len(tmp)
			n.next[cloned.id] = c
			tmp[cloned.id] = &cloned
			existed[*cloned.s] = cloned.id
		} else {
			// blocked
			break
		}
	}

	for i := src + 1; i < len(n.s.hallway); i++ {
		// skip the space immediately outside any room
		if _, ok := cost[i]; ok {
			continue
		}

		if n.s.hallway[i] == 0 {
			cloned := n.clone()
			flag = true

			amphipod := n.s.rooms[src/2-1][deep]
			cloned.s.hallway[i] = amphipod
			cloned.s.rooms[src/2-1][deep] = 0
			c := (i - src + 1 + deep) * cost[amphipod]

			if n.removeDuplicate(&cloned, existed, c) {
				continue
			}

			cloned.id = len(g) + len(tmp)
			n.next[cloned.id] = c
			tmp[cloned.id] = &cloned
			existed[*cloned.s] = cloned.id
		} else {
			// blocked
			break
		}
	}

	return flag
}

func (n *node) roomToRoom(src, srcDeep, dst int, g, tmp graph, existed map[state]int) bool {
	flag := false
	dstDeep, available := n.roomAvailable(dst)
	if !available {
		return flag
	}
	var road []int
	if dst > src {
		road = n.s.hallway[src+1 : dst]
	} else {
		road = n.s.hallway[dst+1 : src]
	}
	for _, block := range road {
		// the road to room is blocked, can't move
		if block != 0 {
			return flag
		}
	}

	cloned := n.clone()
	flag = true

	amphipod := n.s.rooms[src/2-1][srcDeep]
	cloned.s.rooms[src/2-1][srcDeep] = 0
	cloned.s.rooms[dst/2-1][dstDeep] = amphipod
	c := (len(road) + 2 + srcDeep + dstDeep + 1) * cost[amphipod]

	if n.removeDuplicate(&cloned, existed, c) {
		return flag
	}

	cloned.id = len(g) + len(tmp)
	n.next[cloned.id] = c
	tmp[cloned.id] = &cloned
	existed[*cloned.s] = cloned.id

	return flag
}

func (n *node) roomAvailable(amphipod int) (int, bool) {
	roomIdx := amphipod/2 - 1
	if n.s.rooms[roomIdx][0] != 0 {
		return 4, false
	}

	for i := range n.s.rooms[roomIdx] {
		if n.s.rooms[roomIdx][i] != 0 && n.s.rooms[roomIdx][i] != amphipod {
			return 4, false
		}
	}

	res := 0
	for i := range n.s.rooms[roomIdx] {
		if n.s.rooms[roomIdx][i] == 0 {
			res = i
		}
	}

	return res, true
}

func (n *node) roomCanMove(amphipod int) bool {
	roomIdx := amphipod/2 - 1
	for i := range n.s.rooms[roomIdx] {
		if n.s.rooms[roomIdx][i] != 0 && n.s.rooms[roomIdx][i] != amphipod {
			return true
		}
	}

	return false
}

func (n *node) removeDuplicate(cloned *node, existed map[state]int, c int) bool {
	// find duplicated state
	// if exist, s.next point to duplicated state
	// instead of new cloned state
	if id, ok := existed[*cloned.s]; ok {
		n.next[id] = c
		return true
	}

	return false
}

func (n *node) clone() node {
	rooms := [4][4]int{}
	for k, v := range n.s.rooms {
		tmp := [4]int{}
		copy(tmp[:], v[:])
		rooms[k] = tmp
	}

	hallway := [11]int{}
	copy(hallway[:], n.s.hallway[:])
	return node{0, &state{rooms, hallway}, make(map[int]int)}
}

func (n *node) end() bool {
	for _, block := range n.s.hallway {
		if block != 0 {
			return false
		}
	}

	for amphipod, room := range n.s.rooms {
		wanted := amphipod*2 + 2
		for _, occupied := range room {
			if wanted != occupied {
				return false
			}
		}
	}

	return true
}

type node struct {
	id   int
	s    *state
	next map[int]int
}

func (n *node) String() string {
	return fmt.Sprintf("%d %v\n%s", n.id, n.next, n.s.String())
}

type state struct {
	rooms   [4][4]int
	hallway [11]int
}

var rev = map[int]byte{A: 'A', B: 'B', C: 'C', D: 'D'}

func (s *state) String() string {
	res := [5][11]byte{}
	for i := range res[0] {
		if s.hallway[i] == 0 {
			res[0][i] = '.'
		} else {
			res[0][i] = rev[s.hallway[i]]
		}
	}
	for i := 1; i < 5; i++ {
		res[i][0] = ' '
		res[i][10] = ' '
		for j := 1; j < 10; j++ {
			if j%2 == 1 {
				res[i][j] = ' '
			} else {
				if s.rooms[j/2-1][i-1] == 0 {
					res[i][j] = ' '
				} else {
					res[i][j] = rev[s.rooms[j/2-1][i-1]]
				}
			}
		}
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		string(res[0][:]), string(res[1][:]),
		string(res[2][:]), string(res[3][:]), string(res[4][:]))
}

func initialNodeFromStr(s string) node {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Scan()

	scanner.Scan()
	hallway := [11]int{}

	amphipod := map[rune]int{'A': A, 'B': B, 'C': C, 'D': D}

	scanner.Scan()
	rooms := [4][4]int{}
	for i, ch := range scanner.Text()[1:] {
		if ch != '#' {
			rooms[i/2-1][0] = amphipod[ch]
		}
	}

	scanner.Scan()
	for i, ch := range scanner.Text()[1:] {
		if ch != '#' && ch != ' ' {
			rooms[i/2-1][1] = amphipod[ch]
		}
	}

	return node{0, &state{rooms, hallway}, make(map[int]int)}
}

/// Priority queue implementation

type stateWithEnergy struct {
	id, energy, index int
}

type priorityQueue []*stateWithEnergy

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].energy < pq[j].energy
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
	item := x.(*stateWithEnergy)
	item.index = n
	*pq = append(*pq, item)
}
