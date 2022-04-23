package quick

import (
	"testing"
	"testing/quick"

	"github.com/polarode/hska-go-quickcheck/src/testable"
)

func TestAddInt(t *testing.T) {
	f := func(a, b int64) bool {
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		return y == y2
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}

	f2 := func(a, b float64) bool {
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		return y == y2
	}
	if err := quick.Check(f2, nil); err != nil {
		t.Error(err)
	}
}
