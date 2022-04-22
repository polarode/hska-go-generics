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

### Eigene Implementierung für Quickcheck (falls existieriende Bibliotheken Erwartungen nicht entsprechen)

> todo