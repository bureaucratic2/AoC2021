package common

import (
	"os"
	"path"
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
