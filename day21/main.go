package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	s := common.Load()
	g, initialState := gameFromStr(s)
	part1(&g)
	part2(initialState)
}

func part1(g *game) {
	g.play()
	fmt.Println("Part1:", g.rollTimes*g.players[1-g.curr].score)
}

func part2(initialState state) {
	totalState := make(map[state]int64)
	totalState[initialState] = 1

	wins := [2]int64{}

	for {
		tmpState := make(map[state]int64)

		for s, n := range totalState {
			for d1 := 1; d1 <= 3; d1++ {
				for d2 := 1; d2 <= 3; d2++ {
					for d3 := 1; d3 <= 3; d3++ {
						copied := s

						copied.pos[s.curr] =
							(copied.pos[copied.curr]+d1+d2+d3-1)%10 + 1
						copied.score[s.curr] += copied.pos[s.curr]

						if copied.score[s.curr] >= 21 {
							wins[copied.curr] += n
						} else {
							copied.curr = 1 - copied.curr
							tmpState[copied] += n
						}
					}
				}
			}
		}

		totalState = tmpState
		if len(totalState) == 0 {
			break
		}
	}

	winner := wins[0]
	if wins[1] > winner {
		winner = wins[1]
	}

	fmt.Println("Part2:", winner)
}

type game struct {
	players   [2]player
	rollTimes int
	curr      int
}

type state struct {
	pos   [2]int
	score [2]int
	curr  int
}

type player struct {
	pos, score int
}

func (g *game) play() {
	for {
		for i := 0; i < 3; i++ {
			pos := g.players[g.curr].pos
			g.players[g.curr].pos = (pos%10+(g.rollTimes%100+1)%10-1)%10 + 1
			g.rollTimes++
		}

		g.players[g.curr].score += g.players[g.curr].pos

		if g.players[g.curr].score >= 1000 {
			break
		}
		g.curr = 1 - g.curr
	}
}

func gameFromStr(s string) (game, state) {
	var players [2]player
	var initialState state
	for i, s := range strings.Split(s, "\n") {
		pos := 0
		drop := 0
		fmt.Sscanf(s, "Player %d starting position: %d", &drop, &pos)
		players[i].pos = pos
		initialState.pos[i] = pos
	}

	return game{players, 0, 0}, initialState
}
