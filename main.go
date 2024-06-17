package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type BookStats struct {
	words int
	chars map[string]int
}

const (
	bufferSize = 8096
)

func (bs *BookStats) countWords(readString string) {

	words := strings.Fields(readString)
	bs.words += len(words)
}

func (bs *BookStats) countChars(readString string) {

	lowered := strings.ToLower(readString)

	for _, char := range lowered {
		charString := string(char)
		if char > 96 && char < 123 {
			bs.chars[charString]++
		}
	}
}

func (bs *BookStats) printStats() {
	type ToSort struct {
		char  string
		count int
	}

	sortedSlice := make([]ToSort, 0)

	fmt.Println("--- Report of book \"frankenstein.txt\" ---")
	fmt.Printf("A total of %d words were found on this book\n", bs.words)
	fmt.Println()

	for char, count := range bs.chars {
		ts := ToSort{
			char,
			count,
		}

		sortedSlice = append(sortedSlice, ts)

	}
	sort.Slice(sortedSlice, func(i, j int) bool {
		return sortedSlice[i].count > sortedSlice[j].count
	})

	for _, chars := range sortedSlice {
		fmt.Printf("The character %q was found a total of %d times\n", chars.char, chars.count)
	}

	fmt.Println()
	fmt.Println("--- End Report of book \"frankenstein.txt\" ---")
}

func main() {

	file, err := os.Open("books/frankenstein.txt")
	assertNotNil(err, "Could not read file")
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := new(bytes.Buffer)

	bs := BookStats{
		0,
		map[string]int{},
	}
	for {

		bytes := make([]byte, bufferSize)
		_, err := reader.Read(bytes)

		if err == io.EOF {
			fmt.Println("End of file.")
			break
		}

		// Copy to the buffer to make use of the utility methods (read/WriteString)
		buffer.Write(bytes)
		for {
			line, err := buffer.ReadString('\n')

			if err == io.EOF {
				// We might have read a bunch but not found the delimiter, in which case
				// the buffer would return EOF. So we put back the line that was
				// read back into the buffer and continue to the next chunk.
				buffer.WriteString(line)
				break
			}

			bs.countWords(line)
			bs.countChars(line)
		}

	}

	bs.printStats()
}

func assertNotNil(err error, message string) {
	if err != nil {
		log.Fatalf(message+": %q", err)
	}
}
