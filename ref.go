package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// mask.txt
// "# ## #\n"
// "##  ##\n"
// "# ### \n"

// letters.txt
// "abcdef\n"
// "ghijkl\n"
// "mnopqr\n"

// $ go run ref.go ./mask.txt ./letters.txt
// beijnr

// this implementation takes a bit over 30 seconds on the 10.000x10.000 inputs

// big inputs:
// https://raw.githubusercontent.com/SimonMTS/go-grille/master/letters.txt
// https://raw.githubusercontent.com/SimonMTS/go-grille/master/mask.txt

// verify result by getting the md5 of the output (no newlines):
// 05ea1fcfd9473c0ba81a20ee03a68814

// measure performance with:
// perf stat -nr 1 ./program ./mask.txt ./letters.txt > /dev/null

func main() {
	var mask []byte
	var letters []byte

	mask, _ = ioutil.ReadFile(os.Args[1])
	letters, _ = ioutil.ReadFile(os.Args[2])

	for i := 0; i < len(mask); i++ {
		if mask[i] == ' ' {
			fmt.Print(string(letters[i]))
		}
	}
}
