package stringsearch

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/dghubble/trie"
)

type AutocompleteTrie struct {
	trie *trie.RuneTrie
}

type AutocompleteTrieValue struct {
	Value int    `json:"value"` // ignore value in json
	Text  string `json:"text"`
}

func NewAutocompleteTrie(reader io.Reader, maxValuesPerEntry int) *AutocompleteTrie {
	trie := trie.NewRuneTrie()
	// read lines
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// split text and value
		valueAndText := strings.SplitN(line, " ", 2)
		// convert value to int atoi
		value, _ := strconv.Atoi(valueAndText[0])
		text := valueAndText[1]
		// for every character in the line add it to the trie
		for i := 1; i < len(text)+1; i++ {
			addAndSortValue(trie, strings.ToLower(text[:i]), AutocompleteTrieValue{value, text}, maxValuesPerEntry)
		}
	}

	at := &AutocompleteTrie{
		trie: trie,
	}
	return at

}

func addAndSortValue(trie *trie.RuneTrie, key string, newValue AutocompleteTrieValue, maxValuesPerEntry int) {
	if trie.Get(key) == nil {
		trie.Put(key, []AutocompleteTrieValue{})
	}
	currentValue := trie.Get(key).([]AutocompleteTrieValue)
	currentValue = append(currentValue, newValue)
	// sort the autocompletetrievalues
	sort.Slice(currentValue, func(i, j int) bool {
		return currentValue[i].Value > currentValue[j].Value
	})
	// if max values per entry is set, only keep the first N
	if maxValuesPerEntry > 0 && len(currentValue) > maxValuesPerEntry {
		currentValue = currentValue[:maxValuesPerEntry]
	}
	trie.Put(key, currentValue)
}

func (at *AutocompleteTrie) Find(prefix string) ([]AutocompleteTrieValue, bool) {
	prefix = strings.ToLower(strings.TrimSpace(prefix))

	if prefix == "" {
		return nil, false
	}
	value := at.trie.Get(prefix)
	if value == nil {
		return nil, false
	}
	return value.([]AutocompleteTrieValue), true
}

func (at *AutocompleteTrie) FindCaseAware(prefix string) ([]AutocompleteTrieValue, bool) {
	prefix = strings.TrimSpace(prefix)

	if prefix == "" {
		return nil, false
	}

	// find all values for the lowercase prefix and then sort by the original prefix case

	if values, ok := at.Find(strings.ToLower(prefix)); ok {
		// copy values array
		valuesCopy := make([]AutocompleteTrieValue, len(values))
		copy(valuesCopy, values)
		sort.Slice(valuesCopy, func(i, j int) bool {
			iPoints := 0 // points for the i value
			jPoints := 0 // points for the j value

			if strings.HasPrefix(valuesCopy[i].Text, prefix) {
				iPoints += 1
			}
			if strings.HasPrefix(valuesCopy[j].Text, prefix) {
				jPoints += 1
			}
			// if prefix string case insensitive matches add a point
			if strings.EqualFold(valuesCopy[i].Text, prefix) {
				iPoints += 3
			}
			if strings.EqualFold(valuesCopy[j].Text, prefix) {
				jPoints += 3
			}
			// if prefix string case sensitive matches add a point
			if valuesCopy[i].Text == prefix {
				iPoints += 1
			}
			if valuesCopy[j].Text == prefix {
				jPoints += 1
			}

			// if the points are equal sort by value
			if iPoints == jPoints {
				return valuesCopy[i].Value < valuesCopy[j].Value
			}
			// sort by points
			return iPoints > jPoints
		})
		return valuesCopy, true
	} else {
		return nil, false
	}
}
