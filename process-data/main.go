package main

import (
	"fmt"
	"os"
)

func main() {
	// read file from disk
	filename := "test.txt"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read file until space character
	var buffer [1]byte

	words := [3]string{}
	wordIndex := 0
	for {

		_, err := file.Read(buffer[:])
		if err != nil {
			sendWords(words, wordIndex)
			break
		}

		words[wordIndex] += string(buffer[:])

		if buffer[0] == ' ' {
			sendWords(words, wordIndex)
			wordIndex = (wordIndex + 1) % 3
			words[wordIndex] = ""
		}
	}

}

func sendWords(words [3]string, wordIndex int) {
	// print words
	for i := 0; i < len(words); i++ {
		fmt.Printf("%s", words[(i+wordIndex+1)%len(words)])
	}
	fmt.Println()
}
