package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/searchspring/autocomplete/process-data/disksort"
	"github.com/searchspring/autocomplete/process-data/ngramify"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	inputText := "This is the first file and it's longer than you might expect.  Triumphant last sentence in this file.  This is the third sentence."
	var buf *bytes.Buffer = bytes.NewBuffer([]byte{})
	ngramify, err := ngramify.New(buf)
	require.NoError(t, err)
	ngramify.Ngramify(inputText, 3)
	require.Equal(t, 55, len(bytes.Split(buf.Bytes(), []byte(" "))))

	lines := strings.Split(buf.String(), "\n")
	stringChan := make(chan string, 2)
	go func() {
		for _, line := range lines {
			stringChan <- line
		}
		close(stringChan)
	}()
	buf.Reset()
	err = disksort.Sort(stringChan, buf, true)
	require.NoError(t, err)
	require.Equal(t, "you might expect", strings.Split(buf.String(), "\n")[51])
}
