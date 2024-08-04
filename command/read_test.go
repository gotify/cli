package command

import (
	"strings"
	"testing"
)

// Polyfill for slices.Equal for Go 1.20
func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func readChanAll[T any](c chan T) []T {
	var res []T
	for s := range c {
		res = append(res, s)
	}
	return res
}

func TestReadMessage(t *testing.T) {
	var split rune = '\x00'

	// Test case 1: message set via arguments
	output := make(chan string)
	go readMessage([]string{"Hello", "World"}, nil, output, nil)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello World"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello World"}, res)
	}

	// Test case 2: message set via arguments should not split on 'split' character
	output = make(chan string)
	go readMessage([]string{"Hello\x00World"}, nil, output, &split)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello\x00World"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello\x00World"}, res)
	}

	// Test case 3: message set via stdin
	output = make(chan string)
	go readMessage([]string{}, strings.NewReader("Hello\x00World"), output, &split)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello", "World"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello", "World"}, res)
	}
}
