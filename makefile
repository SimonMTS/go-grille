all: run clean

build:
	go build -o grille main.go

clean:
	rm -f ./grille

run: build
	./grille ./mask.txt ./letters.txt | md5sum
	# 05ea1fcfd9473c0ba81a20ee03a68814
	# is correct

bench:
	go test -bench=. -benchtime 10x -benchmem -memprofile memprofile.out -cpuprofile profile.out

bench-cpu:
	go tool pprof -http 127.0.0.1:8081 profile.out

bench-mem:
	go tool pprof -http 127.0.0.1:8082 memprofile.out

perf: build
	perf stat -nr 5 sh -c './grille ./mask.txt ./letters.txt > /dev/null'
