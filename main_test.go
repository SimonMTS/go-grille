package main

import "testing"

func BenchmarkOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		optimized("./mask.txt", "letters.txt", write{})
	}
}

func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simple("./mask.txt", "letters.txt", write{})
	}
}

type write struct{}

func (write) Write(p []byte) (n int, err error) { return }
