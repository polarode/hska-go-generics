# Projektarbeit: Go Quickcheck und Generics

## Grundlagen

### Quickcheck in Haskell

todo: Grundlagen zu Quickcheck anhand von Haskell kurz zusammengefasst

### Generics in Go 1.18

todo: Aktuellen Stand der Generics in Go v1.18 (Funktionalität wurde teilweise nach 1.19 verschoben)

Generics in Version 1.18:
- Funktions- und Typdeklarationen erlauben jetzt Typparameter (Eckige Klammern)
- Interfaces als type constraints: Ein Interface kann jetzt eine Menge von Typen definieren
- Vordefinierter Identifier "any": alias für das leere Interface (matcht alle Typen)
- Vordefinierter Identifier "comparable": Menge aller Typen die mit == und != verglichen werden können

Fehlend in Version 1.18:
- Typdeklarationen innerhalb von generic functions
- Nur Methoden, die explizit im Interface des Typ-Parameters definiert sind, können genutzt werden
- Zugriff auf Felder in stuct nicht möglich

## Quickcheck in Go

### Stand der Technik

todo: Welche Bibliotheken gibt es bereits für Quickcheck? Wie unterscheiden sie sich? Welche Qualität/Welchen Nutzen haben diese? Unterstützen sie Generics?

- [testing/quick](https://pkg.go.dev/testing/quick)
  - relativ einfach
  - feature-frozen
  - shrinkage wird nicht unterstützt
- [github.com/leanovate/gopter](https://pkg.go.dev/github.com/leanovate/gopter)
  - komplexer
  - unterstützt shrinkage (benötigt aber wohl Code durch user)
  - letzter commit auf master: Mai 2021
- [pgregory.net/rapid](https://pkg.go.dev/pgregory.net/rapid)
  - unterstützt shrinkage
  - unterstützt "stateful" oder "model-based" testing
  - letzter commit auf master: September 2021

### Vergleich zweier ausgewählter Bibliotheken oder eigene Implementierung für Quickcheck (falls existieriende Bibliotheken Erwartungen nicht entsprechen)

todo