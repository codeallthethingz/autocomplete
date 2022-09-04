package disksort

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleArray(t *testing.T) {
	input := []string{"a", "c", "b", "a"}
	expected := []string{"2 a", "1 b", "1 c"}
	actual := sortArray(t, input)
	require.Equal(t, expected, actual)
}

func sortArray(t *testing.T, input []string) []string {
	buf := &bytes.Buffer{}
	stringChan := make(chan string, 2)
	go func() {
		for _, word := range input {
			stringChan <- word
		}
		close(stringChan)
	}()
	err := Sort(stringChan, buf)
	require.NoError(t, err)
	actual := strings.Split(strings.TrimSpace(buf.String()), "\n")
	return actual
}
