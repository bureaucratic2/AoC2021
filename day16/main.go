package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	s := common.Load()
	msg := packetsFromStr(s)
	part1(msg)
	part2(msg)
}

func part1(msg message) {
	pkt := parsePacket(msg)
	fmt.Println("Part1:", pkt.versionSum())
}

func part2(msg message) {
	pkt := parsePacket(msg)
	fmt.Println("Part2:", pkt.evaluate())
}

type packet struct {
	version int
	typeID  int
	len     int
	val     int
	sub     []packet
}

func (pkt *packet) evaluate() int {
	res := 0
	switch pkt.typeID {
	case 0:
		res = 0
		for i := range pkt.sub {
			res += pkt.sub[i].evaluate()
		}
	case 1:
		res = 1
		for i := range pkt.sub {
			res *= pkt.sub[i].evaluate()
		}
	case 2:
		res = common.MaxInt
		for i := range pkt.sub {
			res = common.IntMin(pkt.sub[i].evaluate(), res)
		}
	case 3:
		res = common.MinInt
		for i := range pkt.sub {
			res = common.IntMax(pkt.sub[i].evaluate(), res)
		}
	case 4:
		res = pkt.val
	case 5:
		if len(pkt.sub) != 2 {
			panic("sub packets number conflict for comparison")
		}
		if pkt.sub[0].evaluate() > pkt.sub[1].evaluate() {
			res = 1
		} else {
			res = 0
		}
	case 6:
		if len(pkt.sub) != 2 {
			panic("sub packets number conflict for comparison")
		}
		if pkt.sub[0].evaluate() < pkt.sub[1].evaluate() {
			res = 1
		} else {
			res = 0
		}
	case 7:
		if len(pkt.sub) != 2 {
			panic("sub packets number conflict for comparison")
		}
		if pkt.sub[0].evaluate() == pkt.sub[1].evaluate() {
			res = 1
		} else {
			res = 0
		}
	}
	pkt.val = res
	return res
}

func (pkt *packet) versionSum() int {
	ver := 0
	ver += pkt.version
	for i := range pkt.sub {
		ver += pkt.sub[i].versionSum()
	}
	return ver
}

func (pkt packet) String() string {
	if pkt.sub == nil {
		return fmt.Sprintf("[ver: %v] [type ID: %v] [%v] [%v]", pkt.version, pkt.typeID, pkt.len, pkt.val)
	} else {
		s := []string{fmt.Sprintf("[ver: %v] [type ID: %v] [%v] [%v]{", pkt.version, pkt.typeID, pkt.len, pkt.val)}
		for _, pkt := range pkt.sub {
			s = append(s, pkt.String())
		}
		return fmt.Sprintf("%v\n}", strings.Join(s, "\n"))
	}
}

func parsePacket(msg message) packet {
	pkt := packet{}
	pkt.version = binaryToDecimal(msg[:3])
	pkt.typeID = binaryToDecimal(msg[3:6])
	if pkt.typeID == 4 {
		// parse literal
		pkt.len, pkt.val = parseLiteral(msg[6:])
		pkt.len += 6
		return pkt
	} else {
		// parse operator
		pkt.sub = make([]packet, 0)
		pkt.len += 6
		parseOperator(&pkt, msg[6:])
		// todo
		return pkt
	}
}

func parseLiteral(msg message) (int, int) {
	start := 0
	val := 0
	for {
		val <<= 4
		val += binaryToDecimal(msg[start+1 : start+5])

		if msg[start] == 0 {
			break
		}
		start += 5
	}
	return start + 5, val
}

func parseOperator(pkt *packet, msg message) {
	var start, end int
	if msg[start] >= 2 {
		panic("length type ID bigger than 1")
	}
	pkt.len++
	if msg[start] == 0 {
		start++
		pkt.len += 15
		end = start + 15
		pktsLen := binaryToDecimal(msg[start:end])
		start = end
		end += pktsLen
		for start < end {
			sub := parsePacket(msg[start:])
			start += sub.len
			pkt.sub = append(pkt.sub, sub)
			pkt.len += sub.len
			if end-start < 4 {
				// padding zero, ignore them
				break
			}
		}
	} else {
		start++
		pkt.len += 11
		end = start + 11
		pktsNum := binaryToDecimal(msg[start:end])
		start = end
		for i := 0; i < pktsNum; i++ {
			sub := parsePacket(msg[start:])
			start += sub.len
			pkt.sub = append(pkt.sub, sub)
			pkt.len += sub.len
		}
	}
}

func binaryToDecimal(arr []byte) int {
	res := 0
	for i := 0; i < len(arr); i++ {
		res <<= 1
		res += int(arr[i])
	}
	return res
}

type message []byte

func packetsFromStr(s string) message {
	msg := make(message, 0)

	for i := range s {
		n := s[i]
		if n >= 65 {
			n -= '7'
		} else {
			n -= '0'
		}
		l := len(msg)
		msg = append(msg, 0, 0, 0, 0)
		for i := 3; i >= 0; i-- {
			msg[l+i] = n % 2
			n >>= 1
		}
	}

	return msg
}
