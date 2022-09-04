package ngramify

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var testNgramify *Ngramify
var buf *bytes.Buffer
var ngramLen = 3

func setup() {
	buf = bytes.NewBuffer([]byte{})
	n, err := New(buf)
	if err != nil {
		panic(err)
	}
	testNgramify = n
}

func test(t *testing.T, text string) []string {
	err := testNgramify.Ngramify(text, ngramLen)
	require.NoError(t, err)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func TestBracketedText(t *testing.T) {
	setup()
	expected := []string{
		"foo",
		"foo bar",
		"foo bar baz",
		"bar",
		"bar baz",
		"baz",
	}
	actual := test(t, "foo (bar  baz)")
	require.Equal(t, expected, actual)
}

func TestOneWord(t *testing.T) {
	setup()
	require.Equal(t, []string{"a"}, test(t, "a"))
}
func TestTwoWords(t *testing.T) {
	setup()
	require.Equal(t, []string{"a", "a b", "b"}, test(t, "a b"))
}
func TestThreeWords(t *testing.T) {
	setup()
	require.Equal(t, []string{"a", "a b", "a b c", "b", "b c", "c"}, test(t, "a b c"))
}
func TestFourWords(t *testing.T) {
	setup()
	require.Equal(t, []string{"a", "a b", "a b c", "b", "b c", "b c d", "c", "c d", "d"}, test(t, "a b c d"))
}
func TestCommas(t *testing.T) {
	setup()
	require.Equal(t, []string{"this", "this should", "this should work", "should", "should work", "work"}, test(t, "this, should work"))
}
func TestSentences(t *testing.T) {
	setup()
	require.Equal(t, []string{"This", "This should", "This should end",
		"should", "should end", "end",
		"then", "then start", "then start here",
		"start", "start here",
		"here"}, test(t, "This should end. then start here."))
}
func TestSentenceWithAncronym(t *testing.T) {
	setup()
	require.Equal(t, []string{
		"Dr", "Dr Roberts'", "Dr Roberts' desk",
		"Roberts'", "Roberts' desk",
		"desk"}, test(t, "Dr. Roberts' desk"))
}
