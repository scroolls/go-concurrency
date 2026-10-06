package main

import (
	"reflect"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	var got []int
	for value := range squarer(generator(-2, 2)) {
		got = append(got, value)
	}
	if !reflect.DeepEqual(got, []int{4, 1, 0, 1, 4}) {
		t.Fatalf("squares = %v", got)
	}
}

func TestGeneratorAtMaxInt(t *testing.T) {
	maximum := int(^uint(0) >> 1)
	ch := generator(maximum, maximum)
	if got := <-ch; got != maximum {
		t.Fatalf("got %d", got)
	}
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("generator overflowed instead of closing")
		}
	case <-time.After(time.Second):
		t.Fatal("generator did not close")
	}
}
