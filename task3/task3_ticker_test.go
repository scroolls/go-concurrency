package main

import (
	"testing"
	"time"
)

func TestTickerStopsWithoutReceiver(t *testing.T) {
	done := make(chan struct{})
	ch := runTicker(0, done)
	// Let the producer reach an unbuffered send without consuming it.
	time.Sleep(20 * time.Millisecond)
	close(done)
	time.Sleep(20 * time.Millisecond)
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("ticker remained blocked on a send after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("ticker did not stop")
	}
}

func TestTickerFirstTick(t *testing.T) {
	done := make(chan struct{})
	ch := runTicker(1, done)
	defer func() {
		close(done)
		for range ch {
		}
	}()
	select {
	case tick := <-ch:
		if tick[0] != 1 || tick[1] < 1 {
			t.Fatalf("unexpected tick %v", tick)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no tick received")
	}
}
