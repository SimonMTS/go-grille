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
