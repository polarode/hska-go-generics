package gopter

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/polarode/hska-go-quickcheck/src/testable"
)

func TestAddSymmetric(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("add is symmetric (int64)", prop.ForAll(
		func(a, b int64) bool {
			y := testable.Add(a, b)
			y2 := testable.Add(b, a)
			return y == y2
		},
		gen.Int64(),
		gen.Int64(),
	))
	properties.Property("add is symmetric (float64)", prop.ForAll(
		func(a, b float64) bool {
			y := testable.Add(a, b)
			y2 := testable.Add(b, a)
			return y == y2
		},
		gen.Float64(),
		gen.Float64(),
	))
	properties.TestingRun(t)
}
