package disksort

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleArray(t *testing.T) {
	input := []string{"a", "c", "b", "a"}
	expected := []string{"a", "a", "b", "c"}
	actual := sortArray(t, input, false)
	require.Equal(t, expected, actual)
}

func TestSimpleArrayDedup(t *testing.T) {
	input := []string{"a", "c", "b", "a"}
	expected := []string{"a", "b", "c"}
	actual := sortArray(t, input, true)
	require.Equal(t, expected, actual)
}

func sortArray(t *testing.T, input []string, dedup bool) []string {
	buf := &bytes.Buffer{}
	stringChan := make(chan string, 2)
	go func() {
		for _, word := range input {
			stringChan <- word
		}
		close(stringChan)
	}()
	err := Sort(stringChan, buf, dedup)
	require.NoError(t, err)
	actual := strings.Split(strings.TrimSpace(buf.String()), "\n")
	return actual
}
