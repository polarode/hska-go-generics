package rapid

import (
	"testing"

	"github.com/polarode/hska-go-quickcheck/src/testable"
	"pgregory.net/rapid"
)

func TestAddSymmetric(t *testing.T) {
	propertyInt := func(t *rapid.T) {
		a := rapid.Int64().Draw(t, "a").(int64)
		b := rapid.Int64().Draw(t, "b").(int64)
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		if y != y2 {
			t.Error("falsified: add is symmetric (int64)")
		}
	}
	propertyFloat := func(t *rapid.T) {
		a := rapid.Float64().Draw(t, "a").(float64)
		b := rapid.Float64().Draw(t, "b").(float64)
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		if y != y2 {
			t.Error("falsified: add is symmetric (float64)")
		}
	}
	rapid.Check(t, propertyInt)
	rapid.Check(t, propertyFloat)
}
