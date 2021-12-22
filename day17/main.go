package main

import (
	"common"
	"fmt"
	"math"
)

const (
	onComming int = iota
	yTooFast
	outOfBound
)

func main() {
	s := common.Load()
	// lowerBound := math.Sqrt(2*float64(5)) - 1
	// fmt.Println(lowerBound)
	// fmt.Println(int(lowerBound))
	t := targetFromStr(s)
	part1(t)
	part2(t)
}

func part1(t target) {
	highest, _ := lauch(t, false)
	fmt.Println("Part1:", highest)
}

func part2(t target) {
	_, count := lauch(t, true)
	fmt.Println("Part2:", count)
}

func lauch(t target, duplicate bool) (int, int) {
	lowerBound := int(math.Sqrt(2*float64(t.minX))) + 1
	upperBound := t.maxX + 1
	endCondition := 0
	existFallfree := false

	highest := 0
	count := 0

	// find possible x that could do free-fall motion
	freeFallX := make(map[int]bool)
	for x := lowerBound; x < upperBound; x++ {
		res := x * (x + 1) / 2
		if res <= t.maxX && res >= t.minX {
			freeFallX[x] = true
		}
	}

	if len(freeFallX) > 0 {
		existFallfree = true
	}

	// first, consider projectile motion
	// assume y velocity always start from 0 or from negative number if we need
	// to count total options
	y := 0
	if duplicate {
		y = t.minY
	}
outer:
	for ; ; y++ {
		found := false
	inner:
		for x := lowerBound; x < upperBound; x++ {
			p := probe{coordinate{0, 0}, velocity{x, y}}
			for !t.within(p.c) {
				switch t.outOfBound(p.c) {
				case onComming:
					p.step()
				case yTooFast:
					if y < 0 {
						// actually x is too fast
						continue outer
					}
					if endCondition == 0 {
						if p.v.x > 0 && !found {
							endCondition = x
							continue outer
						}
						continue inner
					} else {
						if x >= endCondition {
							// upper bound of y has found
							break outer
						} else {
							continue outer
						}
					}

				case outOfBound:
					continue inner
				}
			}
			highest = common.IntMax(highest, y*(y+1)/2)
			found = true

			if duplicate {
				if y < 0 {
					count++
					// don't count free fall point here, it's incomplete
				} else if !freeFallX[x] {
					count++
				}
			} else {
				break
			}
		}
	}

	// then, consider free-fall motion
	if existFallfree {
		freeFallCount := 0
		// y velocity start from -1
		for y := -1; y >= t.minY; y-- {
			p := probe{coordinate{t.minX, 0}, velocity{0, y}}
			for {
				p.step()
				if t.within(p.c) {
					initialY := -y - 1
					highest = common.IntMax(highest, initialY*(initialY+1)/2)
					freeFallCount++
					break
				}
				if p.c.y < t.minY {
					break
				}
			}
		}
		// count free fall point here
		count += len(freeFallX) * freeFallCount
	}
	return highest, count
}

type target struct {
	minX, maxX, minY, maxY int
}

func (t *target) within(c coordinate) bool {
	if c.x <= t.maxX && c.x >= t.minX && c.y <= t.maxY && c.y >= t.minY {
		return true
	}
	return false
}

func (t *target) outOfBound(c coordinate) int {
	if c.y > t.maxY && c.x >= t.maxX {
		return yTooFast
	}
	if c.y <= t.minY || c.x >= t.maxX {
		return outOfBound
	}
	return onComming
}

type probe struct {
	c coordinate
	v velocity
}

func (p *probe) step() {
	p.c.x += p.v.x
	p.c.y += p.v.y
	p.v.step()
}

type coordinate struct {
	x, y int
}

type velocity struct {
	x, y int
}

func (v *velocity) step() {
	v.y--
	if v.x == 0 {
		return
	}
	if v.x > 0 {
		v.x--
	} else {
		v.x++
	}
}

func targetFromStr(s string) target {
	t := target{}
	fmt.Sscanf(s, "target area: x=%d..%d, y=%d..%d", &t.minX, &t.maxX, &t.minY, &t.maxY)
	return t
}
