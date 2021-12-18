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
	commands := part1(s)
	part2(commands)
}

func part1(s string) []command {
	scanner := bufio.NewScanner(strings.NewReader(s))
	commands := make([]command, 0)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		direction, distance := parts[0], parts[1]
		val, err := strconv.Atoi(distance)
		common.Check(err)
		var command command
		command.val = val

		switch direction {
		case "forward":
			command.direction = forward
		case "down":
			command.direction = down
		case "up":
			command.direction = up
		default:
			panic("unreachable")
		}

		commands = append(commands, command)
	}

	var x, y int
	for _, cmd := range commands {
		switch cmd.direction {
		case forward:
			x += cmd.val
		case down:
			y += cmd.val
		case up:
			y -= cmd.val
		}
	}

	fmt.Println("Part1:", x*y)

	return commands
}

func part2(commands []command) {
	var x, y, aim int
	for _, cmd := range commands {
		switch cmd.direction {
		case forward:
			x += cmd.val
			y += aim * cmd.val
		case down:
			aim += cmd.val
		case up:
			aim -= cmd.val
		}
	}

	fmt.Println("Part2:", x*y)
}

type direction int

const (
	forward direction = iota
	down
	up
)

type command struct {
	direction direction
	val       int
}
