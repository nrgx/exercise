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
		input    []byte
		expected []byte
	}{
		{[]byte("Hello, World!"), []byte("hello  world ")},
		{[]byte("123ABCxyz"), []byte("   abcxyz")},
		{[]byte("!@#$%^&*()"), []byte("          ")},
		{[]byte("Lorem Ipsum"), []byte("lorem ipsum")},
	}

	for _, tc := range testCases {
		actual := process(tc.input)
		if !bytes.Equal(actual, tc.expected) {
			t.Errorf("process(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestIncrEntry(t *testing.T) {
	entries := []Entry{
		{Word: []byte("hello"), Freq: 2},
		{Word: []byte("world"), Freq: 1},
	}

	incrEntry(&entries, []byte("hello"))
	expected := []Entry{
		{Word: []byte("hello"), Freq: 3},
		{Word: []byte("world"), Freq: 1},
	}

	if !entriesEqual(entries, expected) {
		t.Errorf("Unexpected result after incrementing entry")
	}

	incrEntry(&entries, []byte("new"))
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
	// Just copied alphabet letters 8 times because removed threshold variable in main.
	// Outputs only 20 letters in order.
	fileContent := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\na\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\na\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\na\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\na\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\na\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\na\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\na\nb\nc\nd\ne\nf\ng\nh\ni\nj\nk\nl\nm\nn\no\np\nq\nr\ns\nt\nu\nv\nw\nx\ny\nz\n"
	// Create a tmp file for testing purposes
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(file.Name())
	file.WriteString(fileContent)
	file.Close()

	os.Args = []string{"", file.Name()}

	// Actual result returns formated code, should be the same in expected.
	expectedOutput := fmt.Sprintf(
		"%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
		fmt.Sprintf("%7d %-5s\n", 8, "a"),
		fmt.Sprintf("%7d %-5s\n", 8, "b"),
		fmt.Sprintf("%7d %-5s\n", 8, "c"),
		fmt.Sprintf("%7d %-5s\n", 8, "d"),
		fmt.Sprintf("%7d %-5s\n", 8, "e"),
		fmt.Sprintf("%7d %-5s\n", 8, "f"),
		fmt.Sprintf("%7d %-5s\n", 8, "g"),
		fmt.Sprintf("%7d %-5s\n", 8, "h"),
		fmt.Sprintf("%7d %-5s\n", 8, "i"),
		fmt.Sprintf("%7d %-5s\n", 8, "j"),
		fmt.Sprintf("%7d %-5s\n", 8, "k"),
		fmt.Sprintf("%7d %-5s\n", 8, "l"),
		fmt.Sprintf("%7d %-5s\n", 8, "m"),
		fmt.Sprintf("%7d %-5s\n", 8, "n"),
		fmt.Sprintf("%7d %-5s\n", 8, "o"),
		fmt.Sprintf("%7d %-5s\n", 8, "p"),
		fmt.Sprintf("%7d %-5s\n", 8, "q"),
		fmt.Sprintf("%7d %-5s\n", 8, "r"),
		fmt.Sprintf("%7d %-5s\n", 8, "s"),
		fmt.Sprintf("%7d %-5s\n", 8, "t"),
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
