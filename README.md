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
	- [2.2 Summary](#22-summary)
		- [2.2.1 Generators](#221-generators)
		- [2.2.2 Properties](#222-properties)
		- [2.2.3 Conclusion](#223-conclusion)

## 1. Essentials

Normal unit tests check the code, they are written for, based on explicitly provided examples. Only the examples, that are provided in the form of unit tests by a developer, can be tested.

Property based testing, which originated in Haskell's library QuickCheck, uses another approach. Instead of explicit examples, generators are provided, that can generate valid input data. These generators allow us to automate the test inputs.
To also automate the validation of the tested inputs, properties are specified, that need to hold for all valid inputs.
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

But first the generator for the inputs of this function is needed. For QuickCheck this can be done by providing an instance for the class `Arbitrary` for the required type (in this case `String`). For some common types, these instances are already given, but if they are missing or for new types, they could be provided as follows:

```haskell
instance Arbitrary Char where 
   arbitrary = elements (['a'..'z'] ++ [' '])

instance Arbitrary String where
   arbitrary = listOf arbitrary
```
The functions `elements` and `listOf` used here are defined by the QuickCheck package. `elements` allows to generate a random single value from a list of possible values, while `listOf` allows to generate a list of random elements. The implementation can be looked up in the [QuickCheck Source](https://hackage.haskell.org/package/QuickCheck-2.14.2/docs/src/Test.QuickCheck.Gen.html):
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

With the generators provided by `arbitrary` we can define properties. The example function from above was `count`, which allows counting words in a string. One expectation to this function is, that for the reverse of a string `count` should yield the same number of words as for the string itself. This can be expressed with this property:

```haskell
-- | Reversing the string yields the same number of words.
prop :: String -> Bool
prop s = count s == count (reverse s)
```

The defined property can then be tested for example via the terminal with `quickCheck prop`. QuickCheck will check the property 100 times by generating random inputs for its functions according to the instances of Arbitrary, to see if the property holds.
If these tests pass, the result might look like this:
```
*Main> quickCheck prop
+++ OK, passed 100 tests.
```

Looking at the type definition of `quickCheck` tells, that the property, that we provide, has to be `Testable`.

```haskell
quickCheck :: Testable prop => prop -> IO ()
```

This is a class defined by QuickCheck and it has multiple instances. The one, that matches the property defined above is the following:

```haskell
instance (Arbitrary a, Show a, Testable prop) => Testable (a -> prop) where 
    ...
```

This recursive definition shows, that there has to be an instance for `Arbitrary` for the type of the parameter(s) of the property. A property with a parameter of some type without an instance of `Arbitrary` would not be allowed and noticed by the compiler. This means that property based testing in Haskell with QuickCheck is **type safe**.

If one test fails, quickCheck will use *shrinkage* to try and find simpler examples, for which the test also fails. This can make it easier for a human to debug the failure.

Such a failure can be found with this property:
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
- function and type declarations now allow type parameters (square bracket)
- interfaces can be used as type constraints: an interface can define a set of types
- predefined identifier "any": alias for the empty interface (matches all types)
- predefined identifier "comparable": set of all types, that can be compared with  == and !=

Missing in version 1.18:
- type-declarations within generic functions
- only methods, that are explicitly declared in the interface of the type parameter, can be used within the method
- accessing struct fields is not possible

For the full changelog, see the [release notes for v1.18](https://go.dev/doc/go1.18).

## 2. QuickCheck in Go

The concept of QuickCheck (or property based testing) has also been adopted to Go to some extent. This will be analyzed in the following section.

### 2.1 Existing work

There exist three popular libraries, that allow property based testing, similar to that of Quickcheck in Haskell.
- [testing/quick](https://pkg.go.dev/testing/quick)
- [github.com/leanovate/gopter](https://pkg.go.dev/github.com/leanovate/gopter)
- [pgregory.net/rapid](https://pkg.go.dev/pgregory.net/rapid)

The library *testing/quick* is an official package that comes with Go and is relatively easy to use. However, because it is kept rather simple, it is also limited in its functionality. For example, it does not support shrinkage.

*gopter* is described as a more sophisticated version of the *testing/quick* package. For example, it allows shrinkage, it allows better control of the generators and there is support for stateful tests.

*rapid* is similar to *gopter*. It also allows shrinkage, better generators and stateful tests. It claims to have a simpler API than *gopter* and it does not require user code to minimize failing tests.

The three libraries will be tested and compared on a few examples in the next sections.

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

The property that is being tested is again, that the reversed string and the original string have the same count of words.
As a quick reminder, this can be tested in Haskell using quickcheck with this property:

```haskell
-- | Reversing the string yields the same number of words.
prop :: String -> Bool
prop s = count s == count (reverse s)
```

The following function is used in the upcoming examples to reverse a string in Go:

```go
// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
```

The same property can be tested in Go with the following implementations for the different libraries:

**testing/quick:**
```go
func TestReverseSameNumberOfWords(t *testing.T) {
	property := func(x string) bool {
		y := Count(x)
		y2 := Count(Reverse(x))
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

The property that is supposed to be tested is provided as a function. The parameters of this function are used as input values for the functions that is being tested. Those values can be generated outside the property function. For that *testing/quick* needs an additional configuration, that contains the generators for the input values of the property function. Since *testing/quick* does not come with a generator for strings, this has to be implemented manually:

```go
func RandomStringGenerator(r *rand.Rand, size int, alphabet string) reflect.Value {
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := r.Intn(len(alphabet))
		buffer.WriteString(string(alphabet[index]))
	}
	return reflect.ValueOf(buffer.String())
}
```

This generator function is provided via the `Config` struct. This struct looks as follows:

```go
type Config struct {
	// MaxCount sets the maximum number of iterations.
	// If zero, MaxCountScale is used.
	MaxCount int
	// MaxCountScale is a non-negative scale factor applied to the
	// default maximum.
	// A count of zero implies the default, which is usually 100
	// but can be set by the -quickchecks flag.
	MaxCountScale float64
	// Rand specifies a source of random numbers.
	// If nil, a default pseudo-random source will be used.
	Rand *rand.Rand
	// Values specifies a function to generate a slice of
	// arbitrary reflect.Values that are congruent with the
	// arguments to the function being tested.
	// If nil, the top-level Value function is used to generate them.
	Values func([]reflect.Value, *rand.Rand)
}
```
It can also be used to configure the number of iterations and to set a source for random numbers. But to only set the generators, only Values function needs to be defined. For each parameter of the property function, the `values` array needs to have an entry with a `reflect.Value`.

```go
config := quick.Config{
	Values: func(values []reflect.Value, r *rand.Rand) {
		values[0] = RandomStringGenerator(r, 16, "abcxyz ")
	}
}
```
Reflective values have to be used, because inheritance is not supported in Go. With reflictive values, the types of these generators are only checked during runtime. But this also means, that compared to the Arbitrary instances in Haskell, the generators in Go are **not type safe**.

*testing/quick* also defines an interface for a `Generator` that can be used:
```go 
type Generator interface {
	// Generate returns a random instance of the type on which it is a
	// method using the size as a size hint.
	Generate(rand *rand.Rand, size int) reflect.Value
}
```

The property function is returning a boolean value, that contains the result for the tested parameters. If a set of parameters evaluates to false, *testing/quick* handles the error and prints according error messages. This function is evaluated multiple times for random inputs according to the given configuration on the call `quick.Check(property, &config)`. This function can return an error, if the check failed for a set of parameters. This error is then forwarded to the testing object `t` together with a description of the property.

**gopter:**
```go
func TestReverseSameNumberOfWords(t *testing.T) {
	properties := gopter.NewProperties(nil)
	property := func(x string) bool {
		y := Count(x)
		y2 := Count(Reverse(x))
		return y == y2
	}
	properties.Property("reverse string has the same number of words",
		prop.ForAll(property, gen.RegexMatch("[a-zA-Z ]+")))
	properties.TestingRun(t)
}
```

For *gopter* the properties that should be tested are again provided as functions. The parameters of this function again serve as input of the function being tested and are generated outside using generators. *gopter* comes with some basic generators, that can be used here. In this example a generator is used, that can generate strings based on a regex.
*gopter* defines its own type `Gen` for generators. This type is basically a function, that takes some `GenParameters` and provides a `GenResult`. 

```go
type Gen func(*GenParameters) *GenResult
```

By looking at those types again, we can see some similarities to *testing/quick*:

```go
type GenParameters struct {
	MinSize        int
	MaxSize        int
	MaxShrinkCount int
	Rng            *rand.Rand
}
```

The GenParameters contain a source of random numbers, just like Config for *testing/quick*. Additionally, it can hold values for minimum and maximum size, which can be used for strings or slices and a value for the maximum number of shrinks.

```go
type GenResult struct {
	Labels     []string
	Shrinker   Shrinker
	ResultType reflect.Type
	Result     interface{}
	Sieve      func(interface{}) bool
}
```

*GenResult* contains the result itself (with the empty interface as type to match all possible types), and the reflective result type. This is again similar to the reflective value, that is used as result type of the generator functions in *testing/quick*. Additionally *GenResult* contains a sieve function that can check if a value is valid, a shrinker and labels.

In contrary to *testing/quick*, multiple properties can be provided at once and are tested in a testing run. For each property, generators have to be provided for each input parameter of the function. 

The properties function returns a boolean value, that contains the result for the tested parameters. If a set of parameters evaluates to false, *gopter* handles the error and prints according error messages based on the provided name (first parameter of `properties.Property`). 
On the call `properties.TestingRun(t)` all the properties provided within `properties` are evaluated multiple times with random inputs.

**rapid:**
```go
func TestReverseSameNumberOfWords(t *testing.T) {
	property := func(t *rapid.T) {
		x := rapid.StringMatching("[a-zA-Z ]+").Draw(t, "words").(string)
		y := Count(x)
		y2 := Count(Reverse(x))
		if y != y2 {
			t.Error("falsified: reverse string has the same number of words")
		}
	}
	rapid.Check(t, property)
}
```

The property to be tested is again provided as a function. However, the parameters are not used for the generation of input values outside this function. The input values are generated inside it, but *rapid* comes with its own generators, like *gopter*. For this example again a generator is used, that can generate a string based on a regex.

The `Generator` for *rapid* is defined as a struct as follows:
```go
type Generator struct {
	impl    generatorImpl
	typ     reflect.Type
	strOnce sync.Once
	str     string
}
```
The internal details are a bit more complicated here, but from this it is already visible, that the reflect package is used. And this is the most important fact needed for the comparison with the QuickCheck implementation in Haskell later.

If the check does not succeed for a set of generated parameters, an error has to be thrown within the property function. Unlike the other two libraries, the result can not just be returned as a boolean value.
The property function is evaluated multiple times on the call `rapid.Check(t, property)`.

#### 2.1.2 Comparison generic functions

For generic functions we test a simple add function, that can add two numbers of type int or float. This function is defined as follows:

```go
func Add[t Number](a, b t) t {
	return a + b
}
```

The definition uses an interface `Number` that is used as type constraint for the generic type t. It is defined as the union of the types `int64` and `float64`:

```go
type Number interface {
	int64 | float64
}
```

In this example we test for symmetry of this function.
In Haskell this property can be tested using quickcheck with this property:

```haskell
prop :: (Eq a, Num a) => a -> a -> Bool
prop a b = a + b == b + a
```

The same property can be tested in Go with the following implementations for the different libraries:

**testing/quick:**
```go
func TestAddSymmetric(t *testing.T) {
	propertyInt := func(a, b int64) bool {
		y := Add(a, b)
		y2 := Add(b, a)
		return y == y2
	}
	propertyFloat := func(a, b float64) bool {
		y := Add(a, b)
		y2 := Add(b, a)
		return y == y2
	}
	if err := quick.Check(propertyInt, nil); err != nil {
		t.Error("falsified: add is symmetric (int64)", err)
	}
	if err := quick.Check(propertyFloat, nil); err != nil {
		t.Error("falsified: add is symmetric (float64)", err)
	}
}
```

As can be seen from this example, the property can not be written in a generic way. This means it has to be written for all possible type of type constraint and the Check call needs to be executed for all of them.
In contrary to the string type from the non-generic example, no generator has to be provided for primitive types like the ones used here, therefore no configuration needs to be specified. The correct generators are used implicitly.

**gopter:**
```go
func TestAddSymmetric(t *testing.T) {
	properties := gopter.NewProperties(nil)
	propertyInt := func(a, b int64) bool {
		y := Add(a, b)
		y2 := Add(b, a)
		return y == y2
	}
	propertyFloat := func(a, b float64) bool {
		y := Add(a, b)
		y2 := Add(b, a)
		return y == y2
	}
	properties.Property("add is symmetric (int64)",
		prop.ForAll(propertyInt, gen.Int64(), gen.Int64()))
	properties.Property("add is symmetric (float64)",
		prop.ForAll(propertyFloat, gen.Float64(), gen.Float64()))
	properties.TestingRun(t)
}
```

For *gopter*, the same problem can be observed: The property can not be defined in a generic way. But in this example we can now use the properties object and add multiple properties to it. At the end of the test, we execute the check on all the properties.

**rapid:**
```go
func TestAddSymmetric(t *testing.T) {
	propertyInt := func(t *rapid.T) {
		a := rapid.Int64().Draw(t, "a").(int64)
		b := rapid.Int64().Draw(t, "b").(int64)
		y := Add(a, b)
		y2 := Add(b, a)
		if y != y2 {
			t.Error("falsified: add is symmetric (int64)")
		}
	}
	propertyFloat := func(t *rapid.T) {
		a := rapid.Float64().Draw(t, "a").(float64)
		b := rapid.Float64().Draw(t, "b").(float64)
		y := Add(a, b)
		y2 := Add(b, a)
		if y != y2 {
			t.Error("falsified: add is symmetric (float64)")
		}
	}
	rapid.Check(t, propertyInt)
	rapid.Check(t, propertyFloat)
}
```

*rapid* again has the same issue as the previous two examples: The property has to be written multiple times, because it can not be written in a generic way.

### 2.2 Summary

#### 2.2.1 Generators
The three implementations are using reflective values for their generators. This means they are not type safe, because the types are only checked during runtime. QuickCheck in Haskell on the other hand is type safe, because it is implemented using overloading. All types, for which an instance for the Arbitrary class is provided, can be generated. This means, type safety can be checked already during compile time.

To show that, consider this alternative generator, that generates random integer values instead of strings:

```go
func RandomIntGenerator(r *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(rand.Intn(size))
}
```

Since its return value has the same type `reflect.Value`, it can replace the string generator from the example of *testing/quick*:

```go
func TestReverseSameNumberOfWords(t *testing.T) {
	property := func(x string) bool {
		y := testable.Count(x)
		y2 := testable.Count(stringutil.Reverse(x))
		return y == y2
	}
	config := quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = RandomIntGenerator(r, 16)
		}}
	if err := quick.Check(property, &config); err != nil {
		t.Error("falsified: reverse string has the same number of words", err)
	}
}
```

The Go compiler doesn't recognize, that this doesn't work and so we only get an error message during runtime when executing the test:

```text
=== RUN   TestReverseSameNumberOfWords
--- FAIL: TestReverseSameNumberOfWords (0.00s)
panic: reflect: Call using int as type string [recovered]
        panic: reflect: Call using int as type string
```


There are also differences between the three libraries regarding generators. *testing/quick* comes only with basic generators, like for primitive types. *gopter* and *rapid* both offer more options out of the box, like the string generator based on a regex used in the earlier examples. All the available, predefined generators can be found in the Go documentation ([*gopter*](https://pkg.go.dev/github.com/leanovate/gopter/gen?utm_source=godoc), [*rapid*](https://pkg.go.dev/pgregory.net/rapid#pkg-index)). In all cases own, more specific generators can be defined.

#### 2.2.2 Properties

All three libraries are using functions to define a property, but there are still differences. For *testing/quick* and *gopter* a property function has to return a boolean value, that says, if this property was checked successfully or not. This is the same way as it is implemented for QuickCheck in Haskell. For *rapid*, property functions are provided without a return value. If the property does not apply, an error has to be thrown.
Property functions in Go show a current limitation when used in combination with generics. It is not possible to provide a generic function to a property based testing library, so that the library can check it for the different types on its own. This is possible with QuickCheck in Haskell. The "next best thing" possible in Go is to define the property as a generic function and then check it multiple times with the specific type parameters. For example with *gopter*:
```go
func property_AddSymmetry[t Number](a, b t) bool {
	y := Add(a, b)
	y2 := Add(b, a)
	return y == y2
}
```
```go
properties.Property("add is symmetric (int64)",
	prop.ForAll(property_AddSymmetry[int64], gen.Int64(), gen.Int64()))
properties.Property("add is symmetric (float64)",
	prop.ForAll(property_AddSymmetry[float64], gen.Float64(), gen.Float64()))
```
In Haskell however, this is as simple as the following code example:
```haskell
prop :: (Eq a, Num a) => a -> a -> Bool
prop a b = a + b == b + a
```

#### 2.2.3 Conclusion

The three packages that are available in Go implement property based testing as good as possible, but there are some limitations coming with the language. For example it is not possible to provide generators for properties in a way that is type safe, because it requires the use of reflective values. Furthermore it is not possible to cover a generic function with a single property/check.
Two of the packages (gopter and rapid) in Go also come with support for shinkage like it is possible with QuickCheck in Haskell. However, those two are not official packages and might not be updated by their maintainers in the future.
But there is also a small thing, that gopter and rapid provide, that is not available in Haskell: They provide some more complex generators for common types out of the box, like the string generator, that matches a regex.

Finally, here is a table quickly showing the most obvious facts and differences between QuickCheck in Haskell and the three packages in Go:

<table>
	<tr>
		<th rowspan=2></th>
		<th>Haskell</th>
		<th colspan=3>Go</th>
	</tr>
	<tr>
		<th>QuickCheck</th>
		<th>testing/quick</th>
		<th>gopter</th>
		<th>rapid</th>
	</tr>
	<tr>
		<td>generators type safe</td>
		<td>✔</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>
	<tr>
		<td>comes with complex generators</td>
		<td></td>
		<td></td>
		<td>✔</td>
		<td>✔</td>
	</tr>
	<tr>
		<td>shinkage</td>
		<td>✔</td>
		<td></td>
		<td>✔</td>
		<td>✔</td>
	</tr>
	<tr>
		<td>simplicity to define tests</td>
		<td>✔</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>
	<tr>
		<td>support for generics</td>
		<td>✔</td>
		<td></td>
		<td></td>
		<td></td>
	</tr>
	<tr>
		<td>last version release*</td>
		<td>14.11.2020</td>
		<td>01.06.2022</td>
		<td>09.11.2020</td>
		<td>05.07.2021</td>
	</tr>
</table>
(*) at the time of this work