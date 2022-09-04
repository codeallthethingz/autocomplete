package stringsearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// test finding a word in a string using trie
func TestTrie(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 hello world\n2 help me\n")
	test := NewAutocompleteTrie(reader)
	values, ok := test.Find("h")
	require.True(t, ok)
	require.Equal(t, []string{"help me", "hello world"}, values)
	values, ok = test.Find("hell")
	require.True(t, ok)
	require.Equal(t, []string{"hello world"}, values)
	t.Fail()
}
