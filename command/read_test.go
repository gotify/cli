package command

import (
	"bufio"
	"bytes"
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
	if bytes.IndexByte([]byte("Hello\x00World"), '\x00') != len("Hello") {
		t.Errorf("Expected %v, but got %v", len("Hello"), bytes.IndexByte([]byte("Hello\x00World"), '\x00'))
	}
	rdr := bufio.NewReader(strings.NewReader("Hello\x00World"))
	if s, _ := rdr.ReadString('\x00'); s != "Hello\x00" {
		t.Errorf("Expected %x, but got %x", "Hello\x00", s)
	}
	// Test case 1: message set via arguments
	output := make(chan string)
	go readMessage([]string{"Hello", "World"}, nil, output, false)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello World"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello World"}, res)
	}

	// Test case 2: message set via arguments should not split on 'split' character
	output = make(chan string)
	go readMessage([]string{"Hello\x00World"}, nil, output, true)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello\x00World"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello\x00World"}, res)
	}

	// Test case 3: message set via stdin
	output = make(chan string)
	go readMessage([]string{}, strings.NewReader("Hello\x00World"), output, true)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello", "World"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello", "World"}, res)
	}

	// Test case 4: multiple null bytes should be split as one
	output = make(chan string)
	go readMessage([]string{}, strings.NewReader("Hello\x00\x00World"), output, true)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello", "World"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello", "World"}, res)
	}

	// Test case 5: multiple null bytes at the end should be split as one
	output = make(chan string)
	go readMessage([]string{}, strings.NewReader("Hello\x00\x00"), output, true)

	if res := readChanAll(output); !(slicesEqual(res, []string{"Hello"})) {
		t.Errorf("Expected %v, but got %v", []string{"Hello"}, res)
	}

	// Test case 6: multiple null bytes at the start should be split as one
	output = make(chan string)
	go readMessage([]string{}, strings.NewReader("\x00\x00World"), output, true)

	if res := readChanAll(output); !(slicesEqual(res, []string{"World"})) {
		t.Errorf("Expected %v, but got %v", []string{"World"}, res)
	}

}
