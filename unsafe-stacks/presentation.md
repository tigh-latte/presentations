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

# Explaining the output of these casts

```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 15
    end: 19
```

<!-- stop -->

When we write `10` to `*i`, because `*i` is an `*int`, `10` is written as an `int` would be.

Let's imagine `*i` points to memory address `0xc000020158`. Writing `10` to memory touches this memory address and the 7 following it.

<!-- stop -->


| 0xc..58 | 0xc..59 | 0xc..5a | 0xc..5b | 0xc..5c | 0xc..5d | 0xc..5e | 0xc..5f |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 1010    | 0       | 0       | 0       | 0       | 0       | 0       |  0      |

---

### The same thing applies to `256`.

```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 20
    end: 25
```

<!-- stop -->

| 0xc..58 | 0xc..59 | 0xc..5a | 0xc..5b | 0xc..5c | 0xc..5d | 0xc..5e | 0xc..5f |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 0...    | 1       | 0       | 0       | 0       | 0       | 0       |  0      |

---

### Writing `10` again but this time to a `*float64`

```file
path: src/examples/unsafe_ptr/cmd/int_to_bool/main.go
lang: go
transform: sed 's/\t/    /g'
lines:
    start: 30
    end: 34
```

<!-- stop -->

| 0xc..58 | 0xc..59 | 0xc..5a | 0xc..5b | 0xc..5c | 0xc..5d | 0xc..5e | 0xc..5f |
|---------|---------|---------|---------|---------|---------|---------|---------|
| 0       | 0       | 0       | 0       | 0       | 0       | 100100  | 1000000 |

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

# Casting `[]byte` to a `string`

Golang, by design, doesn't have many built in operations that are `O(N)`, as this can lead to writing expensive code without even realising.

Casting `[]byte` to `string` is one of these few exceptions (and visa-versa).

<!-- stop -->

Given that a `string` in golang is immutable but a `[]byte` isn't, in order to ensure integrity, `string([]byte)` makes a full clone of the slice of bytes, converting that clone to a `string`.

```go
    bb := []byte("Hello world!")
    s := string(bb)
    bb[0] = 'z'
    fmt.Println(s) // Output: Hello world
```

---

# Explaination

```go
    bb := []byte("Hello world!") // data starts at 0x20
```

Create data in memory.

| 0x20 | 0x21 | 0x22 | 0x23 | 0x24 | 0x25 | 0x26 | 0x27 | 0x28 | 0x29 | 0x2a | 0x2b |
|------|------|------|------|------|------|------|------|------|------|------|------|
| H    | e    | l    | l    | o    |      | w    | o    | r    | l    | d    | !    |

<!-- stop -->

```go
    s := string(bb) // data starts at 0xa3
```

We clone that data and have golang treat it as a `string`.

| 0xa3 | 0xa4 | 0xa5 | 0xa6 | 0xa7 | 0xa8 | 0xa9 | 0xaa | 0xab | 0xac | 0xad | 0xae |
|------|------|------|------|------|------|------|------|------|------|------|------|
| H    | e    | l    | l    | o    |      | w    | o    | r    | l    | d    | !    |

<!-- stop -->

Modifications to the `[]byte` don't modify the cloned string.

---

# The unsafe way

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
    bb := []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33}
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

[]custom -> []string conversion

compiler no escape hack

-->

this is a test
