package ngramify

import (
	"io"
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
	sentences := n.sentenceTokenizer.Tokenize(s)
	for _, sentence := range sentences {
		text := strings.ReplaceAll(sentence.Text, ",", "")
		text = strings.ReplaceAll(text, ".", "")
		words := strings.Fields(text)
		for l, word := range words {
			if _, err := n.w.Write([]byte(word + "\n")); err != nil {
				return err
			}
			ngram := word
			for j := 1; j < i && len(words) > l+j; j++ {
				ngram += " " + words[l+j]
				if _, err := n.w.Write([]byte(ngram + "\n")); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
