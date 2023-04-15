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

<!--stop -->

To dereference an `unsafe.Pointer`, we must first cast it to a type, and then the data is read as it would be for this cast type.

<!-- stop -->

Take the following:
```go
i := new(int64) // 0xf40000002c

fmt.Println(*i) // Output: 0
```

When `i` is dereferenced, golang dereferences it as an `int64`, meaning `0xf40000002c` and the following 7 memory addresses are read too. The binary data is read as a signed integer.

Using `unsafe.Pointer` as a medium, we cast this integer pointer to pointer of another type, letting us read the data in memory __as that type__.

---

# Casting an `unsafe.Pointer`


```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 6
    end: null
```

This comes in very handy, as we'll see shortly.

---

# Then we have `uintptr`

This is an integer representation of an unsafe pointer, letting us do pointer arithmetic.

```go
i := new(int)  // 0x40ffab45
*i = 256
ptr := uintptr(unsafe.Pointer(i)) // 40ffab45
ptr++ // 40ffab46
b := *(*byte)(unsafe.Pointer(ptr)) // 0x40ffab46 (second byte of i, 1)
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
transform: sed 's/\t/    /g'
lines:
    start: 8
    end: null
```

```terminal-ex
command: zsh -il
rows: 10
init_text: cd src/examples/ptr_math/; go run cmd/gc_danger/main.go
init_wait: '> '
init_codeblock_lang: zsh
```

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
