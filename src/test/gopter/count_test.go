package gopter

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/polarode/hska-go-quickcheck/src/stringutil"
	"github.com/polarode/hska-go-quickcheck/src/testable"
)

func TestReverseSameNumberOfWords(t *testing.T) {
	properties := gopter.NewProperties(nil)
	property := func(x string) bool {
		y := testable.Count(x)
		y2 := testable.Count(stringutil.Reverse(x))
		return y == y2
	}
	properties.Property("reverse string has the same number of words",
		prop.ForAll(property, gen.RegexMatch("[a-zA-Z ]+")))
	properties.TestingRun(t)
}

func TestNumberOfWordsGoeZero(t *testing.T) {
	properties := gopter.NewProperties(nil)
	property := func(x string) bool {
		y := testable.Count(x)
		return y >= 0
	}
	properties.Property("count is greater or equal than zero",
		prop.ForAll(property, gen.RegexMatch("[a-zA-Z ]+")))
	properties.TestingRun(t)
}

func TestConcatenatedStringDoubleNumberOfWords(t *testing.T) {
	properties := gopter.NewProperties(nil)
	property := func(x string) bool {
		y := testable.Count(x)
		y2 := testable.Count(x + x)
		return y*2 == y2
	}
	properties.Property("concatenated string has double number of words",
		prop.ForAll(property, gen.RegexMatch("[a-zA-Z ]+")))
	properties.TestingRun(t)
}
