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
	Count int    `json:"count"`
	Text  string `json:"text"`
}

// config options
type Config struct {
	RelevanceCaseAware     bool
	RelevanceExactMatch    bool
	RecallTrimLeadingSpace bool
}

var DefaultConfig = Config{true, true, true}

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
		return currentValue[i].Count > currentValue[j].Count
	})
	// if max values per entry is set, only keep the first N
	if maxValuesPerEntry > 0 && len(currentValue) > maxValuesPerEntry {
		currentValue = currentValue[:maxValuesPerEntry]
	}
	trie.Put(key, currentValue)
}

func (at *AutocompleteTrie) Find(prefix string) ([]AutocompleteTrieValue, bool) {
	return at.FindWithConfig(prefix, DefaultConfig)
}

func (at *AutocompleteTrie) FindWithConfig(prefix string, config Config) ([]AutocompleteTrieValue, bool) {
	if config.RecallTrimLeadingSpace {
		if strings.HasSuffix(prefix, " ") {
			prefix = strings.TrimSpace(prefix)
			prefix += " "
		} else {
			prefix = strings.TrimSpace(prefix)
		}
	}

	if strings.TrimSpace(prefix) == "" {
		return nil, false
	}

	values := at.trie.Get(strings.ToLower(prefix))
	if values != nil {
		// copy values array
		valuesCopy := make([]AutocompleteTrieValue, len(values.([]AutocompleteTrieValue)))
		copy(valuesCopy, values.([]AutocompleteTrieValue))
		sort.Slice(valuesCopy, func(i, j int) bool {
			iPoints := 0 // points for the i value
			jPoints := 0 // points for the j value

			// case sensitive prefix points
			if config.RelevanceCaseAware {
				if strings.HasPrefix(valuesCopy[i].Text, prefix) {
					iPoints += 1
				}
				if strings.HasPrefix(valuesCopy[j].Text, prefix) {
					jPoints += 1
				}
			}

			// case insensitive exact match
			if strings.EqualFold(valuesCopy[i].Text, prefix) {
				iPoints += 2
			}
			if strings.EqualFold(valuesCopy[j].Text, prefix) {
				jPoints += 2
			}

			// case sensitive exact match points
			if config.RelevanceCaseAware {
				if valuesCopy[i].Text == prefix {
					iPoints += 3
				}
				if valuesCopy[j].Text == prefix {
					jPoints += 3
				}
			}

			// if the points are equal sort by value
			if iPoints == jPoints {
				return valuesCopy[i].Count > valuesCopy[j].Count
			}
			// sort by points
			return iPoints > jPoints
		})
		return valuesCopy, true
	} else {
		return nil, false
	}
}
