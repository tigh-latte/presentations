---
theme: ./theme.json
title: Practicing Unsafe Stacks
author: Tighearnán Carroll
styles:
  table:
      column_spacing: 3
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

- I am by no means advocating anything in this presentation.

---

# What is `unsafe`?

A pointer in golang has an associated type and points to some address in memory.

Using unsafe, we can point some address in memory, but without an associated type. We can then cast a type to this pointer.

<!-- stop -->

## Quick warning.

You don't need to know any of what I'm about to tell your in order to write good go.

---

# What is `unsafe`?

As of go1.20, `unsafe` exposes the following:

```go

func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr
func Sizeof(x ArbitraryType) uintptr

func String(ptr *byte, len IntegerType) string
func StringData(str string) *byte
func Slice(ptr *ArbitraryType, len IntegerType) []ArbitraryType
func SliceData(slice []ArbitraryType) *ArbitraryType

type Pointer *ArbitraryType
func Add(ptr Pointer, len IntegerType) Pointer
```

Let's address the first three.

---

# What is `unsafe`?

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

# What is `unsafe`?

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

# What is `unsafe`?

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

# What is `unsafe`?

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

# What is `unsafe`?

## `func Alignof(x ArbitraryType) uintptr`

Prints the alignment value of a given type.

Not really worth showing.

---

# What about these functions are `unsafe`?

Not much on their own.

Rob Pike has even opened a proposal to move these functions into another package in go2: https://github.com/golang/go/issues/5602.

<!-- stop -->

The real lunacy begins when you use `unsafe.Pointer`.

---

# What is `unsafe.Pointer`?

As seen before, the standard library defines this as:

```go
type Pointer *ArbitraryType
```

This has a bit of an unusal look about it, just means that we are pointing to a memory address without any type assumptions.

To understand this, we will first look at a typed pointer.


---

# What is `unsafe.Pointer`?

## Typed pointers


If we have an `*int` that points to `0xc000020158`, golang isn't __just__ concerned with that memory address. It is instead concerned 8 memory addresses, starting with `0xc000020158`.

<!-- stop -->


```go
    i := 34562
    _ = &i // Mem addr: 0xc000020158
```


```
 &i points to                                                             golang reads until
     |                                                                            |
     v                                                                            v
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000010 | 10000111 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# What is `unsafe.Pointer`?

## Typed pointers

If we have a `*string` that points to `0xc000020158`, golang isn't __just__ concerned with that memory address. It is instead concerned 16 memory addresses, starting with `0xc000020158`.

```go
    s := strings.Repeat("holy hell", 42)
    _ = &s // Mem addr: 0xc000020158
```


```
 &s points to
     |
     v
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 10000000 | 00001011 | 0        | 11000000 | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
+---------------------------------------------------------------------------------------+
| 0xc..160 | 0xd..161 | 0xc..162 | 0xc..163 | 0xc..164 | 0xc..165 | 0xc..166 | 0xc..167 |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 01111010 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
                                                                                   ^
                                                                                   |
                                                                           golang reads until
```

---

# What is `unsafe.Pointer`?

## Typed pointers

So if we convert a typed pointer to an `unsafe.Pointer`, we are pretty much just pointing to that single memory address.

```go
    i := 34562
    p := unsafe.Pointer(&i) // Mem addr: 0xc000020158
```

---

# What is `unsafe.Pointer`?

## Typed pointers

So if we convert a typed pointer to an `unsafe.Pointer`, we are pretty much just pointing to that single memory address.

```go
 >  i := 34562
    p := unsafe.Pointer(&i) // Mem addr: 0xc000020158
```

```
 &i points to                                                               i is read until
     |                                                                            |
     v                                                                            v
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000010 | 10000111 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# What is `unsafe.Pointer`?

## Typed pointers

So if we convert a typed pointer to an `unsafe.Pointer`, we are pretty much just pointing to that single memory address.

```go
    i := 34562
 >  p := unsafe.Pointer(&i) // Mem addr: 0xc000020158
```

```
 p points to
     |
     v
+----------+
| 0xc..158 |
+----------+
| 00000010 |
+----------+
```

<!-- stop -->

This will remain the case until we read the `unsafe.Pointer` as a type, and herein lies the madness.

---

# What is `unsafe.Pointer`?

With this `unsafe.Pointer` we can have the binary data stored at `*i` be read as another type.

To dereference an `unsafe.Pointer`, we first cast it to a type, and then the data is read as it would be for this cast type.

<!-- stop -->

Continuing on from the previous example, using `unsafe.Pointer` as a medium, we can read `i` as a `string`:

```go
    i := 34562
    str := *(*string)(unsafe.Pointer(&i))
```

```
 &i points to                                                               i is read until
 &str points to                                                                   |
     |                                                                            |
     v                                                                            v
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000010 | 10000111 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
+---------------------------------------------------------------------------------------+
| 0xc..160 | 0xd..161 | 0xc..162 | 0xc..163 | 0xc..164 | 0xc..165 | 0xc..166 | 0xc..167 |
+----------+----------+----------+----------+----------+----------+----------+----------+
|   ????   |   ????   |   ????   |   ????   |   ????   |   ????   |   ????   |   ????   |
+----------+----------+----------+----------+----------+----------+----------+----------+
                                                                                  ^
                                                                                  |
                                                                           str is read until
```

---

# What is `unsafe.Pointer`?

## Wtaf is `str := *(*string)(unsafe.Pointer(&i))`

Let's break it down.

---

# What is `unsafe.Pointer`?

## Wtaf is `str := *(*string)(unsafe.Pointer(&i))`

We have `i`.
```go
    i
```

---

# What is `unsafe.Pointer`?

## Wtaf is `str := *(*string)(unsafe.Pointer(&i))`

We get the pointer to `i`.
```go
    &i
```

---

# What is `unsafe.Pointer`?

## Wtaf is `str := *(*string)(unsafe.Pointer(&i))`

We cast this pointer to an `unsafe.Pointer`:
```go
    unsafe.Pointer(&i)
```

---

# What is `unsafe.Pointer`?

## Wtaf is `str := *(*string)(unsafe.Pointer(&i))`

We cast this `unsafe.Pointer` to a `*string`:
```go
    (*string)(unsafe.Pointer(&i))
```

---

# What is `unsafe.Pointer`?

## Wtaf is `str := *(*string)(unsafe.Pointer(&i))`

We dereference this cast:
```go
    *(*string)(unsafe.Pointer(&i))
```

---

# What is `unsafe.Pointer`?

## Why would we read an `*int` to a `*string`?

Simple answer, we wouldn't. Don't ever do this.

The uses of cross-type casting to improve performance applies to other data structures, which we will get into later.

However though, it can be interesting to mess around with `unsafe.Pointer` hacks to see how golang reads data.

---

# Some fun with `unsafe.Pointer`


```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g;15,34d;'
lines:
    start: 6
    end: null
```

<!-- stop -->

```
 i points to                                                                 i reads until
 b points to
 f points to                                                                 f reads until
 a points to                                                                 a reads until
     |                                                                             |
     V                                                                             V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 0        | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
     ^
     |
 b reads until
```


---

# Some fun with `unsafe.Pointer`

```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g;13,14d;20,34d;'
lines:
    start: 6
    end: null
```

<!-- stop -->

```
 i points to                                                                 i reads until
 b points to
 f points to                                                                 f reads until
 a points to                                                                 a reads until
     |                                                                             |
     V                                                                             V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00001010 | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
     ^
     |
 b reads until
```

---

# Some fun with `unsafe.Pointer`

```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g;13,19d;25,34d;'
lines:
    start: 6
    end: null
```

<!-- stop -->

```
 i points to                                                                 i reads until
 b points to
 f points to                                                                 f reads until
 a points to                                                                 a reads until
     |                                                                             |
     V                                                                             V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
     ^
     |
 b reads until
```

---

# Some fun with `unsafe.Pointer`

```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g;13,24d;30,34d;'
lines:
    start: 6
    end: null
```

<!-- stop -->

```
 i points to                                                                 i reads until
 b points to
 f points to                                                                 f reads until
 a points to                                                                 a reads until
     |                                                                             |
     V                                                                             V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 11110110 | 11111111 | 11111111 | 11111111 | 11111111 | 11111111 | 11111111 | 11111111 |
+----------+----------+----------+----------+----------+----------+----------+----------+
     ^
     |
 b reads to
```

---

# Then we have `uintptr`

This is an integer representation of a memory address.

An `unsafe.Pointer` can be cast to a `uintptr`, letting us do pointer arithmetic.

```go
    i := new(int)  // 0xc000020158
    *i = 256
    ptr := uintptr(unsafe.Pointer(i)) // 824633852248
    ptr++ // 824633852249
    b := *(*byte)(unsafe.Pointer(ptr)) // 0xc000020159 (second byte of i, 1)
```

---

# Then we have `uintptr`

Let's explain that real quick.

<!-- stop -->

```go
 >  i := new(int)  // 0xc000020158
    *i = 256
    ptr := uintptr(unsafe.Pointer(i)) // 824633852248
    ptr++ // 824633852249
    b := *(*byte)(unsafe.Pointer(ptr)) // 0xc000020159 (second byte of i, 1)
```

```
   
 i points to
     |
     V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 0        | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# Then we have `uintptr`

We write `256` to `*i`

```go
    i := new(int)  // 0xc000020158
 >  *i = 256
    ptr := uintptr(unsafe.Pointer(i)) // 824633852248
    ptr++ // 824633852249
    b := *(*byte)(unsafe.Pointer(ptr)) // 0xc000020159 (second byte of i, 1)
```

```
   
 i points to
     |
     V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# Then we have `uintptr`

We get the `uintptr` of `i`

```go
    i := new(int)  // 0xc000020158
    *i = 256
 >  ptr := uintptr(unsafe.Pointer(i)) // 824633852248
    ptr++ // 824633852249
    b := *(*byte)(unsafe.Pointer(ptr)) // 0xc000020159 (second byte of i, 1)
```

```
 ptr points to
 i points to
     |
     V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# Then we have `uintptr`

We increment the `uintptr`.

```go
    i := new(int)  // 0xc000020158
    *i = 256
    ptr := uintptr(unsafe.Pointer(i)) // 824633852248
 >  ptr++ // 824633852249
    b := *(*byte)(unsafe.Pointer(ptr)) // 0xc000020159 (second byte of i, 1)
```

```
          ptr points to
 i points to    |
     |          |
     V          V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# Then we have `uintptr`

We read that point as a `byte`.

```go
    i := new(int)  // 0xc000020158
    *i = 256
    ptr := uintptr(unsafe.Pointer(i)) // 824633852248
    ptr++ // 824633852249
 >  b := *(*byte)(unsafe.Pointer(ptr)) // 0xc000020159 (second byte of i, 1)
```

```
          ptr points to
 i points to    |
     |          |
     V          V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

<!-- stop -->

__NEVER__ user this.

---

# Don't tell me what to do!

<!-- stop -->

Pointer arithmetic + a garbage collector == bad time.

Take the following:

```file
path: src/examples/ptr_math/cmd/gc_danger/main.go
lang: go
transform: sed 's/\t/    /g;12,13d;27,29d'
lines:
    start: 10
    end: null
```

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 10
init_text: (cd src/examples/ptr_math/; go run cmd/gc_danger/main.go)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# We get it this is dangerous! How can this boost performance?

`unsafe` can boost performance in two main ways:

<!-- stop -->

1. Avoiding allocations, meaning less garbage needs collected
<!-- stop -->

2. Avoiding data copying, meaning less cpu.
<!-- stop -->

 a. In fact, this can turn a few O(N) operations into O(1) (And that O is a very small one)

---

# Casting slices

While this is less relevant in a world with generics, quite a bit of the std lib and third party libs don't use generics.

<!-- stop -->

Say we have the following type:

```go
type MyCoolType string
```

We want to pass a slice of this type, `[]MyCoolType`, to `strings.Join`. `strings.Join` takes a `[]string`.

What can we do?

```go
type MyCoolType string

func MyCoolJoiner(mm []MyCoolType) MyCoolType {
    // What goes here?
}
```

---

# Casting slices

Sadly, we can't just do this:

```go
type MyCoolType string

func MyCoolJoiner(mm []MyCoolType, div string) MyCoolType {
    return MyCoolType(strings.Join(mm, div))
}
```
```

err: cannot use mm (variable of type []MyCoolType) as []string value in argument to strings.Join
```

<!-- stop -->

Or this:

```go
type MyCoolType string

func MyCoolJoiner(mm []MyCoolType, div string) MyCoolType {
    return MyCoolType(strings.Join([]string(mm), div))
}
```
```
err: cannot convert mm (variable of type []MyCoolType) to type []string
```

<!-- stop -->

We can't even use generics to fool the type system:

```go
type MyCoolType string

func MyCoolJoiner[T ~string](tt []T, div string) T {
    return T(strings.Join(tt, div))
}
```

```
err: cannot use tt (variable of type []T) as []string value in argument to strings.Join
```

---

# Casting slices

The only solution that we really have is to build a whole new slice of `[]string` and transfer the data over:
```go
type MyCoolType string

func MyCoolJoiner(mm []MyCoolType, div string) MyCoolType {
    ss := make([]string, len(mm))
    for i, m := range mm {
        ss[i] = string(m)
    }
    return MyCoolType(strings.Join(ss, div))
}
```

<!-- stop -->

We've all probably had to do something similar before, and we've all felt worried over how inefficient this could be when writing it.

---

# Casting slices

## Benchmarking the safe approach

<!-- stop -->

Take the following benchmark tests:
```file
path: src/examples/casting/cmd/customslice_to_slice/bench/bench_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 8
    end: 36
```

---

# Casting slices

## Benchmarking the safe approach

```terminal-ex
command: zsh -il
rows: 20
init_text: (cd src/examples/casting/; go test -run=x -bench=SafeSlice ./cmd/customslice_to_slice/bench/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Casting slices

The unsafe solution:

```go
type MyCoolType string

func MyCoolJoiner(mm []MyCoolType, div string) MyCoolType {
    ss := *(*[]string)(unsafe.Pointer(&mm))
    return MyCoolType(strings.Join(ss, div))
}
```

<!-- stop -->

How does this perform in comparison?

---

# Casting slices

## Benchmarking the unsafe approach

<!-- stop -->

Take the following benchmark tests:
```file
path: src/examples/casting/cmd/customslice_to_slice/bench/bench_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 37
    end: null
```

---

# Casting slices

## Benchmarking the unsafe approach

```terminal-ex
command: zsh -il
rows: 20
init_text: (cd src/examples/casting/; go test -run=x -bench=UnsafeSlice ./cmd/customslice_to_slice/bench/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Casting slices

## All together now!

```terminal-ex
command: zsh -il
rows: 28
init_text: (cd src/examples/casting/; go test -run=x -bench=Slice ./cmd/customslice_to_slice/bench/... -benchmem -benchtime=2s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# String headers

## Holy hell, but why?

To understand this, we need to understand a few things:

<!-- stop -->

1. String headers.
<!-- stop -->

1. How golang converts to and from string types.
<!-- stop -->

1. Slice headers.


---

# Casting slices

## String headers

So what happens when we convert a string to a custom type?

```go
    s := "hello"
    m := MyCoolType(s)
```

<!-- stop -->

Remember when we mentioned that a `*string` cares about 16 memory addresses? These 16 bytes make up a the string header.

This means a pointer to a string does not point to the actual characters of that string.

<!-- stop -->

## String header in memory

```go
    s := strings.Repeat("holy hell", 42)
    _ = &s
```

```
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f | <- Data pointer
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 10000000 | 00001011 | 00000000 | 11000000 | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
+---------------------------------------------------------------------------------------+
| 0xc..160 | 0xd..161 | 0xc..162 | 0xc..163 | 0xc..164 | 0xc..165 | 0xc..166 | 0xc..167 | <- Length
+----------+----------+----------+----------+----------+----------+----------+----------+
| 01111010 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# Casting slices

## String headers

A struct representation of a string header exists within the `reflect` package:

```go
type StringHeader struct {
    Data uintptr // first 8 bytes
    Len  int     // second 8 bytes
}
```

The indirection of a string's content via `StringHeader.Data` is how golang guarantees adding a `string` to a struct will only increase that structs size by 16 bytes, despite a `string` being variable in length.

<!-- stop -->

```go
    myString := "holy hell"
```

<!-- stop -->

```go

    +----------------+                       +----------+----------+----------+
    | myString       | --------------------> | 0xc..158 | 0xc....  | 0xc..160 |
    | Data: 0xc..158 |                       +----------+----------+----------+
    | Len: 9         |                       | 'h'      | . . .    | 'l'      |
    +----------------+                       +----------+----------+----------+
```

<!-- stop -->

If we clone `myString` via `strings.Clone`, the underlying character data would be cloned in memory, and a new `StringHeader` pointing to this cloned data is be created:

```go
    myString := "holy hell"
    myClone := strings.Clone(myString)
```

<!-- stop -->

```go
    +----------------+                       +----------+----------+----------+
    | myString       | --------------------> | 0xc..158 | 0xc....  | 0xc..160 |
    | Data: 0xc..158 |                       +----------+----------+----------+
    | Len: 9         |                       | 'h'      | . . .    | 'l'      |
    +----------------+                       +----------+----------+----------+

    +----------------+                       +----------+----------+----------+
    | myClone        | --------------------> | 0xc..170 | 0xc....  | 0xc..178 |
    | Data: 0xc..170 |                       +----------+----------+----------+
    | Len: 9         |                       | 'h'      | . . .    | 'l'      |
    +----------------+                       +----------+----------+----------+
```

---

# Casting slices

## String headers

However, if we cast this `myString` to a custom type the underlying data isn't cloned. Instead a new `StringHeader` is created pointing to the same character array in memory:

```go
    myString := "holy hell"
    myCoolType := MyCoolType(myString)
```

<!-- stop -->

```go
    +----------------+                       +----------+----------+----------+
    | myString       | ----------+---------> | 0xc..158 | 0xc....  | 0xc..160 |
    | Data: 0xc..158 |           |           +----------+----------+----------+
    | Len: 9         |           |           | 'h'      | . . .    | 'l'      |
    +----------------+           |           +----------+----------+----------+
                                 |
    +----------------+           |
    | myCoolType     | ----------+
    | Data: 0xc..158 |
    | Len: 9         |
    +----------------+
```

---

# Casting slices

## String headers

So for every iteration in our loop, we're cloning 16 bytes of data:

```go
func MyCoolJoiner(mm []MyCoolType, div string) MyCoolType {
    ss := make([]string, len(mm))
    for i, m := range mm {
        ss[i] = string(m)
    }
    return MyCoolType(strings.Join(ss, div))
}
```

This combination of the allocation from `make([]string, len(mm))`, interating every item, and making a copy of the string header for each item, all adds up to a degrading O(N) execution time.

---

# Casting slices

## Why unsafe is faster

Unsafe is faster as it takes advantage of how there are no in memory differences between a `[]<type>` and a `[]<custom type>`.

<!-- stop -->

Let's look at the unsafe solution again:

```go
    ss := *(*[]string)(unsafe.Pointer(&mm))
```

This is casting an entire slice from one type to another. How is this possible?

Because it isn't casting the __slice__, it's casting the __slice header__.

---

# Casting slices

## Slice headers

A slice header is like a `StringHeader`, only they are 24 bytes long. The extra 8 bytes determine the slice's capacity.

```go
    bb := make([]string, 5, 10)
```


```go
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f | <- Data pointer
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 10000000 | 00001011 | 00000000 | 11000000 | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
+---------------------------------------------------------------------------------------+
| 0xc..160 | 0xd..161 | 0xc..162 | 0xc..163 | 0xc..164 | 0xc..165 | 0xc..166 | 0xc..167 | <- Length
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000101 | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+

+---------------------------------------------------------------------------------------+
| 0xc..168 | 0xd..169 | 0xc..16a | 0xc..16b | 0xc..16c | 0xc..16d | 0xc..16e | 0xc..16d | <- Capacity
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00001010 | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
```

---

# Casting slice

## Slice headers

And likewise, a struct representation of a slice header exists in `reflect`:

```go
type SliceHeader struct {
    Data uintptr // bytes 0-7
    Len  int     // bytes 8-15
    Cap  int     // bytes 16-23
}
```


<!-- stop -->

So with the following:
```go
    mm := []MyCoolType{"hello", "there"}
    ss := *(*[]string)(unsafe.Pointer(&m))
```

When we cast `mm` to a `[]string` using unsafe, we do 0 duplication of data. Instead, we create `ss` and have it point to __the exact same header__ as `mm`.

<!-- stop -->

```go
    mm points to
      |   ss points to
      |       |
      V       V
    +----------------+                       +--------------+------------+--------------+
    | Data: 0xc..158 | --------------------> | 0xc..158     | 0xc...     | 0xc..68      |
    | Len: 2         |                       +--------------+------------+--------------+
    | Cap: 2         |                       | StringHeader | . . .      | StringHeader |
    +----------------+                       +--------------+------------+--------------+
```

---

# Any questions so far?

---

# Casting `[]byte` to a `string`

Golang, by design, doesn't have many built in operations that are `O(N)`, as this can lead to writing expensive code without even realising.

Casting `[]byte` to `string` (and visa-versa) is one of these few exceptions.

Given that a `string` in golang is immutable but a `[]byte` isn't, in order to ensure integrity, `string([]byte)` makes a full clone of the `[]byte`.

---

# Casting `[]byte` to a `string`

Let's explain the following:

```go
    bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
    s := string(bb)
    bb[0] = 'z'
    fmt.Println(s) // Output: Hello world
```

---

# Casting `[]byte` to a `string`

We initalise a `[]byte` containing some data.


```go
 >  bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
    s := string(bb)
    bb[0] = 'z'
    fmt.Println(s) // Output: Hello world
```

```
    bb points to
         |
         V
    +----------------+                       +----------+----------+----------+
    | Data: 0xc..158 | --------------------> | 0xc..158 | 0xc....  | 0xc..163 |
    | Len: 11        |                       +----------+----------+----------+
    | Cap: 11        |                       | 01001000 | . . .    | 01100100 |
    +----------------+                       +----------+----------+----------+

```

---

# Casting `[]byte` to a `string`

We clone `bb` and reference the cloned data as a `string`.


```go
    bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
 >  s := string(bb)
    bb[0] = 'z'
    fmt.Println(s) // Output: Hello world
```

```
         b
         |
         V
    +----------------+                       +----------+----------+----------+
    | Data: 0xc..158 | --------------------> | 0xc..158 | 0xc....  | 0xc..164 |
    | Len: 11        |                       +----------+----------+----------+
    | Cap: 11        |                       | 01001000 | . . .    | 01100010 |
    +----------------+                       +----------+----------+----------+

         s
         |
         V
    +----------------+                       +----------+----------+----------+
    | Data: 0xc..178 | --------------------> | 0xc..178 | 0xc....  | 0xc..184 |
    | Len: 11        |                       +----------+----------+----------+
    |                |                       | 01001000 | . . .    | 01100010 |
    +----------------+                       +----------+----------+----------+

```

---

# Casting `[]byte` to a `string`

We change the first character of `[]byte` to whatever the byte value of `'z'` is, and `s` is untouched.


```go
    bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
    s := string(bb)
 >  bb[0] = 'z'
 >  fmt.Println(s) // Output: Hello world
```

```
         b
         |
         V
    +----------------+                       +----------+----------+----------+
    | Data: 0xc..158 | --------------------> | 0xc..158 | 0xc....  | 0xc..164 |
    | Len: 11        |                       +----------+----------+----------+
    | Cap: 11        |                       | 01111010 | . . .    | 01100010 |
    +----------------+                       +----------+----------+----------+

         s
         |
         V
    +----------------+                       +----------+----------+----------+
    | Data: 0xc..178 | --------------------> | 0xc..178 | 0xc....  | 0xc..184 |
    | Len: 11        |                       +----------+----------+----------+
    |                |                       | 01001000 | . . .    | 01100010 |
    +----------------+                       +----------+----------+----------+

```

<!-- stop -->

I'm sure yous can guess, this clone isn't cheap.

---

# Casting `[]byte` to a `string`

## Benchmarking the safe approach

<!-- stop -->

```file
path: src/examples/casting/cmd/byteslice_to_string/bench/bench_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 10
    end: 29
```

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 20
init_text: (cd src/examples/casting/; go test -run=x -bench=SafeCast ./cmd/byteslice_to_string/bench/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Casting `[]byte` to a `string`

## The unsafe way

Truth be told there are three ways to do this.

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 16
    end: 22
```

<!-- stop -->

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 23
    end: 28
```

<!-- stop -->

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 29
    end: null
```

---

# Casting `[]byte` to a `string`

## Benchmarking the unsafe approaches

```file
path: src/examples/casting/cmd/byteslice_to_string/bench/bench_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 30
    end: 47
```

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 20
init_text: (cd src/examples/casting/; go test -run=x -bench=UnsafeCast ./cmd/byteslice_to_string/bench/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Casting `[]byte` to a `string`

## Benchmarking the unsafe approaches

```file
path: src/examples/casting/cmd/byteslice_to_string/bench/bench_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 48
    end: 65
```

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 20
init_text: (cd src/examples/casting/; go test -run=x -bench=UnsafeString ./cmd/byteslice_to_string/bench/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Casting `[]byte` to a `string`

## Benchmarking the unsafe approaches

```file
path: src/examples/casting/cmd/byteslice_to_string/bench/bench_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 66
    end: null
```

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 20
init_text: (cd src/examples/casting/; go test -run=x -bench=UnsafeHeader ./cmd/byteslice_to_string/bench/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Casting `[]byte` to a `string`

## All together now!

```terminal-ex
command: zsh -il
rows: 40
init_text: (cd src/examples/casting/; go test -run=x -bench=. ./cmd/byteslice_to_string/bench/... -benchmem -benchtime=2s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Casting `[]byte` to a `string`

## The direct cast

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 16
    end: 22
```

Here, similar to casting over a custom list, we create a var `s` and tell it to read `bb`'s data, but as a `string`. There is however, one interesting quirk.

<!-- stop -->

We are casting a slice to `string`, thus casting a `SliceHeader` (24 bytes) to a `StringHeader` (16 bytes). This results in the final 8 bytes of the `SliceHeader` (the capacity) being ignored by our new string.

<!-- stop -->

```go
                +---------------------------------------------------------------------------------------+
 bb starts -->  | 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
                +----------+----------+----------+----------+----------+----------+----------+----------+
 s starts --->  | 00000000 | 10000000 | 00001011 | 00000000 | 11000000 | 0        | 0        | 0        |
                +----------+----------+----------+----------+----------+----------+----------+----------+
                +---------------------------------------------------------------------------------------+
                | 0xc..160 | 0xd..161 | 0xc..162 | 0xc..163 | 0xc..164 | 0xc..165 | 0xc..166 | 0xc..167 | <- s ends
                +----------+----------+----------+----------+----------+----------+----------+----------+
                | 00001011 | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
                +----------+----------+----------+----------+----------+----------+----------+----------+
                +---------------------------------------------------------------------------------------+
                | 0xc..168 | 0xd..169 | 0xc..16a | 0xc..16b | 0xc..16c | 0xc..16d | 0xc..16e | 0xc..16d | <- bb ends
                +----------+----------+----------+----------+----------+----------+----------+----------+
                | 00001011 | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
                +----------+----------+----------+----------+----------+----------+----------+----------+
```

<!-- stop -->

This is mostly fine.

---

# Casting `[]byte` to a `string`

## Using unsafe package helpers

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 23
    end: 28
```

This approach, under the hood, is similar to casting a `string` to a custom type (such as `MyCoolType`).

<!-- stop -->

`unsafe.SliceData` takes a slice and returns a pointer to the starting element of that slice (calculated from `&bb[:1][0]`), so here it will return a pointer to where `72` is in memory. Its functionality similar to `&bb[0]` but with a couple of safety checks around `nil`.

<!-- stop -->

`unsafe.String` takes a pointer and a length, and then builds and returns a new `StringHeader`, with the provided pointer being used for `StringHeader.Data`, and the length as `StringHeader.Len`.

<!-- stop -->

This results in a `string` being created which points at the same underlying data as `bb`.

```go
          +----------------+                       +----------+----------+----------+
  bb -->  | Data: 0xc..158 | ----------+---------> | 0xc..158 | 0xc....  | 0xc..162 |
          | Len: 11        |           |           +----------+----------+----------+
          +----------------+           |           | 'H'      | . . .    | 'd'      |
                                       |           +----------+----------+----------+
                                       |
          +----------------+           |
  s --->  | Data: 0xc..158 | ----------+
          | Len: 11        |
          +----------------+
```

---

# Casting `[]byte` to a `string`

## Manually creating the header

Let's go through this line by line.

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 29
    end: null
```

---

# Casting `[]byte` to a `string`

## Manually creating the header

First we create our slice of bytes.

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g;31s/    / >  /'
lines:
    start: 29
    end: null
```

```go
          +----------------+                       +----------+----------+----------+
  bb -->  | Data: 0xc..158 | --------------------> | 0xc..158 | 0xc....  | 0xc..162 |
          | Len: 11        |                       +----------+----------+----------+
          +----------------+                       | 'H'      | . . .    | 'd'      |
                                                   +----------+----------+----------+
```

---

# Casting `[]byte` to a `string`

## Manually creating the header

Next we create an empty string. This creates a zeroed `StringHeader`, pointing to `0` with a `Len` of `0`.

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g;33s/    / >  /'
lines:
    start: 29
    end: null
```

```go
          +----------------+                       +----------+----------+----------+
  bb -->  | Data: 0xc..158 | --------------------> | 0xc..158 | 0xc....  | 0xc..162 |
          | Len: 11        |                       +----------+----------+----------+
          +----------------+                       | 'H'      | . . .    | 'd'      |
                                                   +----------+----------+----------+

          +----------------+
  s --->  | Data: 0        |
          | Len: 0         |
          +----------------+
```

---

# Casting `[]byte` to a `string`

## Manually creating the header

We cast `s` and `bb` to string and slice headers, so their binary data is accessable.

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g;34,35s/    / >  /'
lines:
    start: 29
    end: null
```

```go
          +----------------+                       +----------+----------+----------+
  bb -->  | Data: 0xc..158 | --------------------> | 0xc..158 | 0xc....  | 0xc..162 |
          | Len: 11        |                       +----------+----------+----------+
          +----------------+                       | 'H'      | . . .    | 'd'      |
                                                   +----------+----------+----------+

          +----------------+
  s --->  | Data: 0        |
          | Len: 0         |
          +----------------+
```

---

# Casting `[]byte` to a `string`

## Manually creating the header

We assign `bb`'s `SliceHeader.Data` to `s`'s `SliceHeader.Data`.

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g;37s/    / >  /'
lines:
    start: 29
    end: null
```

```go
          +----------------+                       +----------+----------+----------+
  bb -->  | Data: 0xc..158 | ---------+----------> | 0xc..158 | 0xc....  | 0xc..162 |
          | Len: 11        |          |            +----------+----------+----------+
          +----------------+          |            | 'H'      | . . .    | 'd'      |
                                      |            +----------+----------+----------+
                                      |
          +----------------+          |
  s --->  | Data: 0xc..158 |----------+
          | Len: 0         |
          +----------------+
```

---

# Casting `[]byte` to a `string`

## Manually creating the header

We assign `bb`'s `SliceHeader.Len` to `s`'s `SliceHeader.Len`.

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g;38s/    / >  /'
lines:
    start: 29
    end: null
```

```go
          +----------------+                       +----------+----------+----------+
  bb -->  | Data: 0xc..158 | ---------+----------> | 0xc..158 | 0xc....  | 0xc..162 |
          | Len: 11        |          |            +----------+----------+----------+
          +----------------+          |            | 'H'      | . . .    | 'd'      |
                                      |            +----------+----------+----------+
                                      |
          +----------------+          |
  s --->  | Data: 0xc..158 |----------+
          | Len: 11        |
          +----------------+
```

---

# Casting `[]byte` to a `string`

## So which should I use?

The safest of all of these is the third method, casting to the header, just because we have the most control over the casting.

However we should probably take inspirition from the std lib here.

<!-- stop -->

Here is the code for `strings.Clone` as of `go version go1.20.3 linux/amd64`:

```go
func Clone(s string) string {
    if len(s) == 0 {
        return ""
    }
    b := make([]byte, len(s))
    copy(b, s)
    return unsafe.String(&b[0], len(b))
}
```

<!-- stop -->

## Moral of the story

We're already far off the beaten track so who cares?

---

# Casting a `string` to `[]byte`

This is largely similar to casting a `[]byte` to a `string`, but there is one massive difference. One of the methods is invalid and, while creating code that will execute, is __very__ dangerous.

<!-- stop -->


```file
path: src/examples/casting/cmd/string_to_byteslice/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 16
    end: 22
```


```file
path: src/examples/casting/cmd/string_to_byteslice/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 23
    end: 28
```


```file
path: src/examples/casting/cmd/string_to_byteslice/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 29
    end: null
```

---

# Casting a `string` to `[]byte`

## And the winner is...

<!-- stop -->

```file
path: src/examples/casting/cmd/string_to_byteslice/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 16
    end: 22
```

<!-- stop -->

Why? Those pesky headers again.

---

# Casting a `string` to `[]byte`

## Casting a `StringHeader` to a `SliceHeader`.

Once agian, a `StringHeader` is 16 bytes, and a `SliceHeader` is 24 bytes.

<!-- stop -->

So, if we are to cast a `StringHeader` to a `SliceHeader`, our read will overshoot by 8 bytes, and will read random data into the slice's capacity.

<!-- stop -->


```go
                +---------------------------------------------------------------------------------------+
  s starts -->  | 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
                +----------+----------+----------+----------+----------+----------+----------+----------+
 bb starts -->  | 00000000 | 10000000 | 00001011 | 00000000 | 11000000 | 0        | 0        | 0        |
                +----------+----------+----------+----------+----------+----------+----------+----------+
                +---------------------------------------------------------------------------------------+
                | 0xc..160 | 0xd..161 | 0xc..162 | 0xc..163 | 0xc..164 | 0xc..165 | 0xc..166 | 0xc..167 | <- s ends
                +----------+----------+----------+----------+----------+----------+----------+----------+
                | 00001011 | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
                +----------+----------+----------+----------+----------+----------+----------+----------+
                +---------------------------------------------------------------------------------------+
                | 0xc..168 | 0xd..169 | 0xc..16a | 0xc..16b | 0xc..16c | 0xc..16d | 0xc..16e | 0xc..16d | <- bb ends
                +----------+----------+----------+----------+----------+----------+----------+----------+
                |   ????   |   ????   |   ????   |   ????   |   ????   |   ????   |   ????   |   ????   |
                +----------+----------+----------+----------+----------+----------+----------+----------+
```

This can lead to all sorts of chaos, from undefined behaviour to memory leaks to data corruption.

---

# Any questions so far?

---

# Avoiding allocations

One of the sad parts of go, for me, is balancing these three facts:

<!-- stop -->

1. Interfaces are awesome are allow for easy mocking.
<!-- stop -->

1. Allocating memory to the heap is very expensive.
<!-- stop -->

1. Passing a pointer to an interface always allocates.

<!-- stop -->

So, if we want to write easily mockable services, we may have to accept that we'll be poking the GC some more.

---

# Avoiding allcations

## But WHY does it allocate?

From the compiler's point of view it's very simple.

Let's say you have an interface:

```go
type Relayer interface {
    Relay(ctx context.Context, msg *Message) error
}
```

<!-- stop -->

Let's call this `Relay` function:

```go
func doRelay(ctx context.Context, relayer Relayer, msg *Message) error {
    return relayer.Relay(ctx, msg)
}
```

We have no idea what this `relayer` is doing with `msg`.

---

# Avoiding allcations

## But WHY does it allocate?

It might save the message:

```go
func (s *hyperthymesiaRelay) Relay(ctx context.Context, msg *Message) error {
    // some code

    s.relayedMessages = append(s.relayedMessages, msg)

    // some more code
}
```

---

# Avoiding allcations

## But WHY does it allocate?

It might fire it down a channel to who knows where:

```go
func (s *bufferedRelay) Relay(ctx context.Context, msg *Message) error {
    // some code

    s.buf = <-msg

    // some more code
}
```

---

# Avoiding allcations

## But WHY does it allocate?

It might fire it off in a go routine.

```go
func (s *backgroundRelay) Relay(ctx context.Context, msg *Message) error {
    // some code

    go handle(msg)

    // some more code
}
```

<!-- stop -->

The compiler cannot guarantee that the pointer won't live longer than, or remain inside, the scope of the interface's function.

---

# Avoiding allocations

## Can anything be done?

Against all odds, yes.

<!-- stop -->

Though you should probably never do this.

---

# Avoiding allocations

## Can anything be done?

Take the following __very__ contrived code.

We have service that handles `Person` entities. This service calls to a repo, and this repo saves the `Person` entities in a `Store`.

```file
path: src/examples/heap/person.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 14
    end: 29
```

---

# Avoiding allocations

## Can anything be done?

The implementation of the store uses a map in a non-threadsafe way, because demo.

```file
path: src/examples/heap/store/personstore.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 11
    end: null
```

---

# Avoiding allocations

## Can anything be done?

We implement the repo and service as you would expect.

```file
path: src/examples/heap/repo/person.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 13
    end: 34
```
---

# Avoiding allocations

## Can anything be done?

We implement the repo and service as you would expect.

```file
path: src/examples/heap/service/person.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 13
    end: 25
```

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 25
init_text: (cd src/examples/heap/; go test -run=x -bench=Safe ./service/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---

# Avoiding allocations

## The unsafe way

So I have the exact same repo and service interfaces implemented in an unsafe way.

First we'll look at the benchmarks, then some background. Then we will look at the code.

<!-- stop -->

```terminal-ex
command: zsh -il
rows: 25
init_text: (cd src/examples/heap/; go test -run=x -bench=Unsafe ./service/... -benchmem -benchtime=5s)
init_wait: '> '
init_codeblock_lang: zsh
```

---


# Avoid allocations

## That's neat! But how far from gods light must we stray for this?

<!-- stop -->

Very.

<!-- stop -->

For this we leverage two compiler directives. The first being `//go:noescape`

---

# Avoiding allocations

## //go:noescape

`go:noescape` must be placed before a function is declared, and it disables escape analysis on a function's parameters and return values. Meaning any pointers passed in won't leak to the heap.

There is, however, a catch. Straight from the docs:

```
The //go:noescape directive must be followed by a function declaration without a body (meaning that the function has an implementation not written in Go).
```

Example:

```go
//go:noescape
func hello(str *string) error

//go:noescape
func goodbye(str *string) error
```

---

# Avoiding allocations

## //go:linkname

This directive is pretty cool, but like all cool things, it can only be used if you import `unsafe` because it completely ignores the type system.

This lets us stub the body of one function out with another.

<!-- stop -->

Let's say we have this following func:

```file
path: src/examples/linkname/util/util.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 5
    end: null
```

<!-- stop -->

With the following test:
```file
path: src/examples/linkname/util/util_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 15
    end: null
```

<!-- stop -->

What if we got the compliler to overwrite the body of `rand.Int` with a custom func?
```file
path: src/examples/linkname/util/util_test.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 11
    end: 15
```

---

# Avoiding allocations

## WAT

Another use, using `time.Sleep` as an example.

Has anyone ever gotten curious and looked at the implementation of `time.Sleep` in the std lib?

<!-- stop -->

```go
// Sleep pauses the current goroutine for at least the duration d.
// A negative or zero duration causes Sleep to return immediately.
func Sleep(d Duration)
```

Very underwhelming.

<!-- stop -->

BUT, if we do a bit more digging, we find something interesting inside the `runtime` package:

```go
// timeSleep puts the current goroutine to sleep for at least ns nanoseconds.
//
//go:linkname timeSleep time.Sleep
func timeSleep(ns int64) {
    if ns <= 0 {
        return
    }

    gp := getg()
    t := gp.timer
    if t == nil {
        t = new(timer)
        gp.timer = t
    }
    t.f = goroutineReady
    t.arg = gp
    t.nextwhen = nanotime() + ns
    if t.nextwhen < 0 { // check for overflow.
        t.nextwhen = maxWhen
    }
    gopark(resetForSleep, unsafe.Pointer(t), waitReasonSleep, traceEvGoSleep, 1)
}
```

`time.Sleep` has no body, and `go:linkname` let us stub a body in. See where this is going?

---

# Avoiding allocations

## Bringing it all together.

Even though `//go:noescape` forces us to create a bodyless function, it doesn't forbid the body of that function to be stubbed with `//go:linkname`.

So, let's take a look at our unsafe implementations of the small person service again.

<!-- stop -->

```file
path: src/examples/heap/service/person.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 30
    end: null
```

---

# Avoiding allocations

## Bringing it all together.

And the repo layer.

<!-- stop -->

```file
path: src/examples/heap/repo/person.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 38
    end: null
```

---

# ....and that's us.

Sources:
- https://go101.org
- https://github.com/golang/go

Github:
- https://github.com/tigh-latte

Thanks!

