all: build

build:
	go build -o bin/optimized    -gcflags=-B ./cmd/optimized
	go build -o bin/optimized-v2 -gcflags=-B ./cmd/optimized-v2
	go build -o bin/instrumented -gcflags=-B ./cmd/instrumented
	go build -o bin/ref ./cmd/ref

# build:
# 	go build -o grille -gcflags=-B ./main.go

# clean:
# 	rm -f ./grille

# run: build
# 	./grille ./mask.txt ./letters.txt | md5sum
# 	# 05ea1fcfd9473c0ba81a20ee03a68814
# 	# is correct

# compare:
# 	go test -bench=. -benchtime 20x -benchmem

# bench:
# 	go test -bench="Optimized" -benchtime 10x -benchmem -memprofile memprofile.out -cpuprofile profile.out

# bench-cpu:
# 	go tool pprof -http 127.0.0.1:8081 profile.out

# bench-mem:
# 	go tool pprof -http 127.0.0.1:8082 memprofile.out

# bench-loop:
# 	go test -bench="Loop" -benchmem -memprofile memprofile.out -cpuprofile profile.out

# bench-branchless:
# 	go test -bench="Branchless" -benchmem -memprofile memprofile.out -cpuprofile profile.out

# perf: build
# 	perf stat -nr 100 ./grille ./mask.txt ./letters.txt > /dev/null

