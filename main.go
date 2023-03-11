package main

import (
	"os"
	"runtime"

	"github.com/edsrzf/mmap-go"
)

func main() {
	optimized(os.Args[1], os.Args[2], os.Stdout)
}

func optimized(maskFile, letterFile string, out *os.File) {
	// Read input
	mf, _ := os.Open(maskFile)
	mask, _ := mmap.Map(mf, 0, 0)
	mask.Lock()

	lf, _ := os.Open(letterFile)
	letters, _ := mmap.Map(lf, mmap.COPY, 0)
	letters.Lock()

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
	done := make(chan struct{})
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

//  " " = 0b00100000
//  "#" = 0b00100011
// "\n" = 0b00001010

func processSection(mask, letters []byte, data *meta, done chan struct{}) {
	count := data.Start
	for i := data.Start; i < data.End; i++ {
		letters[count] = letters[i]
		count += int(^mask[i] & 0b00000010 >> 1)
	}

	data.NewEnd = count
	done <- struct{}{}
}
