package main

import (
	"io"
	"os"
	"runtime"
)

func ProcessItem(out *[]byte, mask, letters []byte, done chan bool) {
	for i, b := range mask {
		if b == ' ' {
			*out = append(*out, letters[i])
		}
	}
	done <- true
}

type Item struct {
	Mask    []byte
	Letters []byte
}

func main() {
	optimized(os.Args[1], os.Args[2], os.Stdout)
}

func optimized(maskFile, letterFile string, out io.Writer) {
	// Read input
	mask, merr := os.ReadFile(maskFile)
	letters, lerr := os.ReadFile(letterFile)
	if merr != nil || lerr != nil || len(mask) != len(letters) {
		panic("bad input")
	}

	// Split input into processable chunks
	var (
		inputSize  = len(mask)
		chunkCount = runtime.NumCPU() * 2
		chunkSize  = inputSize / chunkCount
		leftOver   = inputSize - (chunkSize * chunkCount)
		items      = make([]Item, chunkCount)
	)
	for i := 0; i < chunkCount; i++ {
		items = append(items, Item{
			Mask:    mask[i*chunkSize : (i+1)*chunkSize],
			Letters: letters[i*chunkSize : (i+1)*chunkSize],
		})
	}
	if leftOver > 0 {
		items = append(items, Item{
			Mask:    mask[inputSize-leftOver:],
			Letters: letters[inputSize-leftOver:],
		})
	}

	// Calculate outputs
	results := make([][]byte, len(items))
	done := make(chan bool, len(items))
	for i, item := range items {
		go ProcessItem(&results[i], item.Mask, item.Letters, done)
	}

	// Wait for all processing to be done
	for range results {
		<-done
	}

	// Print outputs
	for _, result := range results {
		out.Write(result)
	}
}
