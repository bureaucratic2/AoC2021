package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"
)

func main() {
	s := common.Load()
	part1(s)
	part2(s)
}

func part1(s string) {
	p := pixelMapFromStr(s)
	p.enhance()
	p.enhance()

	fmt.Println("Part1:", len(p.pixels[p.curr]))
}

func part2(s string) {
	p := pixelMapFromStr(s)
	for i := 0; i < 50; i++ {
		p.enhance()
	}

	fmt.Println("Part2:", len(p.pixels[p.curr]))
}

type coordinate struct {
	x, y int
}

func (c *coordinate) add(adder int) {
	c.x += adder
	c.y += adder
}

type pixelMap struct {
	algo     []bool
	pixels   []map[coordinate]bool
	curr     int
	min, max coordinate
	outside  bool
	rev      bool
	ch       chan bool
}

func (p *pixelMap) enhance() {
	p.min.add(-1)
	p.max.add(1)

	curr := p.curr

	for y := p.min.y; y <= p.max.y; y++ {
		for x := p.min.x; x <= p.max.x; x++ {
			idx := 0
			p.center(x, y)
		outer:
			for {
				select {
				case pixel := <-p.ch:
					idx <<= 1
					if pixel {
						idx += 1
					}
				default:
					break outer
				}
			}

			p.pixels[1-curr][coordinate{x, y}] = p.algo[idx]
		}
	}

	if p.rev {
		p.outside = !p.outside
	}
	p.curr = 1 - curr
	p.compact()
}

func (p *pixelMap) center(x, y int) {
	for j := -1; j < 2; j++ {
		for i := -1; i < 2; i++ {
			c := coordinate{x + i, y + j}

			if c.x <= p.min.x || c.x >= p.max.x || c.y <= p.min.y || c.y >= p.max.y {
				p.ch <- p.outside
			} else {
				p.ch <- p.pixels[p.curr][c]
			}
		}
	}
}

func (p *pixelMap) compact() {
	for k, v := range p.pixels[p.curr] {
		if !v {
			delete(p.pixels[p.curr], k)
		}
	}
}

func (p *pixelMap) String() string {
	ss := make([]string, 0)
	tmp := make([]byte, p.max.x-p.min.x+1)
	for y := p.min.y; y <= p.max.y; y++ {
		for x := p.min.x; x <= p.max.x; x++ {
			if p.pixels[p.curr][coordinate{x, y}] {
				tmp[x-p.min.x] = '#'
			} else {
				tmp[x-p.min.x] = '.'
			}
		}
		ss = append(ss, string(tmp))
	}

	return strings.Join(ss, "\n")
}

func pixelMapFromStr(s string) pixelMap {
	scanner := bufio.NewScanner(strings.NewReader(s))

	scanner.Scan()
	algo := make([]bool, 0)

	for _, ch := range scanner.Text() {
		if ch == '#' {
			algo = append(algo, true)
		} else {
			algo = append(algo, false)
		}
	}

	rev := false
	if algo[0] && !algo[len(algo)-1] {
		rev = true
	}

	// skip blank line
	scanner.Scan()

	pixels := make([]map[coordinate]bool, 0)
	pixels = append(pixels, make(map[coordinate]bool))
	pixels = append(pixels, make(map[coordinate]bool))

	y := 0
	xLen := 0
	for scanner.Scan() {
		if xLen == 0 {
			xLen = len(scanner.Text()) - 1
		}
		for x, ch := range scanner.Text() {
			if ch == '#' {
				pixels[0][coordinate{x, y}] = true
			}
		}
		y++
	}

	return pixelMap{algo, pixels, 0, coordinate{0, 0}, coordinate{xLen, y - 1}, false, rev, make(chan bool, 9)}
}
