package common

import (
	"errors"
	"os"
	"path"
	"unicode"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

const MaxUint64 = ^uint64(0)
const MinUint64 = 0
const MaxInt64 = int64(MaxUint64 >> 1)
const MinInt64 = -MaxInt64 - 1

type Pair struct {
	A, B interface{}
}

func Load() string {
	dir, err := os.Getwd()
	Check(err)
	dir = path.Join(dir, "input")

	data, err := os.ReadFile(dir)
	Check(err)

	return string(data)
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IntMin(i int, j int) int {
	if i < j {
		return i
	} else {
		return j
	}
}

func IntMax(i int, j int) int {
	if i < j {
		return j
	} else {
		return i
	}
}

func IntAbs(i int) int {
	if i < 0 {
		i = -i
	}
	return i
}

type Stack []interface{}

func (s *Stack) Push(elem interface{}) {
	*s = append(*s, elem)
}

func (s *Stack) Pop() (interface{}, error) {
	len := len(*s)
	if len == 0 {
		return 0, errors.New("stack is empty")
	}
	res := (*s)[len-1]
	*s = (*s)[:len-1]
	return res, nil
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}
