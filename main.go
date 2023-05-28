package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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
	// if there is need in control how many results to output
	threshold, err := strconv.Atoi(os.Args[2])
	if err != nil {
		threshold = 20
	}
	defer file.Close()
	var entries []Entry
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("corrupted file: %v\n", err)
		}
		if err == io.EOF {
			break
		}
		line = process(line)
		words := strings.Fields(line)
		for i := range words {
			incrEntry(&entries, words[i])
		}
	}
	// sort in ascending order
	sort.Slice(entries, func(i, j int) bool { return entries[i].Freq > entries[j].Freq })
	for i := 0; i < threshold; i++ {
		fmt.Printf("%7d %-5s\n", entries[i].Freq, entries[i].Word)
	}
}

// cleans line from non-alphabetic chars and lowers them
func process(s string) string {
	var result strings.Builder
	for _, char := range s {
		if unicode.IsLetter(char) {
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(' ')
		}
	}
	return result.String()
}

func incrEntry(entries *[]Entry, word string) {
	for i := range *entries {
		if string((*entries)[i].Word) == word {
			(*entries)[i].Freq++
			return
		}
	}
	*entries = append(*entries, Entry{Word: []byte(word), Freq: 1})
}
