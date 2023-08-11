//go:build linux

package main

import (
	"log"
	"os"
	"runtime/pprof"
	"syscall"
)

// Should be a good bit more than the number of core's
const routines = 100

func main() {
	// Profile
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	if len(os.Args) < 3 {
		panic("2 file arguments required: [mask] [letters]")
	}

	// Map mask file into memory
	maskFile := m(os.Open(os.Args[1]))
	maskStat := m(maskFile.Stat())
	mask := m(syscall.Mmap(
		int(maskFile.Fd()),
		0,
		int(maskStat.Size()),
		syscall.PROT_READ|syscall.MAP_LOCKED|syscall.MAP_NORESERVE,
		syscall.MAP_PRIVATE,
	))

	// Map letter file into memory
	lettersFile := m(os.Open(os.Args[2]))
	lettersStat := m(lettersFile.Stat())
	letters := m(syscall.Mmap(
		int(lettersFile.Fd()),
		0,
		int(lettersStat.Size()),
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.MAP_LOCKED|syscall.MAP_NORESERVE,
		syscall.MAP_PRIVATE,
	))

	// Split input into processable chunks
	var (
		inputSize = len(mask)
		chunkSize = inputSize / routines
		metaData  = make([]meta, 0, routines)
	)

	if inputSize-(chunkSize*routines) != 0 {
		panic("could not split input into evenly sized sections")
	}

	for i := 0; i < routines; i++ {
		metaData = append(metaData, meta{
			Start: i * chunkSize,
			End:   (i + 1) * chunkSize,
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
		os.Stdout.Write(letters[d.Start:d.End])
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

func processSection(mask, letters []byte, data *meta, done chan struct{}) {
	count := data.Start
	for i := data.Start; i < data.End; i += 10 {
		letters[count] = letters[i]
		letters[count] = letters[i+1]
		letters[count] = letters[i+2]
		letters[count] = letters[i+3]
		letters[count] = letters[i+4]
		letters[count] = letters[i+5]
		letters[count] = letters[i+6]
		letters[count] = letters[i+7]
		letters[count] = letters[i+8]
		letters[count] = letters[i+9]

		count += int(((^mask[i]) & 0b10) >> 1)
		count += int(((^mask[i+1]) & 0b10) >> 1)
		count += int(((^mask[i+2]) & 0b10) >> 1)
		count += int(((^mask[i+3]) & 0b10) >> 1)
		count += int(((^mask[i+4]) & 0b10) >> 1)
		count += int(((^mask[i+5]) & 0b10) >> 1)
		count += int(((^mask[i+6]) & 0b10) >> 1)
		count += int(((^mask[i+7]) & 0b10) >> 1)
		count += int(((^mask[i+8]) & 0b10) >> 1)
		count += int(((^mask[i+9]) & 0b10) >> 1)
	}

	data.End = count
	done <- struct{}{}
}

func m[T any](val T, err error) T {
	if err != nil {
		panic(err.Error())
	}
	return val
}
