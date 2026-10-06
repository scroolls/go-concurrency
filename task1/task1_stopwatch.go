package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type GoroutineResult struct {
	ID      int
	SleepMs int
}

func main() {
	n := flag.Int("n", 5, "Number of goroutines")
	m := flag.Int("m", 1000, "Max sleep time in milliseconds")
	flag.Parse()

	if *n <= 0 || *m <= 0 {
		fmt.Println("N and M must be positive integers")
		return
	}

	results := make([]GoroutineResult, *n)
	var wg sync.WaitGroup

	for i := 0; i < *n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sleepMs := rand.Intn(*m) + 1
			time.Sleep(time.Duration(sleepMs) * time.Millisecond)
			results[id] = GoroutineResult{ID: id, SleepMs: sleepMs}
		}(i)
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].SleepMs > results[j].SleepMs
	})

	fmt.Println("Goroutine results (sorted by sleep time descending):")
	for _, r := range results {
		fmt.Printf("Goroutine #%d: %d ms\n", r.ID, r.SleepMs)
	}
}
