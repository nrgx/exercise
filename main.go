package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"unicode"
)

type Entry struct {
	Word []byte
	Freq int
}

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("failed to open file: %v\n", err)
	}
	defer file.Close()
	var entries []Entry
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("corrupted file: %v\n", err)
		}
		if err == io.EOF {
			break
		}
		line = process(line)
		words := bytes.Fields(line)
		for i := range words {
			incrEntry(&entries, words[i])
		}
	}
	// sort in ascending order
	sort.Slice(entries, func(i, j int) bool { return entries[i].Freq > entries[j].Freq })
	for i := 0; i < 20; i++ {
		fmt.Printf("%7d %-5s\n", entries[i].Freq, entries[i].Word)
	}
}

// cleans line from non-alphabetic chars and lowers them
func process(s []byte) []byte {
	var result bytes.Buffer
	for _, char := range s {
		if unicode.IsLetter(rune(char)) {
			result.WriteRune(unicode.ToLower(rune(char)))
		} else {
			result.WriteRune(' ')
		}
	}
	return result.Bytes()
}

func incrEntry(entries *[]Entry, word []byte) {
	for i := range *entries {
		if bytes.Equal((*entries)[i].Word, word) {
			(*entries)[i].Freq++
			return
		}
	}
	*entries = append(*entries, Entry{Word: []byte(word), Freq: 1})
}
