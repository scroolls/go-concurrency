package main

import (
	"flag"
	"fmt"
)

func generator(k, n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := k; i <= n; i++ {
			ch <- i
			if i == n {
				break // Avoid wrapping past the maximum int.
			}
		}
	}()
	return ch
}

func squarer(in <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for v := range in {
			ch <- v * v
		}
	}()
	return ch
}

func main() {
	k := flag.Int("k", 1, "Start of range")
	n := flag.Int("n", 10, "End of range (inclusive)")
	flag.Parse()

	if *k > *n {
		fmt.Println("K must be less than or equal to N")
		return
	}

	genCh := generator(*k, *n)
	sqCh := squarer(genCh)

	for v := range sqCh {
		fmt.Println(v)
	}
}
