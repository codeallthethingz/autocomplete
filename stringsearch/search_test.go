package stringsearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCorrect(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 hello world\n2 help me\n3 hell freezes over\n4 hello kitty\n5 hello darkness my old friend\n6 hard to say")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("helo", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, []string{"hello world", "hello kitty", "hello darkness my old friend"}, []string{values[0].Text, values[1].Text, values[2].Text})
}

// test spell correct where correction is underneath original word relevance
func TestSpellCorrect(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 helo sun\n3 hello world")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("helo", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, []string{"helo sun", "hello world"}, []string{values[0].Text, values[1].Text})
}

// test whitespace left removal
func TestWhitespaceLeft(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 red fox\n2 red")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig(" red  ", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, 1, len(values))
	require.Equal(t, []string{"red fox"}, []string{values[0].Text})
}

// test case sensitive returns case match as first result
func TestCaseSensitive(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 RED\n2 red")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("R", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, []string{"RED", "red"}, []string{values[0].Text, values[1].Text})

	values, ok = test.FindWithConfig("R", Config{RelevanceCaseAware: false})
	require.True(t, ok)
	require.Equal(t, []string{"red", "RED"}, []string{values[0].Text, values[1].Text})
}

// test exact match better than prefix match
func TestExactMatch(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 red\n2 Reddit")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("r", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, []string{"red", "Reddit"}, []string{values[0].Text, values[1].Text})
	values, ok = test.FindWithConfig("red", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, []string{"red", "Reddit"}, []string{values[0].Text, values[1].Text})
}

// test mixed case
func TestMixedCase(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 Hello world\n2 help me\n3 hell freezes over\n4 heLlo kitty\n")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("HeLl", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, 3, len(values))
}

// test edges
func TestEdges(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 hello world")
	test := NewAutocompleteTrie(reader, 5)
	if _, ok := test.FindWithConfig("k", DefaultConfig); ok {
		t.Error("should not find k")
	}
	if _, ok := test.FindWithConfig("", DefaultConfig); ok {
		t.Error("should not find empty value")
	}

	reader = strings.NewReader("1 red\n2 red\n3 red")
	test = NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("red", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, "red", values[0].Text)
	require.Equal(t, 3, values[0].Count)

}

// test finding a word in a string using trie
func TestTrie(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 hello world\n2 help me\n")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("h", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, []string{"help me", "hello world"}, []string{values[0].Text, values[1].Text})
	values, ok = test.FindWithConfig("hell", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, "hello world", values[0].Text)
}

// test max N values
func TestMaxN(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 hello world\n2 help me\n3 hell freezes over\n4 hello kitty\n5 hello darkness my old friend\n6 hard to say")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindWithConfig("h", DefaultConfig)
	require.True(t, ok)
	require.Equal(t, len(values), 5)
}
