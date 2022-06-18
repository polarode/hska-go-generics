package quick

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/polarode/hska-go-quickcheck/src/stringutil"
	"github.com/polarode/hska-go-quickcheck/src/testable"
)

func TestReverseSameNumberOfWords(t *testing.T) {
	property := func(x string) bool {
		y := testable.Count(x)
		y2 := testable.Count(stringutil.Reverse(x))
		return y == y2
	}
	config := quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = RandomStringGenerator(r, 16, "abcxyz ")
		}}
	if err := quick.Check(property, &config); err != nil {
		t.Error("falsified: reverse string has the same number of words", err)
	}
}

func TestNumberOfWordsGoeZero(t *testing.T) {
	property := func(x string) bool {
		y := testable.Count(x)
		return y >= 0
	}
	config := quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = RandomStringGenerator(r, 16, "abcxyz ")
		}}
	if err := quick.Check(property, &config); err != nil {
		t.Error("falsified: count is greater or equal than zero", err)
	}
}

func TestConcatenatedStringDoubleNumberOfWords(t *testing.T) {
	property := func(x string) bool {
		y := testable.Count(x)
		y2 := testable.Count(x + x)
		return y*2 == y2
	}
	config := quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = RandomStringGenerator(r, 16, "abcxyz ")
		}}
	if err := quick.Check(property, &config); err != nil {
		t.Error("falsified: concatenated string has double number of words", err)
	}
}
