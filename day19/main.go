package main

import (
	"bufio"
	"common"
	"fmt"
	"strings"
)

var rs = generateRotations()
var index, positive = generateIndex()

func main() {
	s := common.Load()
	scans := scanFromStr(s)
	infos := part1(scans)
	part2(infos)
}

func part1(scans []scan) []map[int]scanInfo {
	lookupTable := make([]map[coordinate]bool, 0)
	infos := make([]map[int]scanInfo, 0)
	beacons := make(map[coordinate]bool)
	for _, scan := range scans {
		tmp := make(map[coordinate]bool)
		for _, c := range scan {
			tmp[c] = true
		}
		lookupTable = append(lookupTable, tmp)
		infos = append(infos, make(map[int]scanInfo))
	}

	for i := range scans {
		for j := range scans {
			if j == i {
				continue
			}
			info := overlap(scans[i], scans[j], lookupTable[j])
			if info.relativeRotation >= 0 {
				infos[i][j] = info
			}
		}
	}

	// infos is just like adjacency list
	// dfs to traverse all scanners
	res := traverse(infos, beacons, scans)
	fmt.Println("Part1:", res)

	return infos
}

func part2(infos []map[int]scanInfo) {
	visited := make(map[int]bool)
	visited[0] = true
	scanners := make(scan, 0)
	// add zero
	scanners = append(scanners, coordinate{})
	for k, v := range infos[0] {
		visited[k] = true
		scanners = append(scanners, v.relativeCoordinate)
		for _, c := range findScanners(visited, infos, k) {
			scanners = append(scanners, relativeCoordinate(v.relativeRotation, v.relativeCoordinate, c))
		}
	}

	largest := 0
	for i, c1 := range scanners {
		for _, c2 := range scanners[i+1:] {
			largest = common.IntMax(distance(c1, c2), largest)
		}
	}

	fmt.Println("Part2:", largest)
}

func distance(c1, c2 coordinate) int {
	return common.IntAbs(c1.x-c2.x) + common.IntAbs(c1.y-c2.y) + common.IntAbs(c1.z-c2.z)
}

func findScanners(visited map[int]bool, infos []map[int]scanInfo, self int) scan {
	scanners := make(scan, 0)

	for k, v := range infos[self] {
		if visited[k] {
			continue
		}
		visited[k] = true
		scanners = append(scanners, v.relativeCoordinate)
		for _, c := range findScanners(visited, infos, k) {
			scanners = append(scanners, relativeCoordinate(v.relativeRotation, v.relativeCoordinate, c))
		}
	}

	return scanners
}

func traverse(infos []map[int]scanInfo, beacons map[coordinate]bool, scans []scan) int {
	visited := make(map[int]bool)
	for _, c := range scans[0] {
		beacons[c] = true
	}
	visited[0] = true
	for k, v := range infos[0] {
		visited[k] = true
		for k := range dfs(visited, infos, scans, k) {
			beacons[relativeCoordinate(v.relativeRotation, v.relativeCoordinate, k)] = true
		}
	}

	return len(beacons)
}

func dfs(visited map[int]bool, infos []map[int]scanInfo, scans []scan, self int) map[coordinate]bool {
	beacons := make(map[coordinate]bool)
	for _, c := range scans[self] {
		beacons[c] = true
	}

	for k, v := range infos[self] {
		if visited[k] {
			continue
		}
		visited[k] = true
		for k := range dfs(visited, infos, scans, k) {
			beacons[relativeCoordinate(v.relativeRotation, v.relativeCoordinate, k)] = true
		}
	}

	return beacons
}

func overlap(first, second scan, lookupTable map[coordinate]bool) scanInfo {
	for i := range rs {
		for f := range first {
			for s := range second {
				used := make(map[int]bool)
				base := compute(i, second[s], first[f])
				count := 1
			inner:
				for {
					if count >= 12 {
						return scanInfo{i, base}
					}
					for idx, c := range first {
						if idx == f || used[idx] {
							continue
						}
						if verify(i, base, c, lookupTable) {
							count++
							used[idx] = true
							continue inner
						}
					}
					break
				}
			}
		}
	}
	return scanInfo{-1, coordinate{}}
}

func compute(r int, c1, c2 coordinate) coordinate {
	res := make([]int, 3)
	c1Val := make([]int, 0)
	c1Val = append(c1Val, c1.x, c1.y, c1.z)
	c2Val := make([]int, 0)
	c2Val = append(c2Val, c2.x, c2.y, c2.z)
	for i, variable := range index[r] {
		res[variable] = c2Val[variable]*positive[r][variable] - c1Val[i]
		res[variable] *= positive[r][variable]
	}

	return coordinate{res[0], res[1], res[2]}
}

func verify(r int, base, c coordinate, lookup map[coordinate]bool) bool {
	res := make([]int, 3)
	cVal := make([]int, 0)
	cVal = append(cVal, c.x, c.y, c.z)
	baseVal := make([]int, 0)
	baseVal = append(baseVal, base.x, base.y, base.z)
	for i, variable := range index[r] {
		res[i] = positive[r][variable] * (cVal[variable] - baseVal[variable])
	}
	return lookup[coordinate{res[0], res[1], res[2]}]
}

func relativeCoordinate(r int, base, c coordinate) coordinate {
	res := make([]int, 3)
	cVal := make([]int, 0)
	cVal = append(cVal, c.x, c.y, c.z)
	baseVal := make([]int, 0)
	baseVal = append(baseVal, base.x, base.y, base.z)
	for i, variable := range index[r] {
		res[variable] = cVal[i] + positive[r][variable]*baseVal[variable]
		res[variable] *= positive[r][variable]
	}
	return coordinate{res[0], res[1], res[2]}
}

type rotation [][]int

func generateRotations() []rotation {
	r := make(rotation, 0)
	for i := 0; i < 3; i++ {
		tmp := make([]int, 3)
		r = append(r, tmp)
	}
	rs := make([]rotation, 0)
	for x := 0; x < 3; x++ {
		if x != 0 {
			r[0][x-1] = 0
		}
		r[0][x] = 1
		rs = append(rs, r.rotate()...)
		r[0][x] = -1
		rs = append(rs, r.rotate()...)
	}

	return rs
}

func generateIndex() ([][]int, [][]int) {
	index := make([][]int, 0)
	positive := make([][]int, 0)
	for r := range rs {
		idx := make([]int, 3)
		pos := make([]int, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if rs[r][i][j] != 0 {
					idx[i] = j
					pos[i] = rs[r][i][j]
				}
			}
		}
		index = append(index, idx)
		positive = append(positive, pos)
	}

	return index, positive
}

func (r *rotation) clone() rotation {
	replica := make(rotation, len(*r))
	for i := range replica {
		replica[i] = append(replica[i], (*r)[i]...)
	}
	return replica
}

func (r *rotation) rotate() []rotation {
	x := 0
	for idx := range *r {
		if (*r)[0][idx] != 0 {
			x = idx
			break
		}
	}

	rs := make([]rotation, 0)

	for y := 0; y < 3; y++ {
		if y == x {
			continue
		}
		for i := -1; i <= 1; i += 2 {
			replica := r.clone()
			replica[1][y] = i

			for z := 0; z < 3; z++ {
				if z == x || z == y {
					continue
				}
				replica[2][z] = replica[0][x] * replica[1][y]
				if x > y {
					replica[2][z] = -replica[2][z]
				}
				if x-y == 2 || x-y == -2 {
					replica[2][z] = -replica[2][z]
				}
			}
			rs = append(rs, replica)
		}
	}

	return rs
}

func (r rotation) String() string {
	for _, row := range r {
		fmt.Println(row)
	}
	return ""
}

type scan []coordinate

type scanInfo struct {
	relativeRotation   int
	relativeCoordinate coordinate
}

func scanFromStr(s string) []scan {
	scans := make([]scan, 0)

	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Scan()
	tmp := make(scan, 0)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			scans = append(scans, tmp)
			tmp = make(scan, 0)
			scanner.Scan()
			continue
		}
		c := coordinate{}
		fmt.Sscanf(scanner.Text(), "%d,%d,%d", &c.x, &c.y, &c.z)
		tmp = append(tmp, c)
	}
	scans = append(scans, tmp)

	return scans
}

type coordinate struct {
	x, y, z int
}
