package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"
)

func main() {
	s := common.Load()
	ps := pairsFromStr(s)
	part1(ps)
	// reconstruct state
	ps = pairsFromStr(s)
	part2(ps)
}

func part1(pairs []*pair) {
	for i := 1; i < len(pairs); i++ {
		pairs[0] = pairs[0].add(pairs[i])
	}
	fmt.Println("Part1:", pairs[0].magnitude())
}

func part2(pairs []*pair) {
	largest := common.MinInt
	for i := 1; i < len(pairs); i++ {
		for j := 1; j < len(pairs); j++ {
			if i == j {
				continue
			}
			res := pairs[i].clone().add(pairs[j].clone()).magnitude()
			largest = common.IntMax(largest, res)
		}
	}
	fmt.Println("Part2:", largest)
}

type pair struct {
	val, deep int
	children  []*pair
}

func (p *pair) clone() *pair {
	replica := pair{}
	replica.val = p.val
	replica.deep = p.deep

	if p.val == -1 {
		for i := 0; i < 2; i++ {
			replica.children = append(replica.children, p.children[i].clone())
		}
	}

	return &replica
}

func (p *pair) magnitude() int {
	if p.val != -1 {
		return p.val
	}

	return 3*p.children[0].magnitude() + 2*p.children[1].magnitude()
}

func (p *pair) add(other *pair) *pair {
	res := pair{}
	res.val = -1
	res.deep = 1
	p.deeper()
	other.deeper()
	res.children = make([]*pair, 0)
	res.children = append(res.children, p, other)
	res.reduce()

	return &res
}

func (p *pair) deeper() {
	if p.val != -1 {
		return
	}
	p.deep++
	p.children[0].deeper()
	p.children[1].deeper()
}

func (p *pair) reduce() {
	for {
		reduced := false

		s := make(common.Stack, 0)
		reduced = p.traversal(&s, true)
		if reduced {
			continue
		}
		reduced = p.traversal(&s, false)
		if reduced {
			continue
		}

		if !reduced {
			return
		}
	}
}

func (p *pair) traversal(s *common.Stack, first bool) bool {
	if p == nil {
		return false
	}

	if p.val != -1 {
		if p.val >= 10 {
			if first {
				return false
			}
			p.split(s)
			return true
		} else {
			return false
		}
	}

	if p.deep > 4 && p.children[0].val != -1 && p.children[1].val != -1 {
		s.Push(p)
		explode(s, p.children[0].val, p.children[1].val, len(*s))
		s.Pop()
		p.val = 0
		p.children = nil

		return true
	}

	s.Push(p)
	defer s.Pop()
	if p.children[0].traversal(s, first) {
		return true
	}
	if p.children[1].traversal(s, first) {
		return true
	}
	return false
}

func (p *pair) split(s *common.Stack) {
	p.deep = (*s)[len(*s)-1].(*pair).deep + 1
	p.children = make([]*pair, 0)
	if p.val%2 == 1 {
		p.children = append(p.children, &pair{p.val / 2, p.deep, nil}, &pair{p.val/2 + 1, p.deep, nil})
	} else {
		p.children = append(p.children, &pair{p.val / 2, p.deep, nil}, &pair{p.val / 2, p.deep, nil})
	}
	p.val = -1
}

func explode(s *common.Stack, left int, right int, l int) {
	if l <= 1 {
		return
	}
	parent := (*s)[l-2].(*pair)
	for i := 0; i < 2; i++ {
		if i == 0 && left == -1 {
			continue
		}
		if i == 1 && right == -1 {
			continue
		}
		if parent.children[i] == (*s)[l-1].(*pair) {
			continue
		}
		switch i {
		case 0:
			parent.children[0].directionMost(1).val += left
			left = -1
		case 1:
			parent.children[1].directionMost(0).val += right
			right = -1
		}
	}
	if left != -1 || right != -1 {
		explode(s, left, right, l-1)
	}
}

// 0-leftmost,1-rightmost
func (p *pair) directionMost(direction int) *pair {
	if p.val != -1 {
		return p
	} else {
		return p.children[direction].directionMost(direction)
	}
}

func (p pair) String() string {
	if p.val != -1 {
		return fmt.Sprintf("%v", p.val)
	} else {
		return fmt.Sprintf("[%v, %v]", *p.children[0], *p.children[1])
	}
}

func pairFromStr(s string) *pair {
	stack := make(common.Stack, 0)
	deep := 1
	for _, ch := range s {
		if ch == ',' {
			continue
		}
		if ch == '[' {
			p := pair{}
			p.val = -1
			p.deep = deep
			p.children = make([]*pair, 0)
			stack.Push(&p)
			deep++
		} else if ch >= '0' && ch <= '9' {
			p := pair{}
			p.val = int(ch - '0')
			stack.Push(&p)
		} else {
			right, err := stack.Pop()
			common.Check(err)
			left, err := stack.Pop()
			common.Check(err)
			parent, err := stack.Pop()
			common.Check(err)

			p := parent.(*pair)
			l := left.(*pair)
			r := right.(*pair)

			l.deep = p.deep + 1
			r.deep = p.deep + 1
			p.children = append(p.children, left.(*pair))
			p.children = append(p.children, right.(*pair))
			stack.Push(p)
			deep--
		}
	}

	return stack[0].(*pair)
}

func pairsFromStr(s string) []*pair {
	pairs := make([]*pair, 0)
	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		p := pairFromStr(scanner.Text())
		p.reduce()
		pairs = append(pairs, p)
	}

	return pairs
}
