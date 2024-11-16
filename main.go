package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "time"
    "strings"
    "strconv"
    "sync"
)

const MAX_THREADS = 100
const MAX_IP = 256
const BUF_SIZE = 1000

var wg sync.WaitGroup
var channel [MAX_THREADS](chan string)
var partNonUnique [MAX_THREADS]int
var table [MAX_IP][MAX_IP][MAX_IP][MAX_IP / 8]uint8
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

func threadLoop(id int, c chan string) {
    defer wg.Done()

    for address := range c {
        parts := strings.Split(address, ".")

        index := [4]int{}
        for i := 0; i < 4; i += 1 {
            value, _ := strconv.Atoi(parts[i])
            index[i] = value
        }

        p := &table[index[0]][index[1]][index[2]][index[3] / 8]
        shift := uint8(index[3] % 8)

        if hasBit(*p, shift) {
            partNonUnique[id] += 1
        } else {
            *p = setBit(*p, shift)
        }
    }
}

func main() {
    defer trackTime(time.Now(), "main")
    filename := "/Users/clarence/Desktop/task/ip_addresses_full"

    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Failed to open file: %s", err)
    }
    defer file.Close()

    wg.Add(MAX_THREADS)
    for i := 0; i < MAX_THREADS; i += 1 {
        channel[i] = make(chan string, BUF_SIZE)
        go threadLoop(i, channel[i])
    }

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lineCounter += 1
        firstPart, _ := strings.CutPrefix(scanner.Text(), ".")
        firstNum, _ := strconv.Atoi(firstPart)
        id := firstNum / (MAX_IP / MAX_THREADS)
        
        channel[id] <- scanner.Text()
    }

    for i := 0; i < MAX_THREADS; i += 1 {
        close(channel[i])
    }
    wg.Wait()

    for i := 0; i < MAX_THREADS; i += 1 {
        nonUnique += partNonUnique[i];
    }

    fmt.Printf("All lines: %d, Unique lines: %d\n", lineCounter, lineCounter - nonUnique)
}
