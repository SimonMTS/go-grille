//go:build linux

package main

import (
	"os"
	"syscall"
)

const routines = 16

func main() {
	if len(os.Args) < 3 {
		println("2 file arguments required: [mask] [letters]")
		os.Exit(1)
	}

	// Map mask file into memory
	maskFile := m(os.Open(os.Args[1]))
	maskStat := m(maskFile.Stat())
	mask := m(syscall.Mmap(
		int(maskFile.Fd()),
		0,
		int(maskStat.Size()),
		syscall.PROT_READ|
			syscall.MAP_LOCKED|
			syscall.MAP_NORESERVE|
			syscall.MAP_POPULATE,
		syscall.MAP_PRIVATE,
	))

	// Map letter file into memory
	lettersFile := m(os.Open(os.Args[2]))
	lettersStat := m(lettersFile.Stat())
	letters := m(syscall.Mmap(
		int(lettersFile.Fd()),
		0,
		int(lettersStat.Size()),
		syscall.PROT_WRITE|
			syscall.PROT_READ|
			syscall.MAP_LOCKED|
			syscall.MAP_NORESERVE|
			syscall.MAP_POPULATE,
		syscall.MAP_PRIVATE,
	))

	// Split input into processable chunks
	var (
		inputSize = len(mask)
		chunkSize = inputSize / routines
		metaData  = make([]meta, 0, routines)
	)

	if inputSize-(chunkSize*routines) != 0 {
		println("could not split input into evenly sized sections")
		os.Exit(2)
	}

	for i := 0; i < routines; i++ {
		metaData = append(metaData, meta{
			Start: i * chunkSize,
			End:   (i + 1) * chunkSize,
		})
	}

	// Calculate outputs
	out := make([]chan meta, routines)
	for i := range out {
		out[i] = make(chan meta)
		go processSection(mask, letters, &metaData[i], out[i])
	}

	for i := range out {
		m := <-out[i]
		os.Stdout.Write(letters[m.Start:m.End])
	}
}

type meta struct {
	Start int
	End   int
}

//           Use this bit
//                |
//  " " = 0b00100000
//  "#" = 0b00100011
// "\n" = 0b00001010

func processSection(mask, letters []byte, data *meta, out chan meta) {
	count := data.Start
	for i := data.Start; i < data.End; i++ {
		letters[count] = letters[i]
		count += int(((^mask[i]) & 0b10) >> 1)
	}

	data.End = count
	out <- *data
}

func m[T any](val T, err error) T {
	if err != nil {
		println(err.Error())
		os.Exit(3)
	}
	return val
}
