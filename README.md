<h1>Project work: Go QuickCheck and Generics</h1>

<h2>Table of contents</h2>

- [1. Essentials](#1-essentials)
	- [1.1 QuickCheck in Haskell](#11-quickcheck-in-haskell)
		- [1.1.1 Generators](#111-generators)
		- [1.1.2 Properties](#112-properties)
	- [1.2 Generics in Go 1.18](#12-generics-in-go-118)
- [2. QuickCheck in Go](#2-quickcheck-in-go)
	- [2.1 Existing work](#21-existing-work)
		- [2.1.1 Comparison non-generic functions](#211-comparison-non-generic-functions)
		- [2.1.2 Comparison generic functions](#212-comparison-generic-functions)
	- [2.2 Own implementation of quickcheck in Go (if existing libraries don't meet expectations)](#22-own-implementation-of-quickcheck-in-go-if-existing-libraries-dont-meet-expectations)

## 1. Essentials

Normal unit tests check the code, they are written for, based on explicitly provided examples. Only the examples, that are provided in the form of unit tests by a developer, can be tested.

Property based testing, which originated in Haskell's library QuickCheck, uses another approach. Instead of explicit examples, generators are provided, that can generate valid input data. These generators allow us to automate the test inputs.
To also automate the validataion of the tested inputs, properties are specified, that need to hold for all valid inputs.
With both the generation of input data and the validation of results automated, a function can be easily tested for a big number of random inputs, without the need for a developer to specify them in unit tests.

This should not be used as the only testing tool in a project, but it can be useful to supplement regular testing methods.

### 1.1 QuickCheck in Haskell

The following Haskell code specifies a function `count`, which can count the number of words within a string of characters. Each word is a sequence of characters, that is not separated by white spaces.

```haskell
count :: String -> Int
count [] = 0
count (c:cs)
  | c == ' ' = count $ skipBlanks cs
  | otherwise = 1 + count (skipWord cs)

-- | Generic skip function.
skip :: (Char -> Bool) -> String -> String
skip p [] = []
skip p (c:cs)
 | p c       = skip p cs
 | otherwise = c:cs

skipWord   = skip (/= ' ')
skipBlanks = skip (== ' ')
```

This will be the function to be tested using QuickCheck.

#### 1.1.1 Generators

But first the generator for the inputs of this function is needed. For QuickCheck this can be done by providing an instance for the class `Arbitrary` for the required type (in this case `String`). For the some common types, these instances are already given, but if they are missing or for new types, they could be provided as follows:

```haskell
instance Arbitrary Char where 
   arbitrary = elements (['a'..'z'] ++ [' '])

instance Arbitrary String where
   arbitrary = listOf arbitrary
```
The funtions `elements` and `listOf` used here are defined by the QuickCheck package. `elements` allows to generate a random single value from a list of possible values, while `listOf` allows to generate a list of random elements. The implementation can be looked up in the [QuickCheck Source](https://hackage.haskell.org/package/QuickCheck-2.14.2/docs/src/Test.QuickCheck.Gen.html):
```haskell
-- | Generates one of the given values. The input list must be non-empty.
elements :: [a] -> Gen a
elements [] = error "QuickCheck.elements used with empty list"
elements xs = (xs !!) `fmap` chooseInt (0, length xs - 1)

-- | Generates a list of random length. The maximum length depends on the
-- size parameter.
listOf :: Gen a -> Gen [a]
listOf gen = sized $ \n ->
  do k <- chooseInt (0,n)
     vectorOf k gen
```

#### 1.1.2 Properties

With the generators provided by `arbitrary` we can define properties. The example funtion from above was `count`, which allows to counts words in a string. One expectation to this function is, that for the reverse of a string `count` should yield the same number of words as for the string itself. This can be expressed with this property:

```haskell
-- | Reversing the string yields the same number of words.
prop :: String -> Bool
prop s = count s == count (reverse s)
```

The defined property can then be tested for example via the terminal with `quickCheck prop`. QuickCheck will check the property 100 times by generating random inputs for its funcitions according to the instances of Arbitrary, to see if the property holds.
If these tests pass, the result might look like this:
```
*Main> quickCheck prop
+++ OK, passed 100 tests.
```

If one test fails, quickCheck will use *shrinkage* to try and find simpler examples, for wich the test also fails. This can make it easier for a human to debug the failur.

Such a failur can be found with this property:
```haskell
-- | Concatenating the string doubles the number of words.
-- NOTE: This property does not hold in general!
prop2 :: String -> Bool
prop2 s = 2 * count s == count (s ++ s)
```

The result of executing QuickCheck for this property might look like this:
```
*Main> quickCheck prop2
*** Failed! Falsifiable (after 5 tests and 1 shrinks):    
"c"
``` 

### 1.2 Generics in Go 1.18

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

For the full changelog, see the [release notes for v1.18](https://go.dev/doc/go1.18).

## 2. QuickCheck in Go

The concept of QuickCheck (or property based testing) has also been adopted to Go to some extend. This will be analized in the following section.

### 2.1 Existing work

There exist three popular libraries, that allow property based testing, similar to that of Quickcheck in Haskell.
- [testing/quick](https://pkg.go.dev/testing/quick)
- [github.com/leanovate/gopter](https://pkg.go.dev/github.com/leanovate/gopter)
- [pgregory.net/rapid](https://pkg.go.dev/pgregory.net/rapid)

The library testing/quick is an official package that comes with Go and is relativly easy to use. However, because it is kept rather simple, it is also limited in its functionality. For example it does not support shrinkage.

gopter is described as a more sophisticated version of the testing/quick package. For exmpample it allows shrinkage, it allows better control of the generators and there is support for stateful tests.

rapid is similar to gopter. It also allows shrinkage, better generators and stateful tests. It claims to have a simpler API than gopter and it does not require user code to minimize failing tests.

#### 2.1.1 Comparison non-generic functions

As an example for non-generic functions, the Go version of the function `Count` is tested here. It is implemented as follows:

```go
func skip(p func(byte) bool) func(string) string {
	return func(s string) string {
		switch {
		case len(s) == 0:
			return ""
		case p(s[0]):
			return skip(p)(s[1:])
		default:
			return s
		}
	}
}
func Count(s string) int {
	skipBlanks := skip(func(b byte) bool {
		return b == ' '
	})

	skipWord := skip(func(b byte) bool {
		return b != ' '
	})
	switch {
	case len(s) == 0:
		return 0
	case s[0] == ' ':
		return Count(skipBlanks(s))
	default:
		return 1 + Count(skipWord(s))
	}
}
```

The property that is beeing tested is again, that the reversed string and the original string have the same count of words.
As a quick reminder, this can be tested in Haskell using quickcheck with this property:

```haskell
-- | Reversing the string yields the same number of words.
prop :: String -> Bool
prop s = count s == count (reverse s)
```

The same property can be tested in Go with the following implementations for the different libraries:

**testing/quick:**
```go
func TestReverseSameNumberOfWords(t *testing.T) {
	property := func(x string) bool {
		y := testable.Count(x)
		y2 := testable.Count(stringutil.Reverse(x))
		return y == y2
	}
	config := quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(RandomStringGenerator(r, 16, "abcxyz "))
		}}
	if err := quick.Check(property, &config); err != nil {
		t.Error("falsified: reverse string has the same number of words", err)
	}
}
```

The property that is supposed to be tested is provided as a function. The parameters of this function are used as input values for the functions that is being tested. Those values can be generated outside of the property function. For that testing/quick needs an additional configuration, that contains the generators for the input values of the property function. Since testing/quick does not come with a generator for strings, this has to be implemented manually:

```go
func RandomStringGenerator(r *rand.Rand, size int, alphabet string) string {
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := r.Intn(len(alphabet))
		buffer.WriteString(string(alphabet[index]))
	}
	return buffer.String()
}
```
The property function is returning a boolean value, that contains the result for the tested parameters. If a set of parameters evaluates to false, testing/quick handles the error and prints according error messages. This function is evaluated multiple times random inputs according to the given configuration on the call `quick.Check(property, &config)`. This function can return an error, if the check failed for a set of parameters. This error is then forwarded to the testing object `t` together with a describtion of the property.

**gopter:**
```go
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
```

For gopter the properties that should be tested are again provided as functions. The parameters of this function again serve as input of the function being tested and are generage outside using generators. In contrary to testing/quick, multiple properties can be provided at once and are tested in a testing run. For each property, generators have to be provided for each input of the function. gopter comes with some basic generators, that can be used here. In this example a generator is used, that can generate strings based on a regex.
The properties function returns a boolean value, that contains the result for the tested parameters. If a set of parameters evaluates to false, gopter handles the error and prints according error messages based on the provided name (first parameter of `properties.Property`). 
On the call `properties.TestingRun(t)` all the properties provided withing `properties` are evaluated multiple times with random inputs.

**rapid:**
```go
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
```

The property to be tested is again provided as a function. However, the parameters are not used for the generation of input values outside of this function. The input values are generated inside it, but rapid comes with its own generators, like gopter. For this example again a generator is used, that can generate a string based on a regex.
If the check does not succeed for a set of generated parameters, an error has to be thrown within the property function. Unlike the other two libraries, the result can not just be returned as a boolean value.
The property function is evaluated multiple times on the call `rapid.Check(t, property)`.

#### 2.1.2 Comparison generic functions

For generic functions we test a simple add function, that can add two numbers of type int or float. In this example we test for symmetry of this function.

This property can be tested in Haskell using quickcheck with this property:
```haskell
prop :: Number -> Number -> Bool
prop a b = a + b == b + a -- todo: verify if this works in Haskell
```

The same property can be tested in Go with the following implementations for the different libraries:

**testing/quick:**
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

**gopter:**
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

**rapid:**
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

### 2.2 Own implementation of quickcheck in Go (if existing libraries don't meet expectations)

> todo