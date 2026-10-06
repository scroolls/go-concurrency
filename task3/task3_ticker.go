package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Отправляет номер тика и время с момента старта в канал tickCh.
// Завершается когда закрывается канал done.
func runTicker(stepSec uint, done <-chan struct{}) <-chan [2]int64 {
	tickCh := make(chan [2]int64)
	go func() {
		defer close(tickCh)
		start := time.Now()
		tick := 1
		for {
			deadline := time.Now().Add(time.Duration(stepSec) * time.Second)
			for time.Now().Before(deadline) {
				select {
				case <-done:
					return
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}

			select {
			case <-done:
				return
			case tickCh <- [2]int64{int64(tick), int64(time.Since(start).Seconds())}:
				tick++
			}
		}
	}()
	return tickCh
}

func main() {
	k := flag.Uint("k", 1, "Ticker step in seconds")
	flag.Parse()

	if *k == 0 {
		fmt.Println("K must be a positive integer")
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigCh)

	done := make(chan struct{})
	tickCh := runTicker(*k, done)

	fmt.Printf("Ticker started with step %d second(s). Press Ctrl+C to stop.\n", *k)

	for {
		select {
		case tick, ok := <-tickCh:
			if !ok {
				return
			}
			fmt.Printf("Tick %d since %d sec\n", tick[0], tick[1])
		case sig := <-sigCh:
			fmt.Printf("\nReceived signal: %s\n", sig)
			close(done)
			for range tickCh {
			}
			fmt.Println("Termination")
			return
		}
	}
}
