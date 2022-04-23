# Project work: Go Quickcheck and Generics

## Essentials

### Quickcheck in Haskell

> todo: summary on basics of Quickcheck using examples in Haskell

### Generics in Go 1.18

> todo: current state of in Go v1.18 (some functionality was moved to version 1.19)

At the time of this project work, version 1.18 of Go is the most recent. In this version generics are not fully supported. The developers decided to push some functionality to a later version. Here is a short overview, what is included in version 1.18 and what is still to come:

Implemented in version 1.18:
- function and type declarartions now allow typeparameters (square bracket)
- interfaces can be used as type contraints: an interface can define a set of types
- predefined identifier "any": alias for the empty interface (matches all types)
- predefined identifier "comparable": set of all types, that can be compared with  == and !=

Missing in version 1.18:
- typedeclarations within generic functions
- only methods, that are explicitly declared in the interface of the type parameter, can be used within the the method
- accessing struct fields is not possible
> source: [release notes v1.18](https://go.dev/doc/go1.18) 

## Quickcheck in Go

### Existing work

> todo: which libraries do already exist for quickcheck? What are their differences? How good/useful are they? Do they support generics?

There exist three popular libraries, that allow property based testing, similar to that of Quickcheck in Haskell.
- [testing/quick](https://pkg.go.dev/testing/quick)
- [github.com/leanovate/gopter](https://pkg.go.dev/github.com/leanovate/gopter)
- [pgregory.net/rapid](https://pkg.go.dev/pgregory.net/rapid)

The library testing/quick is even a official package that comes with Go and is relativly easy to use. However, because it is kept rather simple, it is also limited in its functionality. For example it does not support shrinkage.
> todo: explain shrinkage in essentials

gopter is described as a more sophisticated version of the testing/quick package. For exmpample it allows shrinkage, it allows better control of the generators and there is support for stateful tests.

rapid is similar to gopter. It also allows shrinkage, better generators and stateful tests. It claims to have a simpler API than gopter and it does not require user code to minimize failing tests.

#### Comparison 

> todo: compare implementations for example tests

**Non-generic functions**

As an example for non-generic functions, a function `count`, that counts the words (separated by a whitespace), is tested here. The property that is beeing tested is, that the reversed string and the original string have the same count of words.
This property can be tested in Haskell using quickcheck with this property:

```haskell
prop :: String -> Bool
prop s = count s == count (reverse s)
```

The same property can be tested in Go with the following implementations for the different libraries:

testing/quick:
```go
func TestReverseSameNumberOfWords(t *testing.T) {
	f := func(x string) bool {
		y := testable.Count(x)
		y2 := testable.Count(stringutil.Reverse(x))
		return y == y2
	}
	config := quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(RandomStringGenerator(r, 16, "abcxyz "))
		}}
	if err := quick.Check(f, &config); err != nil {
		t.Error(err)
	}
}
```

gopter:
```go
func TestReverseSameNumberOfWords(t *testing.T) {
	properties := gopter.NewProperties(nil)
	properties.Property("reverse string has the same number of words", prop.ForAll(
		func(x string) bool {
			y := testable.Count(x)
			y2 := testable.Count(stringutil.Reverse(x))
			return y == y2
		},
		gen.RegexMatch("[a-zA-Z ]+"),
	))
	properties.TestingRun(t)
}
```

rapid:
```go
func TestReverseSameNumberOfWords(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := rapid.StringMatching("[a-zA-Z ]+").Draw(t, "words").(string)
		y := testable.Count(x)
		y2 := testable.Count(stringutil.Reverse(x))
		if y != y2 {
			t.Fatal("falsified: reverse string has the same number of words")
		}
	})
}
```

**Generic functions**

For generic functions we test a simple add function, that can add two numbers of type int or float. In this example we test for symmetry of this function.

This property can be tested in Haskell using quickcheck with this property:
> todo: Haskell aquivalent

The same property can be tested in Go with the following implementations for the different libraries:

testing/quick:
```go
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
```

gopter:
```go
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
```

rapid:
```go
func TestAddSymmetric(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		a := rapid.Int64().Draw(t, "a").(int64)
		b := rapid.Int64().Draw(t, "b").(int64)
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		if y != y2 {
			t.Fatal("falsified: add is symmetric (int64)")
		}
	})
	rapid.Check(t, func(t *rapid.T) {
		a := rapid.Float64().Draw(t, "a").(float64)
		b := rapid.Float64().Draw(t, "b").(float64)
		y := testable.Add(a, b)
		y2 := testable.Add(b, a)
		if y != y2 {
			t.Fatal("falsified: add is symmetric (float64)")
		}
	})
}
```

### Own implementation of quickcheck in Go (if existing libraries don't meet expectations)

> todo