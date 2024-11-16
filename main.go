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

const MAX_THREADS = 400
const MAX_IP = 256
var wg sync.WaitGroup
var mutex sync.Mutex
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

func processLine(address string) {
    defer wg.Done()

    parts := strings.Split(address, ".")

    index := [4]int{}
    for i := 0; i < 4; i += 1 {
        value, _ := strconv.Atoi(parts[i])
        index[i] = value
    }

    mutex.Lock()
    p := &table[index[0]][index[1]][index[2]][index[3] / 8]
    shift := uint8(index[3] % 8)

    if hasBit(*p, shift) {
        nonUnique += 1
    } else {
        *p = setBit(*p, shift)
    }
    mutex.Unlock()
}

func main() {
    defer trackTime(time.Now(), "main")
    filename := "/Users/clarence/Desktop/task/ip_addresses_full"

    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Failed to open file: %s", err)
    }
    defer file.Close()

    buffer := [MAX_THREADS]string{}
    bufCounter := 0

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lineCounter += 1
        buffer[bufCounter] = scanner.Text()
        bufCounter += 1
        if bufCounter == MAX_THREADS {
            wg.Add(MAX_THREADS)
            for i := 0; i < MAX_THREADS; i += 1 {
                go processLine(buffer[i])
            }
            wg.Wait()
            bufCounter = 0
        }
    }

    if bufCounter != 0 {
        wg.Add(bufCounter)
        for i := 0; i < bufCounter; i += 1 {
            go processLine(buffer[i])
        }
        wg.Wait()
    }

    fmt.Printf("All lines: %d, Unique lines: %d\n", lineCounter, lineCounter - nonUnique)
}
