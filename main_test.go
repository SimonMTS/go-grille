package main

import (
	"os"
	"testing"
)

func BenchmarkOptimized(b *testing.B) {
	f, _ := os.Open("/dev/null")
	for i := 0; i < b.N; i++ {
		optimized("./mask.txt", "letters.txt", f)
	}
}

func BenchmarkSimple(b *testing.B) {
	f, _ := os.Open("/dev/null")
	for i := 0; i < b.N; i++ {
		simple("./mask.txt", "letters.txt", f)
	}
}

func BenchmarkLoop(b *testing.B) {
	var (
		lttrsRaw = []byte("abcdefghijklmnopqrstuvwxyz")
		mask     = []byte("## #  #  ##   # #####   ##")
		metaData = meta{Start: 0, End: len(mask)}
		done     = make(chan struct{}, 1)
	)

	for i := 0; i < b.N; i++ {
		letters := lttrsRaw
		processSection(mask, letters, &metaData, done)
		<-done
	}
}

func TestBranchless(t *testing.T) {
	in1 := byte(' ')
	in2 := byte('#')
	in3 := byte('\n')

	f := func(i byte) byte {
		return 0b1 ^ (i<<5)>>6
	}

	out1 := f(in1)
	out2 := f(in2)
	out3 := f(in3)

	if out1 != 1 {
		t.Errorf("%b != 1", out1)
	}

	if out2 != 0 {
		t.Errorf("%b != 0", out2)
	}

	if out3 != 0 {
		t.Errorf("%b != 0", out3)
	}
}

func BenchmarkBranchless1(b *testing.B) {
	a := []byte{' ', '#', '\n'}
	c := []byte{1, 0, 0}
	for i := 0; i < b.N; i++ {
		j := i % 3
		s := a[j]
		r := (^s & 0b00000010 >> 1)
		if r != c[j] {
			b.Fail()
		}
	}
}

func BenchmarkBranchless2(b *testing.B) {
	a := []byte{' ', '#', '\n'}
	c := []byte{1, 0, 0}
	for i := 0; i < b.N; i++ {
		j := i % 3
		s := a[j]
		r := (0b1 ^ (s<<5)>>6)
		if r != c[j] {
			b.Fail()
		}
	}
}
