go:
	go build main.go

c:
	clang -Ofast main.c

build: go c

valid:
	./main ../mask.txt ../letters.txt | md5sum
	./a.out ../mask.txt ../letters.txt | md5sum
	../bin/optimized ../mask.txt ../letters.txt | md5sum

bench:
	perf stat -nr 50 ../bin/optimized ../mask.txt ../letters.txt > /dev/null
	perf stat -nr 50 ../bin/optimized-v2 ../mask.txt ../letters.txt > /dev/null
	perf stat -nr 50 ./main ../mask.txt ../letters.txt > /dev/null
	perf stat -nr 50 ./a.out ../mask.txt ../letters.txt > /dev/null
	perf stat -nr 50 cat ../mask.txt ../letters.txt > /dev/null

perf:
	perf stat -r 50 ./a.out ../mask.txt ../letters.txt > /dev/null
