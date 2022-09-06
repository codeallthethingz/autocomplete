package stringsearch

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/dghubble/trie"
	"github.com/sajari/fuzzy"
)

type AutocompleteTrie struct {
	trie  *trie.RuneTrie
	model *fuzzy.Model
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
	RecallSpellCorrection  bool
}

var DefaultConfig = Config{true, true, true, true}

func NewAutocompleteTrie(reader io.Reader, maxValuesPerEntry int) *AutocompleteTrie {
	trie := trie.NewRuneTrie()
	model := fuzzy.NewModel()
	model.SetDepth(6)
	model.SetThreshold(1)
	trainingSet := []string{}
	// read lines
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// split text and value
		valueAndText := strings.SplitN(line, " ", 2)
		// convert value to int atoi
		value, _ := strconv.Atoi(valueAndText[0])
		text := valueAndText[1]
		// add each word in text to the training set
		words := strings.Split(strings.ToLower(text), " ")
		if len(words) > 1 {
			trainingSet = append(trainingSet, words...)
		}
		// for every character in the line add it to the trie
		for i := 1; i < len(text)+1; i++ {
			addAndSortValue(trie, strings.ToLower(text[:i]), AutocompleteTrieValue{value, text}, maxValuesPerEntry)
		}
	}

	model.Train(trainingSet)

	at := &AutocompleteTrie{
		trie:  trie,
		model: model,
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
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	// excute search in a goroutine for the main search
	values := []AutocompleteTrieValue{}
	go func() {
		defer waitGroup.Done()
		if v, ok := at.executeSearch(prefix, config); ok {
			values = append(values, v...)
		}
	}()

	correctedValues := []AutocompleteTrieValue{}
	if config.RecallSpellCorrection {
		correctedTerms := at.generateCorrectedTerms(prefix, 10)
		waitGroup.Add(len(correctedTerms))
		for _, corrected := range correctedTerms {
			go func(corrected string) {
				defer waitGroup.Done()
				// execute search for the spell correction
				if v, ok := at.executeSearch(corrected, config); ok {
					correctedValues = append(correctedValues, v...)
				}
			}(corrected)
		}
	}

	// wait for all goroutines to finish
	waitGroup.Wait()

	if len(values) == 0 {
		// sort the corrected values by count
		sort.Slice(correctedValues, func(i, j int) bool {
			return correctedValues[i].Count > correctedValues[j].Count
		})

		// dedup the corrected values
		seen := map[string]bool{}
		dedupedCorrectedValues := []AutocompleteTrieValue{}
		for _, v := range correctedValues {
			if _, ok := seen[v.Text]; !ok {
				seen[v.Text] = true

				dedupedCorrectedValues = append(dedupedCorrectedValues, v)
			}
		}

		// truncate deduped to max 10
		if len(dedupedCorrectedValues) > 10 {
			dedupedCorrectedValues = dedupedCorrectedValues[:10]
		}

		// append the deduped values to the main search
		values = append(values, dedupedCorrectedValues...)
	}
	// if empty return false
	if len(values) == 0 {
		return nil, false
	}
	return values, true
}

var alreadyCorrected map[string]string = map[string]string{}

func (at *AutocompleteTrie) generateCorrectedTerms(searchTerms string, max int) []string {
	results := []string{searchTerms}
	words := strings.SplitN(searchTerms, " ", 3)
	for i := 0; i < len(words); i++ {
		corrected := ""
		for j := 0; j < i+1; j++ {
			lowerWord := strings.ToLower(words[j])
			if _, ok := alreadyCorrected[lowerWord]; !ok {
				alreadyCorrected[lowerWord] = at.model.SpellCheck(lowerWord)
			}
			corrected += alreadyCorrected[lowerWord] + " "
		}
		for j := i + 1; j < len(words); j++ {
			corrected += words[j] + " "
		}
		results = append(results, strings.TrimSpace(corrected))
	}

	// dedup the results
	seen := map[string]bool{}
	dedupedResults := []string{}
	for _, v := range results {
		if _, ok := seen[v]; !ok {
			seen[v] = true
			dedupedResults = append(dedupedResults, v)
		}
	}

	// truncate deduped to max
	if len(dedupedResults) > max {
		dedupedResults = dedupedResults[:max]
	}

	return dedupedResults
}

func (at *AutocompleteTrie) executeSearch(prefix string, config Config) ([]AutocompleteTrieValue, bool) {
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
