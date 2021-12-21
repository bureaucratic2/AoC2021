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
	p, is := paperAndInstructionsFromStr(s)
	part1(p, is)
	// reconstruct state, avoid using the same memory
	p, is = paperAndInstructionsFromStr(s)
	part2(p, is)
}

func part1(p paper, is []instruction) {
	i := is[0]

	step(&p, i)

	count := 0
	for _, row := range p {
		for _, mark := range row {
			if mark {
				count++
			}
		}
	}

	fmt.Println("Part1:", count)
}

func part2(p paper, is []instruction) {
	for _, i := range is {
		step(&p, i)
	}

	fmt.Println("Part2:")
	p.toStr()
}

func step(p *paper, i instruction) {
	if i.foldX {
		for y := range *p {
			for x, mark := range (*p)[y] {
				if x <= i.idx || !mark {
					continue
				}
				(*p)[y][2*i.idx-x] = mark
			}
			(*p)[y] = (*p)[y][:i.idx]
		}
	} else {
		for y := range *p {
			if y <= i.idx {
				continue
			}
			for x, mark := range (*p)[y] {
				if !mark {
					continue
				}
				(*p)[2*i.idx-y][x] = mark
			}
		}
		*p = (*p)[:i.idx]
	}
}

type paper [][]bool

func (p *paper) toStr() {
	for y := range *p {
		s := make([]byte, 0)
		for x := range (*p)[y] {
			if (*p)[y][x] {
				s = append(s, '#')
			} else {
				s = append(s, '.')
			}
		}
		fmt.Println(string(s))
	}
}

type coordinate struct {
	x int
	y int
}

type instruction struct {
	foldX bool
	idx   int
}

func paperAndInstructionsFromStr(s string) (paper, []instruction) {
	scanner := bufio.NewScanner(strings.NewReader(s))

	cs := make([]coordinate, 0)
	maxX := 0
	maxY := 0
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		}
		c := strings.Split(scanner.Text(), ",")
		x, err := strconv.Atoi(c[0])
		common.Check(err)
		y, err := strconv.Atoi(c[1])
		common.Check(err)
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		cs = append(cs, coordinate{x, y})
	}

	p := make(paper, 0)
	for i := 0; i <= maxY; i++ {
		p = append(p, make([]bool, maxX+1))
	}

	for _, c := range cs {
		p[c.y][c.x] = true
	}

	is := make([]instruction, 0)
	for scanner.Scan() {
		var i instruction
		sub := strings.Split(strings.Trim(scanner.Text(), "fold ang"), "=")
		if sub[0] == "x" {
			i.foldX = true
		} else {
			i.foldX = false
		}
		idx, err := strconv.Atoi(sub[1])
		common.Check(err)
		i.idx = idx
		is = append(is, i)
	}

	return p, is
}
