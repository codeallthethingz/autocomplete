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
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.Find("h")
	require.True(t, ok)
	require.Equal(t, []string{"help me", "hello world"}, []string{values[0].Text, values[1].Text})
	values, ok = test.Find("hell")
	require.True(t, ok)
	require.Equal(t, "hello world", values[0].Text)
}

// test max N values
func TestMaxN(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 hello world\n2 help me\n3 hell freezes over\n4 hello kitty\n5 hello darkness my old friend\n6 hard to say")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.Find("h")
	require.True(t, ok)
	require.Equal(t, len(values), 5)
}
