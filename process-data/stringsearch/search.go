package stringsearch

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/dghubble/trie"
)

type AutocompleteTrie struct {
	trie *trie.RuneTrie
}

func NewAutocompleteTrie(reader io.Reader) *AutocompleteTrie {
	trie := trie.NewRuneTrie()
	// read lines
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// split text and value
		valueAndText := strings.SplitN(line, " ", 2)
		value := math.atoi(valueAndText[0])
		text := valueAndText[1]
		// for every character in the line add it to the trie
		for i := 1; i < len(text); i++ {
			if trie.Get(line[:i]) == nil {
				trie.Put(line[:i], []string{})
			}
			currentValue := trie.Get(line[0:i]).([]string)
			trie.Put(line[0:i], append(currentValue, line))
			fmt.Printf("putting %s, %s\n", line[0:i], trie.Get(line[0:i]))
		}
	}

	at := &AutocompleteTrie{
		trie: trie,
	}
	return at

}

func (at *AutocompleteTrie) Find(prefix string) ([]string, bool) {
	value := at.trie.Get(prefix)
	if value == nil {
		return nil, false
	}
	return value.([]string), true
}
