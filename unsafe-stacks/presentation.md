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

# What about these functions are `unsafe`?

Not much on their own.

Rob Pike has even opened a proposal to move these functions into another package in go2: https://github.com/golang/go/issues/5602.

<!-- stop -->

The real lunacy begins when you use `unsafe.Pointer`.

---

# What is `unsafe.Pointer`?

The standard library defines this as:

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
 &i points to                                                                golang reads to
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
                                                                            golang reads to
```

---

# What is `unsafe.Pointer`?

So if we convert a typed pointer to an `unsafe.Pointer`, we are pretty much just pointing to that single memory address.

```go
    i := 34562
    p := unsafe.Pointer(&i) // Mem addr: 0xc000020158
```

```
 &i points to                                                                &i is read to
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

So if we convert a typed pointer to an `unsafe.Pointer`, we are pretty much just pointing to that single memory address.

```go
    i := 34562
    p := unsafe.Pointer(&i) // Mem addr: 0xc000020158
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

With this `unsafe.Pointer` we can have the binary data stored at `*i` be read
To dereference an `unsafe.Pointer`, we first cast it to a type, and then the data is read as it would be for this cast type.

<!-- stop -->

Continuing on from the previous example, using `unsafe.Pointer` as a medium, we can read `i` as a `string`:

```go
    i := 34562
    str := *(*string)(unsafe.Pointer(&i))
```

```
 &i points to                                                                &i is read to
 &s points to                                                                     |
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
                                                                            &s is read to
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

The uses of this for performance come from other data structures, which we will get into later.

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
 i points to                                                                   i reads to
 b points to
 f points to                                                                   f reads to
 a points to                                                                   a reads to
     |                                                                             |
     V                                                                             V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 0        | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
     ^
     |
 b reads to
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
 i points to                                                                   i reads to
 b points to
 f points to                                                                   f reads to
 a points to                                                                   a reads to
     |                                                                             |
     V                                                                             V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00001010 | 0        | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
     ^
     |
 b reads to
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
 i points to                                                                   i reads to
 b points to
 f points to                                                                   f reads to
 a points to                                                                   a reads to
     |                                                                             |
     V                                                                             V
+---------------------------------------------------------------------------------------+
| 0xc..158 | 0xc..159 | 0xc..15a | 0xc..15b | 0xc..15c | 0xc..15d | 0xc..15e | 0xc..15f |
+----------+----------+----------+----------+----------+----------+----------+----------+
| 00000000 | 00000001 | 0        | 0        | 0        | 0        | 0        | 0        |
+----------+----------+----------+----------+----------+----------+----------+----------+
     ^
     |
 b reads to
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
 i points to                                                                   i reads to
 b points to
 f points to                                                                   f reads to
 a points to                                                                   a reads to
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

This is an integer representation of an unsafe pointer, letting us do pointer arithmetic.

```go
    i := new(int)  // 0xc000020158
    *i = 256
    ptr := uintptr(unsafe.Pointer(i)) // 824633852248
    ptr++ // 824633852249
    b := *(*byte)(unsafe.Pointer(ptr)) // 0xc000020159 (second byte of i, 1)
```

---

# Then we have `uintptr`

Let's explain this real quick.

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

We want to pass a slice of this type, `[]MyCoolType`, to `strings.Join`. What can we do?

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
    start: 6
    end: 55
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
    start: 56
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

<!-- stop -->

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

And likewise, a struct representation of a sliced header exists in `reflect`:

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

I'm sure you's can guess, this clone isn't cheap.

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

# Method 3

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 29
    end: null
```

---

# Method 3

```go
type StringHeader struct {
    Data uintptr
    Len  int
}
```

<!-- stop -->

In memory, this looks like so:

### Data

| 0xc..01 | 0xc..02 | 0xc..03 | 0xc..04 | 0xc..05 | 0xc..06 | 0xc..07 | 0xc..08 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 80      | 1       | 2       | 0       | 192     | 0       | 0       | 0       |

### Len

| 0xc..09 | 0xc..0a | 0xc..0b | 0xc..0c | 0xc..0d | 0xc..0e | 0xc..0f | 0xc..10 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 12      | 0       | 0       | 0       | 0       | 0       | 0       | 0       |

<!-- stop -->

When you get the memory address of a string, you actually get the memory address of this header.

```go
    s := "Hello world!"
    _ = &s // {Data: 0xc456...01, Len: 12}
```

---

# Method 3

```go
type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
```

<!-- stop -->

In memory, this looks like so:

### Data

| 0xc..01 | 0xc..02 | 0xc..03 | 0xc..04 | 0xc..05 | 0xc..06 | 0xc..07 | 0xc..08 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 80      | 1       | 2       | 0       | 192     | 0       | 0       | 0       |

### Len

| 0xc..09 | 0xc..0a | 0xc..0b | 0xc..0c | 0xc..0d | 0xc..0e | 0xc..0f | 0xc..10 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 12      | 0       | 0       | 0       | 0       | 0       | 0       | 0       |

### Cap

| 0xc..11 | 0xc..12 | 0xc..13 | 0xc..14 | 0xc..15 | 0xc..16 | 0xc..17 | 0xc..18 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 16      | 0       | 0       | 0       | 0       | 0       | 0       | 0       |

<!-- stop -->

When you reference a slice, you actually get the memory address of this header.

```go
    bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
    _ = bb // {Data: 0xcff0...ab, Len: 12, Cap: 16}
```

---

# Method 3

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 29
    end: null
```

<!-- stop -->

---

# Method 2

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 23
    end: 28
```



---

# Method #1

```file
path: src/examples/casting/cmd/byteslice_to_string/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 16
    end: 22
```

`&bb` doesn't pull back the memory address holding `{72, 101, ...}`, but instead points to the __slice header__. A slice header spans 24 memory addresses, holding three pieces of information, each 8 addresses long.

<!-- stop -->

### Addrs 0-7

The "data pointer"; the memory address of the underlying array.

| 0xc..01 | 0xc..02 | 0xc..03 | 0xc..04 | 0xc..05 | 0xc..06 | 0xc..07 | 0xc..08 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 80      | 1       | 2       | 0       | 192     | 0       | 0       | 0       |

<!-- stop -->

### Addrs 8-15

The length of the underlying array.

| 0xc..09 | 0xc..0a | 0xc..0b | 0xc..0c | 0xc..0d | 0xc..0e | 0xc..0f | 0xc..10 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 12      | 0       | 0       | 0       | 0       | 0       | 0       | 0       |

<!-- stop -->

### Addrs 16-23

The capacity of the underlying array.

| 0xc..09 | 0xc..0a | 0xc..0b | 0xc..0c | 0xc..0d | 0xc..0e | 0xc..0f | 0xc..10 |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 12      | 0       | 0       | 0       | 0       | 0       | 0       | 0       |

<!-- stop -->

A string header is only 16 addresses long, needing a data pointer and a length.


---

`bb` points to the starting memory address of the slice of bytes, and because this is a slice, golang will read the next 23 memory addresses as well. These 24 bytes build the slice header.

<!-- stop -->

The first 8 bytes build a memory, pointing to the slice's underlying array.

| 0x00 | 0x01 | 0x02 | 0x03 | 0x04 | 0x05 | 0x06 | 0x07 |
|------|------|------|------|------|------|------|------|
| 80   | 1    | 2    | 0    | 192  | 0    | 0    | 0    |


---

# Working around this

If we don't care about this integrity (because, say, our slice of bytes is about to go out of scope), we can use `unsafe.String` to force a cast.


<!-- stop -->

IF YOU ADVANCE THIS DIES

<!-- stop -->

hello

<!--
things to talk about

byte -> string conversion being BLAZING fast
string -> byte conversion being class but longer

compiler no escape hack

-->

this is a test
