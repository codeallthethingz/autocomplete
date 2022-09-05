package stringsearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// test case sensitive returns case match as first result
func TestCaseSensitive(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 RED\n2 red")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.FindCaseAware("R")
	require.True(t, ok)
	require.Equal(t, []string{"RED", "red"}, []string{values[0].Text, values[1].Text})

	values, ok = test.Find("r")
	require.True(t, ok)
	require.Equal(t, []string{"red", "RED"}, []string{values[0].Text, values[1].Text})
}

// test mixed case
func TestMixedCase(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 Hello world\n2 help me\n3 hell freezes over\n4 heLlo kitty\n")
	test := NewAutocompleteTrie(reader, 5)
	values, ok := test.Find("HeLl")
	require.True(t, ok)
	require.Equal(t, 3, len(values))
}

// test edges
func TestEdges(t *testing.T) {
	// create string reader
	reader := strings.NewReader("1 hello world")
	test := NewAutocompleteTrie(reader, 5)
	if _, ok := test.Find("k"); ok {
		t.Error("should not find k")
	}
	if _, ok := test.Find(""); ok {
		t.Error("should not find empty value")
	}
	if _, ok := test.FindCaseAware("k"); ok {
		t.Error("should not find k")
	}
	if _, ok := test.FindCaseAware(""); ok {
		t.Error("should not find empty value")
	}

}

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
