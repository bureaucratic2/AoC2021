package common

import (
	"os"
	"path"
	"unicode"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

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
