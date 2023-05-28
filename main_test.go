package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

func TestProcess(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "hello  world "},
		{"123ABCxyz", "   abcxyz"},
		{"!@#$%^&*()", "          "},
		{"Lorem Ipsum", "lorem ipsum"},
	}

	for _, tc := range testCases {
		actual := process(tc.input)
		if actual != tc.expected {
			t.Errorf("process(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestIncrEntry(t *testing.T) {
	entries := []Entry{
		{Word: []byte("hello"), Freq: 2},
		{Word: []byte("world"), Freq: 1},
	}

	incrEntry(&entries, "hello")
	expected := []Entry{
		{Word: []byte("hello"), Freq: 3},
		{Word: []byte("world"), Freq: 1},
	}

	if !entriesEqual(entries, expected) {
		t.Errorf("Unexpected result after incrementing entry")
	}

	incrEntry(&entries, "new")
	expected = append(expected, Entry{Word: []byte("new"), Freq: 1})

	if !entriesEqual(entries, expected) {
		t.Errorf("Unexpected result after adding new entry")
	}
}

func entriesEqual(entries1, entries2 []Entry) bool {
	if len(entries1) != len(entries2) {
		return false
	}
	for i := range entries1 {
		if !bytes.Equal(entries1[i].Word, entries2[i].Word) || entries1[i].Freq != entries2[i].Freq {
			return false
		}
	}
	return true
}

func TestMain(t *testing.T) {
	// Create a tmp file for testing purposes
	fileContent := "Hello, World!\nHello, Go World!\n"
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(file.Name())
	file.WriteString(fileContent)
	file.Close()

	// Arguments to run program. Bc of small input (3 words) it'll panic so I put threshold value that limits execution.
	os.Args = []string{"", file.Name(), "3"}

	// Actual result returns formated code, should be the same in expected.
	expectedOutput := fmt.Sprintf(
		"%s%s%s",
		fmt.Sprintf("%7d %-5s\n", 2, "hello"),
		fmt.Sprintf("%7d %-5s\n", 2, "world"),
		fmt.Sprintf("%7d %-5s\n", 1, "go"),
	)
	actualOutput := captureOutput(func() { main() })

	if actualOutput != expectedOutput {
		t.Errorf("Unexpected output. Expected:\n%s\nGot:\n%s", expectedOutput, actualOutput)
	}
}

// captureOutput - a helper function that captures log output on STDOUT
func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
