---
theme: ./theme.json
title: Practicing Unsafe Stacks
author: Tighearnán Carroll
extensions:
  - terminal
  - file_loader
---
# Fáilte!

Practicing Unsafe Stacks.

<!-- stop -->

## Background

How using `unsafe` can reduce allocations (with hilarious consequences!)

<!-- stop -->

## Disclaimers

- I am by no means an expert in writing low allocation code.
<!-- stop -->

- I am by no means an expert in memory safety.
<!-- stop -->

- I am by no means recommending or advocating you do this.

---

# `unsafe`

Unsafe allows us to bypass certain golang memory and type safety operations by giving us a (limited) number of APIs to interface directly with the host machine's memory.

<!-- stop -->

```go

func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
func Sizeof(x ArbitraryType) uintptr
```

---

## `func Offsetof(x ArbitraryType) uintptr`

```file
path: src/examples/safestuff/cmd/offset_misaligned/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 6
    end: null
```

---

## `func Offsetof(x ArbitraryType) uintptr`

```file
path: src/examples/safestuff/cmd/offset_aligned/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 6
    end: null
```

---

## `func Sizeof(x ArbitraryType) uintptr`

```file
path: src/examples/safestuff/cmd/sizeof_misaligned/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 6
    end: null
```

---

## `func Sizeof(x ArbitraryType) uintptr`

```file
path: src/examples/safestuff/cmd/sizeof_aligned/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 6
    end: null
```

---

## `func Alignof(x ArbitraryType) uintptr`

Prints the alignment value of a given type.

Not really worth showing.

---

# Are these `unsafe`?

Probably not. On their own you can't do any damage with them.

Rob Pike has even opened a proposal to move these functions into another package in go2: https://github.com/golang/go/issues/5602.

<!-- stop -->

The lunacy begins when you use `unsafe.Pointer`.

<!-- stop -->


<!-- stop -->

hello

<!--
things to talk about

byte -> string conversion being BLAZING fast
string -> byte conversion being class but longer

[]custom -> []string conversion

compiler no escape hack

-->

this is a test
