#include <stdio.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <unistd.h>
#include <sys/mman.h>

int main(int argc, char **argv) {
    if (argc < 3) {
        printf("2 file arguments required: [mask] [letters]\n");
        return 1;
    }

    int maskFile = open(argv[1], O_RDONLY);
    struct stat maskStat;
    fstat(maskFile, &maskStat);
    char* mask = mmap(NULL,
            maskStat.st_size,
            PROT_READ,
            // MAP_LOCKED doesn't work?
            MAP_PRIVATE|MAP_NORESERVE|MAP_POPULATE,
            maskFile, 0);
    if (mask == MAP_FAILED) {
        perror("mask");
        return 2;
    }

    int letterFile = open(argv[2], O_RDONLY);
    struct stat letterStat;
    fstat(letterFile, &letterStat);
    char* letters = mmap(0,
            letterStat.st_size,
            PROT_WRITE|PROT_READ,
            // MAP_LOCKED doesn't work?
            MAP_PRIVATE|MAP_NORESERVE|MAP_POPULATE,
            letterFile, 0);
    if (letters == MAP_FAILED) {
        perror("letters");
        return 3;
    }

    int count = 0;
    for (int i = 0; i < maskStat.st_size; i++) {
        letters[count] = letters[i];
        count += ((~mask[i]) & 0b10) >> 1;
    }

    write(fileno(stdout), letters, count);
    return 0;
}
