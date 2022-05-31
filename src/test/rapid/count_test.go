package rapid

import (
	"testing"

	"github.com/polarode/hska-go-quickcheck/src/stringutil"
	"github.com/polarode/hska-go-quickcheck/src/testable"
	"pgregory.net/rapid"
)

func TestReverseSameNumberOfWords(t *testing.T) {
	property := func(t *rapid.T) {
		x := rapid.StringMatching("[a-zA-Z ]+").Draw(t, "words").(string)
		y := testable.Count(x)
		y2 := testable.Count(stringutil.Reverse(x))
		if y != y2 {
			t.Error("falsified: reverse string has the same number of words")
		}
	}
	rapid.Check(t, property)
}

func TestNumberOfWordsGoeZero(t *testing.T) {
	property := func(t *rapid.T) {
		x := rapid.StringMatching("[a-zA-Z ]+").Draw(t, "words").(string)
		y := testable.Count(x)
		if !(y >= 0) {
			t.Error("falsified: count is greater or equal than zero")
		}
	}
	rapid.Check(t, property)
}

func TestConcatenatedStringDoubleNumberOfWords(t *testing.T) {
	property := func(t *rapid.T) {
		x := rapid.StringMatching("[a-zA-Z ]+").Draw(t, "words").(string)
		y := testable.Count(x)
		y2 := testable.Count(x + x)
		if y*2 != y2 {
			t.Error("falsified: concatenated string has double number of words")
		}
	}
	rapid.Check(t, property)
}
