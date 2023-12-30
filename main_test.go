package main

import (
	"testing"
)

func TestFunc1(t *testing.T) {

	// Test 1
	if func1(1, 3) != 4 {
		t.Error("func1(1, 3) != 4")
	}

	// Test 2
	if func1(2, 3) == 4 {
		t.Error("func1(1, 3) == 4")
	}

}