package quick

import (
	"testing"
	"testing/quick"

	"github.com/polarode/hska-go-quickcheck/src/testable"
)

func TestAddSymmetric(t *testing.T) {
	propertyInt := func(a, b int64) bool {
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		return y == y2
	}
	propertyFloat := func(a, b float64) bool {
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		return y == y2
	}
	if err := quick.Check(propertyInt, nil); err != nil {
		t.Error("falsified: add is symmetric (int64)", err)
	}
	if err := quick.Check(propertyFloat, nil); err != nil {
		t.Error("falsified: add is symmetric (float64)", err)
	}
}
