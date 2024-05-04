//go:build linux

package main

import (
	"os"
	"syscall"
)

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

	count := 0
	for i := 0; i < len(letters); i++ {
		letters[count] = letters[i]
		count += int(((^mask[i]) & 0b10) >> 1)
	}

	os.Stdout.Write(letters[:count])
}

func m[T any](val T, err error) T {
	if err != nil {
		println(err.Error())
		os.Exit(3)
	}
	return val
}
