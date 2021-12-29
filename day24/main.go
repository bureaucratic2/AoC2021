package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"
)

// program from input can be split into 14 sections by "inp w"
// each section excute like:
// x = z
// x %= 26
// z /= args[0]
// x += args[1]
// if x == input {
// 	  z = z
// } else {
//    z *= 26
//    z += input + args[2]
// }

func main() {
	s := common.Load()
	args := argsFromStr(s)
	e := getEquation(&args)
	v := solveEquation(e)
	part1(&v)
	part2(&v)
}

func part1(v *[]variable) {
	s := []int{}
	for i := range *v {
		s = append(s, (*v)[i].max)
	}
	fmt.Println("Part1:", s)
}

func part2(v *[]variable) {
	s := []int{}
	for i := range *v {
		s = append(s, (*v)[i].min)
	}
	fmt.Println("Part2:", s)
}

type equation struct {
	lVar, rVar variable
	scalar     int
}

type variable struct {
	idx, min, max int
}

func getEquation(args *[3][14]int) []equation {
	s := make(common.Stack, 0)
	e := make([]equation, 0)
	for i := range args[0] {
		if args[0][i] == 1 {
			s.Push(i)
		} else {
			if rVal, err := s.Pop(); err == nil {
				res := args[2][rVal.(int)] + args[1][i]
				if res > 0 {
					fmt.Printf("x%d = x%d + %d\n", i, rVal, res)
				} else {
					fmt.Printf("x%d = x%d - %d\n", i, rVal, -res)
				}
				e = append(e, equation{variable{i, 1, 9}, variable{rVal.(int), 1, 9}, res})
			}
		}
	}
	return e
}

func solveEquation(e []equation) []variable {
	vars := make([]variable, len(e)*2)
	for _, e := range e {
		vars[e.lVar.idx] = e.lVar
		vars[e.rVar.idx] = e.rVar
	}

	for _, e := range e {
		if e.scalar >= 0 {
			vars[e.lVar.idx].min += e.scalar
			vars[e.rVar.idx].max -= e.scalar
		} else {
			vars[e.lVar.idx].max += e.scalar
			vars[e.rVar.idx].min -= e.scalar
		}
	}

	fmt.Println()
	fmt.Println("input variable range: ")
	for i, v := range vars {
		fmt.Printf("%d <= x%d <= %d\n", v.min, i, v.max)
	}
	return vars
}

func argsFromStr(s string) [3][14]int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	args := [3][14]int{}

	idx := [3]int{}

	for scanner.Scan() {
		// remember to fliter confusing lines like "add x z" and "add y w"
		if strings.HasPrefix(scanner.Text(), "div z ") {
			fmt.Sscanf(scanner.Text(), "div z %d", &args[0][idx[0]])
			idx[0]++
			continue
		}

		if strings.HasPrefix(scanner.Text(), "add x ") {
			if scanner.Text() == "add x z" {
				continue
			}
			fmt.Sscanf(scanner.Text(), "add x %d", &args[1][idx[1]])
			idx[1]++
			continue
		}

		if strings.HasPrefix(scanner.Text(), "add y w") {
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "add y %d", &args[2][idx[2]])
			idx[2]++
			continue
		}
	}
	return args
}
