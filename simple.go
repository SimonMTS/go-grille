package main

import (
	"bufio"
	"io"
	"io/ioutil"
)

func simple(maskFile, letterFile string, out io.Writer) {
	grill, _ := ioutil.ReadFile(maskFile)
	cipher, _ := ioutil.ReadFile(letterFile)
	w := bufio.NewWriter(out)
	defer w.Flush()

	for i, b := range grill {
		if b == ' ' {
			w.WriteByte(cipher[i])
		}
	}
}
