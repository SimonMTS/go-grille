package main

import (
	"os"
	"runtime"
)

func main() {
	optimized(os.Args[1], os.Args[2], os.Stdout)
}

func optimized(maskFile, letterFile string, out *os.File) {
	// Read input
	mask, merr := os.ReadFile(maskFile)
	letters, lerr := os.ReadFile(letterFile)
	if merr != nil || lerr != nil || len(mask) != len(letters) {
		panic("bad input")
	}

	// Split input into processable chunks
	var (
		inputSize  = len(mask)
		chunkCount = runtime.NumCPU() * 8
		chunkSize  = inputSize / chunkCount
		leftOver   = inputSize - (chunkSize * chunkCount)
		metaData   = make([]meta, 0, chunkCount)
	)
	for i := 0; i < chunkCount; i++ {
		metaData = append(metaData, meta{
			Start: i * chunkSize,
			End:   (i + 1) * chunkSize,
		})
	}
	if leftOver > 0 {
		metaData = append(metaData, meta{
			Start: inputSize - leftOver,
			End:   inputSize,
		})
	}

	// Calculate outputs
	done := make(chan bool)
	for i := range metaData {
		go processSection(mask, letters, &metaData[i], done)
	}

	// Wait for all processing to be done
	for range metaData {
		<-done
	}

	for _, d := range metaData {
		out.Write(letters[d.Start:d.NewEnd])
	}
}

type meta struct {
	Start  int
	End    int
	NewEnd int
}

func processSection(mask, letters []byte, data *meta, done chan bool) {
	count := data.Start
	for i := data.Start; i < data.End; i++ {
		if mask[i] == ' ' {
			letters[count] = letters[i]
			count++
		}
	}

	data.NewEnd = count
	done <- true
}
