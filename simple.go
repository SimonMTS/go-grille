package main

import (
	"io"
	"os"
)

func simple(maskFile, letterFile string, out io.Writer) {
	mask, _ := os.ReadFile(maskFile)
	letters, _ := os.ReadFile(letterFile)

	for i, b := range mask {
		if b == ' ' {
			out.Write([]byte{letters[i]})
		}
	}
}
