package common

import (
	"os"
	"path"
)

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
