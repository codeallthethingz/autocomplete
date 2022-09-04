package ngramify

import (
	"io"
	"regexp"
	"strings"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
)

type Ngramify struct {
	sentenceTokenizer sentences.SentenceTokenizer
	w                 io.Writer
}

// create new ngramify
func New(w io.Writer) (*Ngramify, error) {
	b, err := data.Asset("data/english.json")
	if err != nil {
		return nil, err
	}
	training, err := sentences.LoadTraining(b)
	if err != nil {
		return nil, err
	}
	return &Ngramify{
		sentenceTokenizer: sentences.NewSentenceTokenizer(training),
		w:                 w,
	}, nil
}

// ngramify takes a string and write a all the ngrams to the writer
func (n *Ngramify) Ngramify(s string, i int) error {
	ngrams := []string{}
	sentences := n.sentenceTokenizer.Tokenize(s)
	for _, sentence := range sentences {
		// replace all non-alphanumeric characters with spaces using a regex
		regexp, err := regexp.Compile("[^a-zA-Z0-9^']+")
		if err != nil {
			return err
		}
		text := regexp.ReplaceAllString(sentence.Text, " ")

		words := strings.Fields(text)
		for l, word := range words {
			ngrams = append(ngrams, word)
			ngram := word
			for j := 1; j < i && len(words) > l+j; j++ {
				ngram += " " + words[l+j]
				ngrams = append(ngrams, ngram)
			}
		}
	}
	b := []byte(strings.Join(ngrams, "\n"))
	_, err := n.w.Write(b)
	return err
}
