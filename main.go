package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "time"
    "strings"
    "strconv"
)

var table [256][256][256][32]uint8
var index [4]int
var lineCounter, nonUnique = 0, 0

func trackTime(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

func hasBit(n uint8, pos uint8) bool {
    val := n & (1 << pos)
    return (val > 0)
}

func setBit(n uint8, pos uint8) uint8 {
    n |= (1 << pos)
    return n
}

func processLine(address string) {
    parts := strings.Split(address, ".")
    for i := 0; i < 4; i += 1 {
        value, _ := strconv.Atoi(parts[i])
        index[i] = value
    }

    p := &table[index[0]][index[1]][index[2]][index[3] / 8]
    shift := uint8(index[3] % 8)

    if hasBit(*p, shift) { // ip address already encountered
        nonUnique += 1
    } else {
        *p = setBit(*p, shift)
    }
}

func main() {
    defer trackTime(time.Now(), "main")

    file, err := os.Open("/Users/clarence/Desktop/task/ip_addresses_full")
    if err != nil {
        log.Fatalf("Failed to open file: %s", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lineCounter += 1
        processLine(scanner.Text())
    }

    fmt.Printf("All lines: %d, Unique lines: %d\n", lineCounter, lineCounter - nonUnique)
}
